package models

import (

	//ext "gehpci/models/extend/job"
	"log"

	"github.com/astaxie/beego"
)

type HPCResource struct {
	Partition   string `json:"partition"`
	State       string `json:"state"`
	Nodes       int64  `json:"nodes"`
	Infos       string `json:"Infos"`
	myScheduler HPCScheduler
}

type HPCScheduler interface {
	Info() ([]*HPCResource, error)     // such as sinfo
	Queue() ([]*HPCJob, error)         // such as squeue
	Submit(*HPCJob) (string, error)    // such as sbatch
	JobInfo(string) ([]*HPCJob, error) // such as sacct -j  ?
	Delete(string) error               // such as scancel
	Name() string                      // Got the name of scheduler
}

var modelsScheduler HPCScheduler

type HPCJob struct {
	Dir       string `json:"dir"`       //"job_wdir" ,　 工作路径
	Name      string `json:"name"`      //"job_name" , 　作业名称
	TimeLimit int64  `json:"timelimit"` // "time_limit" ,　执行时间限制（min） time.Duration int64
	Partition string `json:"partition"` //"partition" ,　　作业分区
	Cores     int64  `json:"cores"`     //"scale_cores" , 　需要多少核
	//"scale_memGB" , 　需要多少内存
	Nodes int64 `json:"nodes"` //"scale_Nodes" , 需要多少节点
	// ? ? ? Scheduler string `json:"scheduler"` //"jobfile_args"，　作业脚本附加参数 　 一个提交例子如下：
	Args        string `json:"args"` // args for scheduler
	JobFile     string `json:"jobfile"`
	JobArgs     string `json:"jobargs"`
	JobID       string `json:"jobid"`
	JobState    string `json:"jobstate"`
	Infos       string `json:"infos"` // may be used for extension
	myScheduler HPCScheduler
}

func ResourceInfo() ([]*HPCResource, error) {
	return modelsScheduler.Info()
}

func JobQueue() ([]*HPCJob, error) {
	return modelsScheduler.Queue()
}

func (j *HPCJob) Submit() (string, error) {
	if j.myScheduler != nil {
		return j.myScheduler.Submit(j)
	}
	return modelsScheduler.Submit(j)
}

func JobInfo(jobid string) ([]*HPCJob, error) {
	return modelsScheduler.JobInfo(jobid)
}

func DeleteJob(jobid string) error {
	return modelsScheduler.Delete(jobid)
}

func (j *HPCJob) Scheduler() string {
	return j.myScheduler.Name()
}

func init() {
	jobmode := beego.AppConfig.String("jobmode")
	//var exapleSD HPCScheduler
	switch jobmode {
	case "shell":
		josshelldir := beego.AppConfig.String("shell::job")
		exapleSD := &shellScheduler{} //&shellScheduler{}
		exapleSD.shelldir = josshelldir
		modelsScheduler = exapleSD
	case "example":
	case "test":
		exapleSD := &exampleScheduler{}
		modelsScheduler = exapleSD
	case "local":
		exapleSD := &localScheduler{}
		exapleSD.returnFile = true
		modelsScheduler = exapleSD
	default:
		exapleSD := &shellScheduler{}
		exapleSD.shelldir = "./shells/exampleScheduler/"
		modelsScheduler = exapleSD
	}
	//modelsScheduler = exapleSD
	if modelsScheduler == nil {
		log.Fatal("Error : Scheduler cannot create!")
	}
	log.Printf("init ok : %#v \n", modelsCommander)
}

type exampleScheduler struct {
}

func (*exampleScheduler) Info() ([]*HPCResource, error) {
	resl := make([]*HPCResource, 2)
	resl[0] = &HPCResource{Partition: "example", Nodes: 1, State: "idle"}
	resl[1] = &HPCResource{Partition: "example", Nodes: 1, State: "down"}
	return resl, nil
}

func (*exampleScheduler) Queue() ([]*HPCJob, error) {
	jobl := make([]*HPCJob, 2)
	jobl[0] = &HPCJob{Name: "job1"}
	jobl[1] = &HPCJob{Name: "job2"}
	return jobl, nil
}

func (*exampleScheduler) Submit(*HPCJob) (jobid string, err error) {
	return "job1", nil
}
func (*exampleScheduler) JobInfo(string) ([]*HPCJob, error) {
	//#var exampljobss []HPCJob{ HPCJob{ Name: "job1"}}
	exampljobss := make([]*HPCJob, 1)
	exampljobss[0] = &HPCJob{Name: "job1"}
	return exampljobss, nil
}
func (*exampleScheduler) Delete(string) error {
	return nil
}
func (*exampleScheduler) Name() string {
	return "exampleScheduler"
}
