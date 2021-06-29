package work

import (
	"cron/common"
	"fmt"
	"time"
)

type Dispatch struct {
    jobEventChan  chan *common.JobEvent  //etcd任务执行队列
    jobTable map[string]*common.JobSchedulePlan //任务调度计划表
    jobImplementTable map[string]*common.Jobimplement //任务执行表
    jobResChan chan *common.JobExiTres
}


var G_Dispatch *Dispatch


func InitDispatch(){
	G_Dispatch = &Dispatch{
		jobEventChan:make(chan *common.JobEvent,1000),
		jobTable:make(map[string]*common.JobSchedulePlan),
		jobImplementTable:make(map[string]*common.Jobimplement),
		jobResChan:make(chan *common.JobExiTres,1024),
			}
	go G_Dispatch.dispatchLoop()
}

//处理任务事件
func (dispatch*Dispatch)HandelJobEvent(jobenvet *common.JobEvent) {
	switch jobenvet.EventType {
	case common.JOB_EVENT_SAVE:
		DispatchPlan,err := common.BuildDispatch(jobenvet.Job)
		if err != nil {
			fmt.Println("解析Cron表达式失败",err)
			return
		}
		dispatch.jobTable[jobenvet.Job.JobName] = DispatchPlan
	case common.JOB_EVENT_DELETE:
		value,jobExisted := dispatch.jobTable[jobenvet.Job.JobName]
		if jobExisted {
			delete(dispatch.jobTable,value.Job.JobName)
		}
	}
}

//调度携程
func (dispatch*Dispatch)dispatchLoop(){
    //初始化调度时间为一秒
	forDispatch := dispatch.ForDispatch()

	dispatchTimer := time.NewTimer(forDispatch)

	for{
		select{
		case jobEventRes := <-dispatch.jobEventChan:
			//对列表做同步
			dispatch.HandelJobEvent(jobEventRes)
		case <-dispatchTimer.C:
		case jobres := <- dispatch.jobResChan:
			fmt.Println("任务结果队列",jobres)
			 dispatch.handelJobres(jobres)
		}
		forDispatch = dispatch.ForDispatch()
		//重新调度间隔
		dispatchTimer.Reset(forDispatch)
	}
}

func (dispatch*Dispatch)PushDispatchCh(jobevent *common.JobEvent)  {
	dispatch.jobEventChan <- jobevent
}

//处理任务结果
func (dispatch*Dispatch)handelJobres(res *common.JobExiTres)  {
	//从任务执行表删除执行完成任务
	  delete(dispatch.jobImplementTable,res.ExitInfo.Job.JobName)

	  //生成执行日志
	  if res.Err != common.ERR_LOCK{
	  	joblog := &common.JobLog{
			 JobName:  res.ExitInfo.Job.JobName,
			 Command:  res.ExitInfo.Job.Command,
			 Output:   string(res.Oputput),
			 Err:      res.Err.Error(),
			 Startime: res.STime.UnixNano()/1000/1000,
			 EndTime:  res.ETime.UnixNano()/1000/1000,
		 }
		 G_log.Append(joblog)
	  }
	  fmt.Println("任务执行完成",res.ExitInfo.Job.JobName)
}

//计算任务调度状态
func (dispatch*Dispatch)ForDispatch() time.Duration {

	//任务表为空
	if len(dispatch.jobTable) == 0{
       timeSlepp := 1 * time.Second
       return timeSlepp
	}

       now := time.Now()
       var nearTIme *time.Time

       for _,jobPlan := range dispatch.jobTable{
       	if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now){
            dispatch.ImplementJob(jobPlan)
       		//尝试执行任务更新下次调度时间
       		jobPlan.NextTime = jobPlan.Expr.Next(now)
		}

		 //要过期的任务时间
		 if nearTIme == nil || jobPlan.NextTime.Before(*nearTIme){
		 	nearTIme = &jobPlan.NextTime
		 }
	   }
   //下次调度时间
	dispatchAfter := (*nearTIme).Sub(now)
	return dispatchAfter
}


//执行任务
func (dispatch*Dispatch)ImplementJob(jobplan *common.JobSchedulePlan)  {
	 //从任务执行表中查看有没有任务
	 _,ok :=dispatch.jobImplementTable[jobplan.Job.JobName]
	 if ok {
	 	//fmt.Println("程序还在执行,",jobplan.Job.JobName)
		 return
	 }

	  //如果不存在,构建执行状态信息
	  jobExitInfo := common.Buildimplement(jobplan)

	  //保存执行状态
	  dispatch.jobImplementTable[jobplan.Job.JobName] = jobExitInfo

	  //执行任务
	  fmt.Println("执行任务",jobExitInfo.Job.JobName)
	  G_implement.ExitJob(jobExitInfo)
}


//任务执行结果
func (dispatch *Dispatch)PushReschan(jobres *common.JobExiTres)  {
	dispatch.jobResChan <- jobres
}

