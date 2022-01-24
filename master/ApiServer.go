package master

import (
	"crontab/common"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)


// HandleJobSave 任务保存接口
// post  job ={"name":"job1","command":"echo hello world","expr":"****}
func HandleJobSave(c *gin.Context)  {
	//任务保存到etcd
	job :=c.PostForm("job")
	var resJob common.Job

	err := json.Unmarshal([]byte(job), &resJob)
	if err != nil {
		fmt.Println("序列化错误")
		return
	}

	oldJob, err := G_jobMgr.SaveJob(&resJob)
	if err != nil {
		c.JSON(200,gin.H{
			"msg": "保存失败",
			"data":nil,
		})
	}
	c.JSON(200,gin.H{
		"msg": "保存成功",
		"data":oldJob,
	})


}

