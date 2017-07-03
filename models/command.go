package models

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"syscall"
	"time"

	"github.com/astaxie/beego"
)

// ErrTimeout for server request timeout
var ErrTimeout = errors.New("Timeout")

// ErrSizeout for size too big
var ErrSizeout = errors.New("SizeOut")

// Command may used for /api/command
type Command struct {
	Name    string   `json:"name"`
	Args    []string `json:"args"`
	Dir     string   `json:"Dir"`
	Env     []string `json:"env"`
	Timeout int      `json:"timeout"`
	Silence bool     `json:"Silence"`
	Commander
	*User
}

type CommandShell struct {
	Command string `json:"command"`
}

// CommandResult for Command.Run return
type CommandResult struct {
	Retcode int    `json:"retcode"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

type Commander interface {
	RunCommand(*Command) (*CommandResult, error)
}

func (cmd *Command) Run() (result *CommandResult, err error) {
	if cmd.Commander != nil {
		return cmd.Commander.RunCommand(cmd)
	}
	return modelsCommander.RunCommand(cmd)
}

func NewCommand() *Command {
	return &Command{Commander: modelsCommander}
}

var modelsCommander Commander

func ShellRun(shellcmd string) (result *CommandResult, err error) {
	cmd := Command{Name: "/bin/bash", Args: []string{"-c", shellcmd}, Commander: modelsCommander}
	return cmd.Run()
}

func init() {
	// if  ....
	// else ...
	cmderD := &CommanderD{}
	setuid := beego.AppConfig.DefaultInt("local::setuid", 0)
	cmderD.setUid = uint32(setuid)
	modelsCommander = cmderD
	log.Printf("init ok : %#v \n", modelsCommander)
}

type CommanderD struct {
	setUid uint32
}

func (c *CommanderD) RunCommand(cmd *Command) (result *CommandResult, err error) {
	uid := cmd.User.GetUID()
	if uid == 0 {
		uid = c.setUid
	}
	return c.RunCommandUID(cmd, uid)
}

func (c *CommanderD) RunCommandUID(cmd *Command, uid uint32) (result *CommandResult, err error) {
	result = &CommandResult{}
	var outbuf, errbuf bytes.Buffer
	sigterm := make(chan int, 1)
	ch := make(chan int, 1)
	const (
		timeout = iota
		sizeout = iota
	)
	oscmd := exec.Command(cmd.Name, cmd.Args...)
	if !cmd.Silence {
		oscmd.Stdout, oscmd.Stderr = &outbuf, &errbuf
	}
	// if setuid = true
	if uid != 0 {
		//if c.setUid != 0 {
		//uid := c.setUid
		sysattr := &syscall.SysProcAttr{
			Credential: &syscall.Credential{Uid: uint32(uid)},
		}
		oscmd.SysProcAttr = sysattr
	}
	oscmd.Env = cmd.Env
	oscmd.Dir = cmd.Dir
	if err = oscmd.Start(); err != nil {
		log.Printf("cmd.Start: %v", err)
		result.Error = err.Error()
		result.Retcode = 127
		return
	}
	// handle timeout

	go func() {
		mytimeout := 60
		// BUG(tonge) TO-DO: read timeoutlimit from config
		//mytimeout, timeouterr := beego.GetConfig(int, "limitcmdtimeout", 5)
		//mytimeout, timeouterr := beego.AppConfig.Int("httpport")
		//if timeouterr != nil {
		//      log.Printf("timeout conf err : %#v \n", timeouterr)
		//}
		if 0 < cmd.Timeout && cmd.Timeout < mytimeout {
			mytimeout = cmd.Timeout
		}
		time.Sleep(time.Duration(mytimeout) * time.Second)
		sigterm <- timeout
	}()

	// handle output ?
	if !cmd.Silence {
		go func() {
			for {
				time.Sleep(1 * time.Second)
				sizelimit := (2 << 18)
				// BUG(tonge) TO-DO: read sizelimit from config
				//log.Printf("len of buf: %#v %#v \n", outbuf.Len(), errbuf.Len())
				if outbuf.Len()+errbuf.Len() > sizelimit {
					sigterm <- sizeout
				}
			}
		}()
	}
	go func() {
		if err := oscmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					result.Retcode = status.ExitStatus() //log.Printf("Exit Status: %d" , status.ExitStatus())
				}
			} else {
				result.Retcode = 126
			}
		}
		ch <- 0
	}()
	select {
	case <-ch:
		result.Output, result.Error = outbuf.String(), errbuf.String()
	case endsig := <-sigterm:
		switch endsig {
		case timeout:
			err = ErrTimeout //result.Error = "timeout"
		case sizeout:
			err = ErrSizeout //result.Error = "sizeout"
		}
	}
	return
}
