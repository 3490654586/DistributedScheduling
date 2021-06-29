package work

import (
	"context"
	"cron/common"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type SetLog struct {
	cli *mongo.Client
	logCollectuib *mongo.Collection
	logChan chan *common.JobLog
	autoCommit chan*common.LogARR
}

var G_log *SetLog

func (setlog *SetLog)SendLog(logarr *common.LogARR)  {
	setlog.logCollectuib.InsertMany(context.Background(),logarr.Logs)
}


func (setlog *SetLog)writlog()  {
	var logarr *common.LogARR
	var commTime *time.Timer
	for{
		select {
		case log := <- setlog.logChan:
			if logarr == nil{
				logarr =&common.LogARR{}
				commTime =time.AfterFunc(time.Duration(1000) *time.Millisecond, func(logarrs *common.LogARR)func(){
					return func() {
                       setlog.autoCommit <- logarrs
					}
				}(logarr),
				)
			}

			logarr.Logs = append(logarr.Logs,log)

			if len(logarr.Logs) >= 100 {
				setlog.SendLog(logarr)
				logarr = nil
				commTime.Stop()
			}
		case timeLogarr :=<- setlog.autoCommit:
			if timeLogarr != logarr{
				continue
			}

			setlog.SendLog(timeLogarr)
			timeLogarr = nil
	  }
	}
}

func InitLog()  {
	url := options.Client().ApplyURI("mongodb://106.52.198.127:27017")
	cli,err := mongo.Connect(context.Background(),url)
	if err != nil {
		fmt.Println(err)
		return
	}

	G_log = &SetLog{
		cli:           cli,
		logCollectuib: cli.Database("cron").Collection("log"),
		logChan:       make(chan *common.JobLog,1000),
		autoCommit:    make(chan *common.LogARR,1000),
	}
	go G_log.writlog()
	G_log.GetJson()
}

func (strlog *SetLog)Append(log *common.JobLog)  {
	select {
	case strlog.logChan <- log:
	default:

	 }
}

func (setlog *SetLog)GetJson()  {
	fmt.Println("我进来查询数据")
 var results []*common.JobLog
	cur, err := setlog.logCollectuib.Find(context.Background(),bson.D{})
	for cur.Next(context.TODO()) {
		// 解码：数据转换
		var u common.JobLog
		err = cur.Decode(&u)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &u)
	}
	// 关闭游标
	cur.Close(context.TODO())

	var i int
	arrLen := len(results)
	for i = 0; i < arrLen; i++ {
		tmp := results[i]
		fmt.Println(i, tmp.JobName, tmp.Command)
	}
}