package models

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"

	"time"

	"github.com/gorilla/websocket"
)

/*
debug js :

ws = new WebSocket("ws://localhost:9090/v1/wc/bindbash/local");

ws.onmessage=function(event){
    console.log(event.data)};

ws.send( JSON.stringify( {"Id":"2a","Kind":"alive"} ))

ws.send( JSON.stringify( {"Id":"2a","Kind":"run","Body":"ls"} ))
*/

/*
TODSs :
Options :
	not bash but other program
	not line buf but bytes
	return response...bash-x$
	more power new function
*/

type WSMessage struct {
	Id      string   // cliend-provied unique id for the process
	Kind    string   // in : run/kill/<alive> out: stdout/stderr/end/<alive>
	Body    string   // message body
	Options *Options `json:",omitempty"`
}

type Options struct {
	options []string
}

type WebSocketMan struct {
	ws           *websocket.Conn
	user         *User
	outputReader *bufio.Reader
	outerrReader *bufio.Reader
	stdin        io.WriteCloser
	basicinfo    string
	cmd          *exec.Cmd
	msgreq       WSMessage
	msgoutl      []byte
	msgerrl      []byte
	muxIn        sync.RWMutex
	muxOut       sync.RWMutex
	muxErr       sync.RWMutex
	chanoutput   chan int
	chanexit     chan struct{}
	limittime    time.Duration
	limitsize    uint32
	limitalive   int
}

const limittimeDef time.Duration = time.Second * 1
const limitsizeDef uint32 = 1024
const limitaliveDef int = 3600

func NewWebSocketMan(conn *websocket.Conn, user *User) *WebSocketMan {
	wsm := &WebSocketMan{}
	wsm.ws = conn
	wsm.user = user
	wsm.msgreq = WSMessage{}
	return wsm
}

const (
	CSRUNOVER = iota
	CSTIMEUPDATE
	CSSIZEUPDATE
	CSTIMEOUT
)

func (wsm *WebSocketMan) BindBash() (err error) {

	cmd := exec.Command("bash", "-i")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if wsm.user != nil && wsm.user.Uid != 0 {
		sysattr := &syscall.SysProcAttr{
			Credential: &syscall.Credential{Uid: wsm.user.Uid},
		}
		cmd.SysProcAttr = sysattr
	}
	wsm.cmd = cmd
	wsm.outputReader = bufio.NewReader(stdout)
	wsm.outerrReader = bufio.NewReader(stderr)
	wsm.stdin = stdin
	cmd.Start()
	defer stdin.Close()
	defer stdout.Close()
	defer stderr.Close()
	stdin.Write([]byte("\n"))
	wsm.outerrReader.ReadByte()
	wsm.outerrReader.UnreadByte()
	for {
		wsm.outerrReader.ReadByte()
		if wsm.outerrReader.Buffered() == 0 {
			break
		}
	}
	stdin.Write([]byte("\n"))
	wsm.basicinfo, _ = wsm.outerrReader.ReadString('\n')
	wsm.basicinfo = wsm.basicinfo[:len(wsm.basicinfo)-1]
	wsm.chanoutput = make(chan int)
	wsm.chanexit = make(chan struct{})
	defer func() {
		//if wsm.chanexit != nil {
		select {
		case _, ok := <-wsm.chanexit:
			if ok {
				close(wsm.chanexit)
			}
		default:
			close(wsm.chanexit)
		}
		//}
		//if wsm.chanoutput != nil {
		select {
		case _, ok := <-wsm.chanoutput:
			if ok {
				close(wsm.chanoutput)
			}
		default:
			close(wsm.chanoutput)
		}

		//}
	}()
	if wsm.limittime == 0 || wsm.limittime < limittimeDef {
		wsm.limittime = limittimeDef
	}
	if wsm.limitsize == 0 || wsm.limitsize > limitsizeDef {
		wsm.limitsize = limitsizeDef
	}
	if wsm.limitalive == 0 || wsm.limitalive > limitaliveDef {
		wsm.limitalive = limitaliveDef
	}
	//fmt.Printf("limits : %v %v \n", limitsize, limittime)

	go wsm.collectOutput()

	// collect error
	go wsm.collectOuterr()

	// output-timeout output
	go func() {
		for {
			select {
			case <-wsm.chanexit:
				return
			default:
				wsm.chanoutput <- CSTIMEUPDATE
				time.Sleep(wsm.limittime)
			}
			//log.Printf("goruntine output-timeout alive...\n %#v %#v", wsm.chanexit, wsm.chanoutput)
		}
	}()

	// close if the cmd is close.
	go func() {
		wsm.cmd.Wait()
		select {
		case <-wsm.chanexit:
			return
		default:
			wsm.CloseBind("bash exited.")
		}
	}()

	// output-noreq timeout
	go func() {
		for {
			select {
			case <-wsm.chanexit:
				return
			default:
				const timestep int = 1
				time.Sleep(time.Duration(timestep) * time.Second)
				if wsm.limitalive <= timestep {
					closeinfo := fmt.Sprintf("Error: timeout as long time no request.%d\n", limitaliveDef)
					wsm.CloseBind(closeinfo)
					return
				} else {

					wsm.limitalive = wsm.limitalive - timestep

				}
			}
			//log.Printf("goruntine output-response alive...\n")

		}

	}()

	// output collect
	go wsm.writeResp()

	// ListenWebSocket
	err = wsm.listenMsg()

	log.Printf("websocket out. %#v  %#v", err, wsm.cmd.Process)
	return err
}

