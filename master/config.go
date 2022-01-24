package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)
var G_config  Config

type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `Json:"apiReadTimeout"`
	ApiWriteTimeout  int `json:"apiWriteTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDiaTimeout"`
}

func InitConfig(fileName string) error {
	content ,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		fmt.Println(err)
		return err
	}
	//json序列化
	err = json.Unmarshal(content, &G_config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}