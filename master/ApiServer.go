package master

import (
	"crontab/common"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandleJobSave 任务保存接口
// post  job ={"name":"job1","command":"echo hello world","expr":"****}
func HandleJobSave(c *gin.Context) {
	//任务保存到etcd
	job := c.PostForm("job")
	var resJob common.Job

	err := json.Unmarshal([]byte(job), &resJob)
	if err != nil {
		fmt.Println("序列化错误")
		return
	}

	oldJob, err := G_jobMgr.SaveJob(&resJob)
	if err != nil {
		c.JSON(200, gin.H{
			"error": 1,
			"msg":   "保存失败",
			"data":  nil,
		})
	}
	c.JSON(200, gin.H{
		"error": 0,
		"msg":   "保存成功",
		"data":  oldJob,
	})
}

func HandleJobDelete(c *gin.Context) {
	jobName := c.Param("jobname")
	oldJob, err := G_jobMgr.DeleteJob(jobName)
	if err != nil {
		c.JSON(200, gin.H{
			"error": 1,
			"msg":   err.Error(),
			"data":  nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"error": 0,
		"msg":   "删除成功",
		"data":  oldJob,
	})
}

func HandleJobList(c *gin.Context) {
	jobList, err := G_jobMgr.GetJobList()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": jobList,
	})
}

func HandleJobKill(c *gin.Context) {
	killJobName := c.Param("name")

	err := G_jobMgr.KillJob(killJobName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
			"msg":   "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "删除成功",
	})
}
