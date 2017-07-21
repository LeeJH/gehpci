package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

const (
	cmsgSTART int = iota
	cmsgSTOP
	cmsgRESFREE
	cmsgRESALLOC
	cmsgQSidle
	cmsgQSpending
	cmsgQSmatch
)

type localJob struct {
	HPCJob
	chanmsg chan int
}

// a simple local job scheduler
// job
type localScheduler struct {
	useSSH     bool
	machines   []string
	joblist    map[int64]*localJob // all jobs in log
	jobqueue   []int64             // pending jobs
	jobrunning map[int64]bool      // running jobs

	resource     map[string]int64 // hostname : tasks
	resourceFree map[string]int64 // hostname : tasks

	lock               sync.RWMutex
	idmax              int64
	returnFile         bool
	channewjob         chan *localJob
	chanqueue          chan int // need to check queues
	queuestate         int
	currentWaittingJob int // currentWaitting job in queue
	nodesize           int64
}

func (ss *localScheduler) Name() string {
	return fmt.Sprintf("localScheduler(ssh:%v)", ss.useSSH)
	//return "localScheduler(ssh:" + string(ss.useSSH) + ")"
}

func NewLocalScheduler() *localScheduler {
	ss := &localScheduler{}
	ss.resource = make(map[string]int64)
	ss.resourceFree = make(map[string]int64)
	ss.machines = make([]string, 0)
	ss.joblist = make(map[int64]*localJob)
	ss.jobqueue = make([]int64, 0)
	ss.jobrunning = make(map[int64]bool)
	ss.channewjob = make(chan *localJob, 5)
	ss.chanqueue = make(chan int, 10)
	// resource local
	local := "local"
	ss.machines = append(ss.machines, local)
	ss.resource[local] = 4
	ss.resourceFree[local] = 4
	ss.nodesize = 4
	go ss.submitD()
	go ss.queueD()
	go func() {
		for {
			time.Sleep(100 * time.Second)
			ss.chanqueue <- cmsgQSmatch

		}
	}()
	return ss
}

// daemon for submit jobs
func (ss *localScheduler) submitD() {
	for {
		job := <-ss.channewjob
		// just push the job to the queue .
		ss.lock.Lock()
		ss.idmax++
		id := ss.idmax
		job.JobID = fmt.Sprint(id)
		//ss.jobqueue[id] = job
		ss.jobqueue = append(ss.jobqueue, id)
		ss.joblist[id] = job
		ss.lock.Unlock()
		job.JobState = "Pending"
		go func() {
			job.chanmsg <- cmsgSTART
		}() // do not need to block
		ss.chanqueue <- cmsgRESALLOC
		//got a new job to run
		// add to queue
	}
}

// daemon for queue manager
// simple FIFO
func (ss *localScheduler) queueD() {
	for {
		msg := <-ss.chanqueue
		log.Printf("queue looop %v %v %v  \n", len(ss.jobqueue), len(ss.resourceFree), msg)

		if len(ss.jobqueue) == 0 {
			ss.queuestate = cmsgQSidle
			continue
		}
		// run job
		if len(ss.resourceFree) == 0 {
			ss.queuestate = cmsgQSpending
			continue
		}

		ss.queuestate = cmsgQSmatch

		switch msg {
		case cmsgRESALLOC:
			if ss.queuestate == cmsgQSpending {
				// no res free , go on pending
				continue
			}
		case cmsgRESFREE:
			if ss.queuestate == cmsgQSidle {
				// no job waiting for res .
				continue
			}
		default:
		}
		//ss.queuestate = cmsgQSpending
		// go to match job and res
		jobid := ss.jobqueue[0]
		// check job info
		job := ss.joblist[jobid]
		// can not deal nodes job now
		if job.Nodes > 1 {
			log.Printf("job %d failed as nodes  %d \n ", jobid, job.Nodes)
			job.JobState = "Failed"
			ss.lock.Lock()
			if len(ss.jobqueue) == 1 {
				ss.jobqueue = make([]int64, 0)
			} else {
				ss.jobqueue = ss.jobqueue[1:]
			}
			ss.lock.Unlock()
			continue
		}
		// do not deal overloap now
		if (job.Cores) > ss.nodesize {
			job.JobState = "Failed"
			log.Printf("job %d failed as cores  %d > %d \n ", jobid, job.Nodes, ss.nodesize)
			ss.lock.Lock()
			if len(ss.jobqueue) == 1 {
				ss.jobqueue = make([]int64, 0)
			} else {
				ss.jobqueue = ss.jobqueue[1:]
			}
			ss.lock.Unlock()
			continue
		}
		var machinename string
		var goalv int64 = job.Cores
		ss.lock.Lock()
		for k, v := range ss.resourceFree {
			if v >= goalv {
				machinename = k
				break
			}
		}
		if machinename == "" {
			ss.lock.Unlock()
			continue
		}
		job.JobState = "START"
		ss.jobqueue = ss.jobqueue[1:]

		ss.resourceFree[machinename] = ss.resourceFree[machinename] - goalv
		if ss.resourceFree[machinename] == 0 {
			delete(ss.resourceFree, machinename)
		}
		ss.lock.Unlock()
		go ss.runJob(jobid, machinename)
	}

}

