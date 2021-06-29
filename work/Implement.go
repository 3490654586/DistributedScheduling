package work

import (
	"context"
	"cron/common"
	"fmt"
	"os/exec"
	"time"
)

type Implement struct {
	Jobimplement *common.Jobimplement
}


//执行一个任务
func (implement *Implement)ExitJob(jobSchedulePlan *common.Jobimplement)  {
      go func() {
      	 //执行shell命令
		  jobExiTres := &common.JobExiTres{
			  ExitInfo: jobSchedulePlan,
			  Oputput:  make([]byte,0),
		  }

		  //初始化锁
		 jobLock := JobMarStruct.CreatLock(jobSchedulePlan.Job.JobName)

		  jobExiTres.STime = time.Now()
		  err :=jobLock.TryLock()
		  defer jobLock.UnLcok()
		  if err != nil {
			  fmt.Println("上锁失败",err)
			  return
		  }else {
		  	 jobExiTres.STime = time.Now()
			  cmd := exec.CommandContext(context.Background(),"/bin/bash","-c",jobSchedulePlan.Job.Command)
			  output,err :=cmd.CombinedOutput()

			  jobExiTres.ETime = time.Now()
			  jobExiTres.Oputput = output
			  jobExiTres.Err = err
		  }

		  //执行完成,返回r给调度器,调度器从任务执行表中删除执行记录
		  G_Dispatch.PushReschan(jobExiTres)
      }()

}

var G_implement *Implement

func InitImplement()  {
	G_implement = &Implement{}
}