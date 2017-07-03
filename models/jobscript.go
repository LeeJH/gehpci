package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"log"
	"os"
	"reflect"
)

type shellScheduler struct {
	shelldir string
}

func (ss *shellScheduler) run(filename string, mydata interface{}) error {
	//log.Println("start " + filename)
	shellcmd := ss.shelldir + "/" + filename //path.Join(ss.shelldir, filename)
	result, err := ShellRun(shellcmd)
	log.Printf("Run shell %s : %d \n %s \n", shellcmd, result.Retcode, result.Output)
	if err != nil {
		return err
	}
	if result.Retcode != 0 {
		return errors.New("Error " + string(result.Retcode) + " : " + result.Error)
	}
	switch mydata.(type) {
	case *string:
		mydata = result.Output
	default:
		json.Unmarshal([]byte(result.Output), mydata)
	}
	log.Printf("Unmarshal result:\n %#v \n", mydata)
	return nil
}

func (ss *shellScheduler) Info() ([]*HPCResource, error) {
	log.Println("start info()")
	resl := make([]*HPCResource, 0)
	err := ss.run("info", &resl)
	return resl, err
}

func (ss *shellScheduler) Queue() ([]*HPCJob, error) {
	jobl := make([]*HPCJob, 0)
	err := ss.run("queue", &jobl)
	return jobl, err
}

func genEnv(prefix string, obj interface{}) (envgen []string) {
	envgen = make([]string, 0)
	s := reflect.ValueOf(obj).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !f.CanInterface() {
			continue
		}
		envi := fmt.Sprintf("%s%s=%v", prefix, s.Type().Field(i).Name, f.Interface())
		envgen = append(envgen, envi)
	}
	return envgen
}

func (ss *shellScheduler) Submit(job *HPCJob) (jobid string, err error) {
	log.Printf("log here")
	cmd := &Command{}
	cmd.Name = ss.shelldir + "/submit"
	//cmd.Args = []string{}
	envgen := genEnv("gehpci_", job)
	cmd.Env = append(os.Environ(), envgen...)
	result, err := cmd.Run()

	if err != nil {
		return
	}
	if result.Retcode != 0 {
		return "", errors.New("Error " + string(result.Retcode) + " : " + result.Error)
	}

	return result.Output, nil
}
func (ss *shellScheduler) JobInfo(jobid string) ([]*HPCJob, error) {
	jobss := make([]*HPCJob, 0)
	err := ss.run("jobinfo "+jobid, &jobss)
	return jobss, err
}
func (ss *shellScheduler) Delete(jobid string) error {
	var resultstr string
	err := ss.run("delete "+jobid, &resultstr)
	return err
}
func (ss *shellScheduler) Name() string {
	return "shell-Scheduler-ej"
}
