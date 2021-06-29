package main

import (
	"cron/master"
	"cron/master/router"
	"flag"
	"fmt"
	"runtime"
)


var(
	PathName string
)

func InitArgs()  {
	flag.StringVar(&PathName,"config","./config.json","指定master节点的配置文件")
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
   master.InitLog()
	master.InitConfig(PathName)

	master.InitJoMar()
	//初始化api接口服务
	router.InitServer()
}