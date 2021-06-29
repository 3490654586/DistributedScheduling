package work

import (
	"context"
	"cron/common"
	"fmt"
	"go.etcd.io/etcd/clientv3"
)

type JobLock struct {
	Kv clientv3.KV
	Lease clientv3.Lease

	jobName string
    canceFunc context.CancelFunc
	leaseID clientv3.LeaseID
	isLock bool
}



func InitJobLock(jobName string,kv clientv3.KV,lease clientv3.Lease) *JobLock {
	jobLock := &JobLock{
		Kv:      kv,
		Lease:   lease,
		jobName: jobName,
	}

	return jobLock
}



func (jobLock*JobLock)TryLock()error  {
	// 创建租约(5秒)
   leaseres,err := jobLock.Lease.Grant(context.Background(),5)
	if err != nil {
		fmt.Println("创建租约失败",err)
		return nil
	}
	//租约id
	leaseId :=leaseres.ID

	//创建可关闭的context上下文
	ctx,chanfunl := context.WithCancel(context.Background())
	defer chanfunl()
	defer jobLock.Lease.Revoke(context.Background(),leaseId)

	//自动续租
	leaseKeppRes,err := jobLock.Lease.KeepAlive(ctx,leaseId)
	if err != nil {
		return err
	}
	//处理续租应答的协程
	go func() {
		for{
			select {
			case leaseKeppCH := <- leaseKeppRes:
				if leaseKeppCH == nil{
					goto END
				}
			}
		}
		END:
	}()

	//创建事务
	txn := jobLock.Kv.Txn(context.Background())
	txnPath := common.JOB_LOCK_DIR + jobLock.jobName

	//抢锁
      txn.If(clientv3.Compare(clientv3.CreateRevision(txnPath),"=",0)).
		Then(clientv3.OpPut(txnPath,"",clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(txnPath))

	//提交事务
	txnres,err := txn.Commit()
	if err != nil {
		fmt.Println("提交事务失败",err)
		return err
	}
	//成功返回, 失败释放租约
      if !txnres.Succeeded{
      	return nil
	  }

	//抢锁成功
	jobLock.leaseID = leaseId
	jobLock.canceFunc = chanfunl
	jobLock.isLock = true
	return nil
}

func (jobLock *JobLock)UnLcok()  {
     if jobLock.isLock {
     	jobLock.canceFunc()
     	jobLock.Lease.Revoke(context.Background(),jobLock.leaseID)
	 }
}