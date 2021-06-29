package work

import (
	"cron/common"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	wath clientv3.Watcher
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
		wath := clientv3.NewWatcher(cli)
		JobMarStruct = &JobMgr{
			client: cli,
			kv:     cliKv,
			lease:  cliLease,
			wath:   wath,
		}
		JobMarStruct.watcJobs()
		
}


/**
 监听任务变化
 1.查询任务列别
 2.从该revision向后监听变化事件
 */

func (jomber *JobMgr)watcJobs()error  {
	//查询任务列表
	getres ,err := jomber.kv.Get(context.Background(),common.JOB_SAVE_DIR,clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return err
	}
      var oldjob *common.Job
	var jobEvents *common.JobEvent

	//获得所有job
	for _,value := range getres.Kvs{
		job,_:=common.UnpackJob(value.Value)
		jobEvents = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)

		G_Dispatch.PushDispatchCh(jobEvents)
		fmt.Println("old=",oldjob)
	}
     //fmt.Println("old=",oldjob)
	//监听
	go func() {
		//fmt.Println("dadawd")
		//获得监听的版本号
		Revision := getres.Header.Revision + 1

		wathch := jomber.wath.Watch(context.Background(),common.JOB_SAVE_DIR,clientv3.WithRev(Revision),clientv3.WithPrefix())

        //fmt.Println("watch",wathch)
		for watchres := range wathch{
			//fmt.Println("watchres",watchres)
			for _,event := range watchres.Events{
				fmt.Println("Type",event.Type)
				switch event.Type {
				case mvccpb.PUT:
					jobEvents = common.BuildJobEvent(common.JOB_EVENT_SAVE,oldjob)
					//G_Dispatch.PushDispatchCh(jobEvents)
				case mvccpb.DELETE:
					jobname := common.ExtractJobName(string(event.Kv.Key))

					oldjob = &common.Job{
						JobName:  jobname,
					}
					//构建一个删除事件
					jobEvents = common.BuildJobEvent(common.JOB_EVENT_DELETE,oldjob)
				}
				fmt.Println("我push的任务",jobEvents)
				G_Dispatch.PushDispatchCh(jobEvents)
			}
		}

	}()

	return nil
}


//创建任务锁

func (jobmar *JobMgr)CreatLock(name string) *JobLock{
	jobLock := InitJobLock(name,jobmar.client,jobmar.lease)
	return jobLock
}