func (wsm *WebSocketMan) listenMsg() error {

	for {
		select {
		case <-wsm.chanexit:
			return fmt.Errorf("exit by wsm.chanexit ")
		default:
			var msg WSMessage // Read in a new message as JSON and map it to a Message object
			err := wsm.ws.ReadJSON(&msg)
			if err != nil {
				return err
			}
			log.Printf("Got msg: %#v\n", msg)
			wsm.limitalive = limitaliveDef
			switch msg.Kind {
			case "run":
				wsm.muxIn.Lock()
				if msg.Id <= wsm.msgreq.Id {
					wsm.ws.WriteJSON(
						WSMessage{
							Id:   msg.Id,
							Kind: "info",
							Body: "u shoud use id > " + wsm.msgreq.Id})
					log.Printf("error id used.")
				} else {
					wsm.stdin.Write([]byte(msg.Body))
					wsm.msgreq = msg
					log.Printf("msg:run:%s", string(msg.Body))
				}
				wsm.muxIn.Unlock()
			case "kill":
				errkill := fmt.Errorf("killed by msg:kill(id:%s)", msg.Id)
				wsm.CloseBind(errkill.Error())
				return errkill
			case "alive": // not really need .
				wsm.ws.WriteJSON(
					WSMessage{
						Id:   msg.Id,
						Kind: "alive"})
			}
		}
	}
}

func (wsm *WebSocketMan) collectOutput() {
	for {
		select {
		case <-wsm.chanexit:
			return
		default:
			output, err := wsm.outputReader.ReadBytes('\n')
			//output, err := wsm.outputReader.ReadString('\n')
			if err != nil {
				break
			}
			wsm.muxOut.Lock()
			wsm.msgoutl = append(wsm.msgoutl, output...)
			wsm.muxOut.Unlock()
			if uint32(len(wsm.msgoutl)) >= wsm.limitsize {
				wsm.chanoutput <- CSSIZEUPDATE
			}

		}
	}
}

func (wsm *WebSocketMan) collectOuterr() {
	for {
		select {
		case <-wsm.chanexit:
			return
		default:
			outerr, err := wsm.outerrReader.ReadBytes('\n')
			if err != nil {
				break
			}
			wsm.muxErr.Lock()
			wsm.msgerrl = append(wsm.msgerrl, outerr...)
			wsm.muxErr.Unlock()
			wsm.outerrReader.ReadByte()
			wsm.outerrReader.UnreadByte()
			bn := wsm.outerrReader.Buffered()
			var sendchan bool = false
			if bn >= len(wsm.basicinfo) {
				pk, _ := wsm.outerrReader.Peek(bn)
				if string(pk[bn-len(wsm.basicinfo):bn]) == wsm.basicinfo {
					log.Printf("once run over")
					wsm.chanoutput <- CSRUNOVER
					sendchan = true
				}
			}
			if !sendchan {
				if uint32(len(wsm.msgerrl)) >= wsm.limitsize {
					wsm.chanoutput <- CSSIZEUPDATE
				}
			}
		}
	}

}

func (wsm *WebSocketMan) writeResp() {
	for {
		select {
		case outputsig, ok := <-wsm.chanoutput:
			if !ok {
				return
			}
			wsm.muxOut.RLock()
			if len(wsm.msgoutl) > 0 {
				wsm.muxIn.RLock()
				wsm.ws.WriteJSON(
					WSMessage{
						Id:   wsm.msgreq.Id,
						Kind: "stdout",
						Body: string(wsm.msgoutl)})
				wsm.muxIn.RUnlock()

			}
			wsm.muxOut.RUnlock()

			wsm.muxErr.RLock()
			if len(wsm.msgerrl) > 1 {
				wsm.muxIn.RLock()
				wsm.ws.WriteJSON(
					WSMessage{
						Id:   wsm.msgreq.Id,
						Kind: "stderr",
						Body: string(wsm.msgerrl)})
				wsm.muxIn.RUnlock()
			}
			wsm.muxErr.RUnlock()

			wsm.muxOut.Lock()
			wsm.muxErr.Lock()
			wsm.msgoutl = make([]byte, 0)
			wsm.msgerrl = make([]byte, 0)
			wsm.muxErr.Unlock()
			wsm.muxOut.Unlock()

			if outputsig == CSRUNOVER {

			}
		case <-wsm.chanexit:
			return
		}

	}

}

func (wsm *WebSocketMan) CloseBind(reason string) {
	select {
	case <-wsm.chanexit:
		// already close .
		return
	default:
		wsm.cmd.Process.Kill()
		wsm.ws.WriteJSON(
			WSMessage{
				Id:   wsm.msgreq.Id,
				Kind: "end",
				Body: reason})
		close(wsm.chanexit)
		wsm.ws.Close()

	}
}
