package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

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
	Info() ([]HPCResource, error)    // such as sinfo
	Queue() ([]HPCJob, error)        // such as squeue
	Submit(*HPCJob) (string, error)  // such as sbatch
	JobInfo(string) ([]HPCJob, error) // such as sacct -j  ?
	Delete(string) error             // such as scancel
	Name() string                    // Got the name of scheduler
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

func ResourceInfo() ([]HPCResource, error) {
	return modelsScheduler.Info()
}

func JobQueue() ([]HPCJob, error) {
	return modelsScheduler.Queue()
}

func (j *HPCJob) Submit() (string, error) {
	if j.myScheduler != nil {
		return j.myScheduler.Submit(j)
	}
	return modelsScheduler.Submit(j)
}

func JobInfo(jobid string) ([]HPCJob, error) {
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
		exapleSD := &shellScheduler{}
		exapleSD.shelldir = josshelldir
		modelsScheduler = exapleSD
	case "example":
	case "test":
		exapleSD := &exampleScheduler{}
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

func (*exampleScheduler) Info() ([]HPCResource, error) {
	resl := make([]HPCResource, 2)
	resl[0] = HPCResource{Partition: "example", Nodes: 1, State: "idle"}
	resl[1] = HPCResource{Partition: "example", Nodes: 1, State: "down"}
	return resl, nil
}

func (*exampleScheduler) Queue() ([]HPCJob, error) {
	jobl := make([]HPCJob, 2)
	jobl[0] = HPCJob{Name: "job1"}
	jobl[1] = HPCJob{Name: "job2"}
	return jobl, nil
}

func (*exampleScheduler) Submit(*HPCJob) (jobid string, err error) {
	return "job1", nil
}
func (*exampleScheduler) JobInfo(string) ([]HPCJob, error) {
        //#var exampljobss []HPCJob{ HPCJob{ Name: "job1"}}
        exampljobss := make([]HPCJob,1)
	exampljobss[0] = HPCJob{ Name: "job1"}
	return exampljobss, nil
}
func (*exampleScheduler) Delete(string) error {
	return nil
}
func (*exampleScheduler) Name() string {
	return "exampleScheduler"
}

type shellScheduler struct {
	shelldir string
}

func (ss *shellScheduler) run(filename string, mydata interface{}) error {
	//log.Println("start " + filename)
	shellcmd := ss.shelldir + "/" + filename //path.Join(ss.shelldir, filename)
	result, err := ShellRun(shellcmd)
        log.Printf("Run shell %s : %d \n %s \n",shellcmd , result.Retcode , result.Output)
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

	// if retjson {
	// 	json.Unmarshal([]byte(result.Output), mydata)
	// } else {
	// 	mytype = result.Output
	// }
	return nil
}

func (ss *shellScheduler) Info() ([]HPCResource, error) {
	log.Println("start info()")
	resl := make([]HPCResource, 0)
	/*
		infoshell := path.Join(ss.shelldir, "info")
		result, err := ShellRun(infoshell)
		if err != nil {
			return resl, err
		}
		if result.Retcode != 0 {
			return resl, errors.New("Error " + string(result.Retcode) + " : " + result.Error)
		}
		log.Printf("result : %#v\n", result)
		json.Unmarshal([]byte(result.Output), &resl)
		log.Printf("result : %#v\n", resl)
		return resl, err
	*/
	err := ss.run("info", &resl)
	return resl, err
}

func (ss *shellScheduler) Queue() ([]HPCJob, error) {
	jobl := make([]HPCJob, 0)
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
		/*
			fmt.Printf("%d : %s %s = %v \n", i,
				s.Type().Field(i).Name,
				f.Type(),
				f.Interface())
		*/
		//jsonb, _ := json.Marshal(f.Interface())
		//fmt.Println(s.Type().Field(i).Name + "=" + string(jsonb))
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
func (ss *shellScheduler) JobInfo(jobid string) ([]HPCJob,  error) {
	jobss := make([]HPCJob, 0)
	err := ss.run("jobinfo "+jobid, jobss)
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
