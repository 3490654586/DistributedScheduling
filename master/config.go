package master

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ApiPost string `json:"api_post"`
	EtcdEndpoints []string `json:"etcd_endpoints"`
	EtcdDialTimeout int `json:"etcd_dial_timeout"`
}
var Con_fig Config

func InitConfig(path string)  {


	content,err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("打开配置文件失败",err)
		return
	}

	err = json.Unmarshal(content,&Con_fig)
	if err != nil {
		log.Println("反序列化错误",err)
	}
}

func GetConfig()*Config{
	return &Con_fig
}
