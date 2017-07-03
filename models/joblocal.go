package models

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
)

type localScheduler struct {
	useSSH     bool
	machines   []string
	joblist    []*HPCJob
	jobqueue   []*HPCJob
	lock       sync.RWMutex
	idmax      int
	returnFile bool
}

func (ss *localScheduler) Name() string {
	return fmt.Sprintf("localScheduler(ssh:%v)", ss.useSSH)
	//return "localScheduler(ssh:" + string(ss.useSSH) + ")"
}

func (ss *localScheduler) failJob(id int, info string) {
	log.Printf("job %d failed : %s \n", id, info)
	ss.joblist[id-1].JobState = "Failed"
	ss.removeQueue(id)
}

func (ss *localScheduler) removeQueue(id int) {
	strid := strconv.Itoa(id)
	ss.lock.Lock()
	queuelen := len(ss.jobqueue)
	for i, job := range ss.jobqueue {
		if job.JobID == strid {
			if queuelen == 1 {
				ss.jobqueue = make([]*HPCJob, 0)
			} else if i == 0 {
				ss.jobqueue = ss.jobqueue[1:]
			} else if i == queuelen-1 {
				ss.jobqueue = ss.jobqueue[:i]
			} else {
				ss.jobqueue = append(ss.jobqueue[:i], ss.jobqueue[i+1:]...)
			}
		}
	}
	ss.lock.Unlock()
}

func (ss *localScheduler) Submit(job *HPCJob) (jobid string, err error) {
	log.Printf("log submit")
	ss.lock.Lock()
	ss.idmax++
	id := ss.idmax
	jobid = fmt.Sprintf("%d", id)
	job.JobState = "Pending"
	job.JobID = jobid
	ss.joblist = append(ss.joblist, job)
	ss.jobqueue = append(ss.jobqueue, job)
	job = ss.joblist[id-1]
	ss.lock.Unlock()
	if job.Name == "" {
		job.Name = path.Base(job.JobFile)
	}
	go func() {
		if !ss.useSSH {
			// just exec job at local machine
			var outbuf, errbuf bytes.Buffer
			cmd := exec.Command(job.JobFile, job.JobArgs)
			cmd.Dir = job.Dir
			// waiting for res ?
			// wait_for_res()
			if !ss.returnFile {
				cmd.Stdout, cmd.Stderr = &outbuf, &errbuf
			} else {
				stdoutfile := fmt.Sprintf("%s.out", job.JobID)
				stderrfile := fmt.Sprintf("%s.err", job.JobID)
				stdoutfile = path.Join(job.Dir, stdoutfile)
				stderrfile = path.Join(job.Dir, stderrfile)
				stdout, err := os.OpenFile(stdoutfile, os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					log.Fatalln(err)
				}
				defer stdout.Close()
				cmd.Stdout = stdout
				stderr, err := os.OpenFile(stderrfile, os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					log.Fatalln(err)
				}
				defer stderr.Close()
				cmd.Stderr = stderr

			}
			err = cmd.Start()
			if err != nil {
				ss.failJob(id, "Job start failed:"+err.Error())
				return
			}
			log.Printf("pid: %#v", cmd.Process)
			job.JobState = "Running"
			err = cmd.Wait()
			if err != nil {
				ss.failJob(id, "Job running failed:"+err.Error())
				return
			}
			//result.Output, result.Error = , errbuf.String()
			if !ss.returnFile {
				job.Infos += "\nOutput:\n" + outbuf.String() + "\nError:\n" + errbuf.String()
			}
			job.JobState = "Success"
			ss.removeQueue(id)
			// cmd.ProcessState.Pid
			// remove from queue
		} else {
			log.Fatal("not imp yet!")
			return
		}
	}()

	return //"job1", nil

}

////////////////////////

func (ss *localScheduler) Info() ([]*HPCResource, error) {
	resl := make([]*HPCResource, 2)
	resl[0] = &HPCResource{Partition: "example", Nodes: 1, State: "idle"}
	resl[1] = &HPCResource{Partition: "example", Nodes: 1, State: "down"}
	return resl, nil
}

func (ss *localScheduler) Queue() ([]*HPCJob, error) {
	//jobl := make([]HPCJob, 2)
	//jobl[0] = HPCJob{Name: "job1"}
	//jobl[1] = HPCJob{Name: "job2"}
	return ss.jobqueue, nil
}

func (ss *localScheduler) JobInfo(jobid string) ([]*HPCJob, error) {
	//#var exampljobss []HPCJob{ HPCJob{ Name: "job1"}}
	exampljobss := make([]*HPCJob, 1)
	id, err := strconv.Atoi(jobid)
	if err != nil {
		return exampljobss, nil
	}
	exampljobss[0] = ss.joblist[id-1]
	return exampljobss, nil
}
func (ss *localScheduler) Delete(jobid string) error {
	return nil
}
