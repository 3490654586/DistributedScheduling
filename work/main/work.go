package main

import (
	"cron/work"
	"flag"
	"fmt"
	"runtime"
	"time"
)


var(
	PathName string
)

func InitArgs()  {
	flag.StringVar(&PathName,"config","./work.json","work.json")
	flag.Parse()
}


func InitCpu()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main()  {
	//初始化cpu
	InitCpu()

	InitArgs()
	fmt.Println(PathName)

	work.InitConfig(PathName)

	work.InitLog()
    //任务执行器
    work.InitImplement()
	//任务调度器
	work.InitDispatch()
	work.InitJoMar()

	for {
		time.Sleep(1 * time.Second)
	}
	//初始化api接口服务
	//router.InitServer()
}