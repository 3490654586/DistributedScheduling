package v1

import (
	"cron/common"
	"cron/master"
	"fmt"
	"github.com/gin-gonic/gin"
)

//添加任务
func HandleJobSave(context *gin.Context)  {
	fmt.Println("ddad")
	var JobReq common.Job
	err := context.ShouldBind(&JobReq)

	if err != nil {
		fmt.Println("解析请求参数失败",err)
		return
	}

	//调用接口保存到ETCD
	reqJob,err := master.JobMarStruct.SaveJob(&JobReq)
	if err != nil {
		context.JSON(200,gin.H{
			"error":err,
			"msg":"保存任务失败",
		})
	}

	context.JSON(200,gin.H{
		"error":err,
		"msg":"保存任务成功",
		"data":reqJob,
	})
}

//删除任务
func HandleDelelJob(contetx *gin.Context)  {
	name := contetx.PostForm("name")

	reqJob,err := master.JobMarStruct.DelelJob(name)
	if err != nil {
		contetx.JSON(200,gin.H{
			"error":err,
			"msg":"删除任务失败",
		})
	}

	contetx.JSON(200,gin.H{
		"error":err,
		"msg":"删除任务成功",
		"data":reqJob,
	})
}

//查询任务列表
func HandeleFindjob(contetx *gin.Context)  {
      res,err :=   master.JobMarStruct.FindJob()
	if err != nil {
		contetx.JSON(200,gin.H{
			"error":err,
			"msg":"查询任务失败",
		})
		return
	}

	contetx.JSON(200,gin.H{
		"error":err,
		"msg":"查询任务成功",
		"data":res,
	})
}

//杀死任务
func HandelKilljob(contetx *gin.Context)  {
	name := contetx.PostForm("name")
	err := master.JobMarStruct.KillJob(name)
	if err != nil {
		contetx.JSON(200,gin.H{
			"error":err,
			"msg":"杀死任务失败",
		})
		return
	}
	contetx.JSON(200,gin.H{
		"error":err,
		"msg":"杀死任务成功",
	})
}

func GetLog(ctx *gin.Context)  {
	res := master.G_log.GetJson()
	ctx.JSON(200,gin.H{
		"error":"",
		"msg":"查询任务成功",
		"data":res,
	})
}