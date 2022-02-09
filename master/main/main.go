package main

import (
	"crontab/master"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	r.LoadHTMLGlob("./webroot/*")
	r.Handle("GET", "/", func(context *gin.Context) {
		// 返回HTML文件，响应状态码200，html文件名为index.html，模板参数为nil
		context.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("job/save/", master.HandleJobSave)
	r.GET("job/delete/:jobname", master.HandleJobDelete)
	r.GET("job/list", master.HandleJobList)
	r.GET("job/kill/:name", master.HandleJobKill)
	r.Run(fmt.Sprintf(":%d", master.G_config.ApiPort))
}
