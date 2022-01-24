package common

import (
	"github.com/gin-gonic/gin"
)

type Job struct {
	Name  string  `json:"name"`
	Command string  `json:"command"`
	CronExpr string  `json:"expr"`
}

type  Response struct {
	Error   int `json:"error"`
	Msg string `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(error int,msg string,Data interface{},c *gin.Context) {

	c.JSON(200,gin.H{
		"msg": msg,
		"data":Data,
	})
}