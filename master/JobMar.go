package master

import (
	"context"
	"cron/common"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}
var JobMarStruct *JobMgr

func InitJoMar(){
	    etcdconfig := GetConfig()

       config := clientv3.Config{Endpoints:[]string{etcdconfig.EtcdEndpoints[0]},DialTimeout:time.Duration(etcdconfig.EtcdDialTimeout)* time.Millisecond}
	     cli,err := clientv3.New(config)

		if err != nil {
			fmt.Println(err)
			return
		}

		cliKv := clientv3.NewKV(cli)

		cliLease := clientv3.NewLease(cli)
		JobMarStruct = &JobMgr{
			client: cli,
			kv:     cliKv,
			lease:  cliLease,
		}
		
}

//保存任务
func (jobMar *JobMgr)SaveJob(job *common.Job)(*common.Job,error) {
        //保存key
        jobKey := common.JOB_SAVE_DIR+job.JobName
         fmt.Println("keuy",jobKey)
        jobValue,err  :=json.Marshal(job)
		if err != nil {
			log.Panicln(err)
			return nil,err
		}

		//保存到etcd
		PutRes,err :=jobMar.kv.Put(context.Background(),jobKey,string(jobValue),clientv3.WithPrevKV())
		if err != nil {
		log.Println(err)
		return nil,err
	}

     var PrevkeJob common.Job
    //不是创建是更新
	if PutRes.PrevKv !=nil{
		json.Unmarshal(PutRes.PrevKv.Value,&PrevkeJob)
		return &PrevkeJob,nil
	}

	return nil,nil
}

//删除任务
func (jobMgr *JobMgr)DelelJob(name string)(*common.Job,error)  {
	   jobkey :=common.JOB_SAVE_DIR + name
       var oldJob common.Job

	   delres,err :=  jobMgr.kv.Delete(context.Background(),jobkey,clientv3.WithPrevKV())
	  if err != nil {
	  	return nil,err
	}

	 if delres.PrevKvs != nil{
	 	json.Unmarshal(delres.PrevKvs[0].Value,&oldJob)
	 	return &oldJob,nil
	 }

    return nil,nil
}

//查询任务列表
func (jobMar *JobMgr)FindJob()([]*common.Job,error){
        jobkey := common.JOB_SAVE_DIR

      getres,err := jobMar.kv.Get(context.Background(),jobkey,clientv3.WithPrefix())
		if err != nil {
			return nil,err
		}

		arrjob :=make([]*common.Job,0)

		for _,value := range getres.Kvs{
			jobval := common.Job{}
			err = json.Unmarshal(value.Value,&jobval)
			if err != nil {
				err = nil
				continue
			}
			arrjob = append(arrjob,&jobval)
		}
     return arrjob,nil
}

//杀死任务
func (jobMar *JobMgr)KillJob(name string)error{
         jobkey := common.JOB_KILL_DIR + name

        leaseRes ,err := jobMar.lease.Grant(context.Background(),1)
		if err != nil {
			log.Println(err)
			return err
		}

		leaseResId := leaseRes.ID

		 _,err = jobMar.kv.Put(context.Background(),jobkey,"",clientv3.WithLease(leaseResId))

		if err != nil {
			return err
		}

	return nil
}