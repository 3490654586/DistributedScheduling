package router

import (
    "cron/master"
    v1 "cron/master/api/v1"
    "github.com/gin-gonic/gin"
)

func InitServer()  {
    r := gin.Default()
    r.POST("/cron/job", v1.HandleJobSave)
    r.POST("/job/delel",v1.HandleDelelJob)
    r.GET("/job/quer",v1.HandeleFindjob)
    r.POST("/job/kill",v1.HandelKilljob)
    r.GET("/job/log",v1.GetLog)
    con_fig := master.GetConfig()
    r.Run(":" + con_fig.ApiPost + "")
}