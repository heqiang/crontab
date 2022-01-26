package main

import (
	"crontab/master"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

var confFile string

func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "配置文件")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	initEnv()
	initArgs()
	err := master.InitConfig(confFile)
	if err != nil {
		return
	}
	err = master.InitJobMag(master.G_config)
	if err != nil {
		return
	}
	r := gin.Default()
	r.POST("job/save/", master.HandleJobSave)
	r.GET("job/delete/:jobname", master.HandleJobDelete)
	r.GET("job/list", master.HandleJobList)
	r.GET("job/kill/:name", master.HandleJobKill)
	r.Run(fmt.Sprintf(":%d", master.G_config.ApiPort))
}
