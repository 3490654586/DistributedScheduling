package common

import (
	"encoding/json"
	"fmt"
	"github.com/gorhill/cronexpr"
	_ "google.golang.org/genproto/googleapis/cloud/ml/v1"
	"strings"
	"time"
)

//日志批量
type LogARR struct {
	Logs []interface{}
}

//任务
type Job struct {
	JobName string `json:"job_name" form:"job_name"` //任务名字
    Command string `json:"command"  form:"command"` //shell命令
    CronExpr string `json:"cron_expr" form:"cron_expr"` //cron表达式
}


// 任务调度计划
type JobSchedulePlan struct {
	Job *Job	// 要调度的任务信息
	Expr *cronexpr.Expression	// 解析好的cronexpr表达式
	NextTime time.Time	// 下次调度时间
}

//任务执行状态
type Jobimplement struct {
	Job *Job
	PlanTime time.Time
	ReaTime time.Time
}

//任务执行结果
type JobExiTres struct {
	ExitInfo *Jobimplement
	Oputput []byte
	Err error
	STime time.Time
	ETime time.Time
}

// 变化事件
type JobEvent struct {
	EventType int //  SAVE, DELETE
	Job *Job
}

type JobLog struct {
	JobName string`bson:"Jobname"`
	Command string `bson:"command"`
	Output string `bson:"ouput"`
	Err string `bson:"err"`
	Startime int64 `bson:"statime"`
	EndTime int64 `bson:"endtime"`
}

func ExtractJobName(jobkey string)string{
	  return strings.TrimPrefix(jobkey,JOB_SAVE_DIR)
}

// 任务变化事件有2种：1）更新任务 2）删除任务
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job: job,
	}
}

//构造执行计划
func BuildDispatch(job *Job) (*JobSchedulePlan,error){
	expr,err := cronexpr.Parse(job.CronExpr)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	jobSchedulePlan :=  &JobSchedulePlan{
		 Job:      job,
		 Expr:     expr,
		 NextTime: expr.Next(time.Now()),
	 }

	 fmt.Println("任务调度计划表",jobSchedulePlan)
	 return jobSchedulePlan,nil
}

//构造执行状态
func Buildimplement(jobDispatch *JobSchedulePlan) *Jobimplement {
	 jobimplement := &Jobimplement{
		Job:      jobDispatch.Job,
		PlanTime: jobDispatch.NextTime,
		ReaTime:  time.Now(),
	}
	return jobimplement
}
// 反序列化Job
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}