func (ss *localScheduler) runJob(jobid int64, machine string) {
	job := ss.joblist[jobid]
	job.JobState = "START"
	// not really run now
	job.JobState = "Running"
	ss.lock.Lock()
	ss.jobrunning[jobid] = true
	ss.lock.Unlock()
	if machine == "local" {
		cmd := &Command{
			Name: "bash",
			Args: []string{},
			Dir:  job.Dir,
			//User : job.
		}
		cmd.Args = append(cmd.Args, job.JobFile)
		cmd.Args = append(cmd.Args, job.JobArgs...)
		cmdresult, err := modelsCommander.RunCommand(cmd)
		if err != nil {
			ss.finishJob(jobid, machine)
			return
		}
		if ss.returnFile == true {
			stdfilename := "log." + job.JobID
			errfilename := "err." + job.JobID
			if job.Dir != "" {
				stdfilename = path.Join(job.Dir, stdfilename)
				errfilename = path.Join(job.Dir, errfilename)
			}
			os.Create(stdfilename)
		} else {
			restb, _ := json.Marshal(cmdresult)
			job.HPCJob.Infos += string(restb)
		}
	}
	time.Sleep(20 * time.Second)
	//job.JobState = "CG"
	ss.finishJob(jobid, machine)
}

func (ss *localScheduler) finishJob(jobid int64, machine string) {

	job := ss.joblist[jobid]
	job.JobState = "Finish"
	ss.lock.Lock()
	delete(ss.jobrunning, jobid)
	ss.lock.Unlock()
	if machine == "" {
		return
	}
	ss.lock.Lock()
	v, ok := ss.resourceFree[machine]
	if !ok {
		ss.resourceFree[machine] = job.Cores
	} else {
		ss.resourceFree[machine] = v + job.Cores
	}
	ss.lock.Unlock()
	ss.chanqueue <- cmsgRESFREE
}

func (ss *localScheduler) Submit(job *HPCJob) (jobid string, err error) {
	ljob := &localJob{HPCJob: *job}
	ljob.chanmsg = make(chan int, 2)
	ss.channewjob <- ljob
	<-ljob.chanmsg

	return ljob.JobID, nil

}

////////////////////////

func (ss *localScheduler) Info() ([]*HPCResource, error) {
	resl := make([]*HPCResource, 2)
	resl[0] = &HPCResource{Partition: "example", Nodes: 1, State: "idle"}
	resl[1] = &HPCResource{Partition: "example", Nodes: 1, State: "down"}
	return resl, nil
}

func (ss *localScheduler) Queue() ([]*HPCJob, error) {
	jobl := make([]*HPCJob, 0)
	for _, jobid := range ss.jobqueue {
		jobl = append(jobl, &ss.joblist[jobid].HPCJob)
	}
	for jobid, _ := range ss.jobrunning {
		jobl = append(jobl, &ss.joblist[jobid].HPCJob)
	}

	return jobl, nil
}

func (ss *localScheduler) JobInfo(jobid string) ([]*HPCJob, error) {
	exampljobss := make([]*HPCJob, 1)
	idint, err := strconv.ParseInt(jobid, 0, 64)
	if err != nil {
		return exampljobss, err
	}
	job := ss.joblist[idint]
	exampljobss = append(exampljobss, &job.HPCJob)
	return exampljobss, err
}
func (ss *localScheduler) Delete(jobid string) error {
	return nil
}
