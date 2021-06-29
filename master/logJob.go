package master

import (
	"context"
	"cron/common"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type SetLog struct {
	cli *mongo.Client
	logCollectuib *mongo.Collection
	logChan chan *common.JobLog
}

var G_log *SetLog


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
	}

}


func (setlog *SetLog)GetJson() []*common.JobLog {
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

	return results

}