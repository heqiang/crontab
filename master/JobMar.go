package master

import (
	"context"
	"crontab/common"
	"encoding/json"
	"fmt"
	client3 "go.etcd.io/etcd/client/v3"
	"time"
)
var (
	G_jobMgr *JobMgr
)

type JobMgr struct {
	Client  *client3.Client
	kv  client3.KV
	lease  client3.Lease
}
func InitJobMag(conf Config)error {
	config := client3.Config{
		Endpoints:   conf.EtcdEndpoints,
		DialTimeout: time.Duration(conf.EtcdDialTimeout)*time.Millisecond,
	}
	//
	client ,err:=client3.New(config)
	if err!=nil{
		fmt.Println(err)
		return err
	}
	G_jobMgr = &JobMgr{
		Client: client,
		kv: client3.NewKV(client),
		lease: client3.NewLease(client),

	}
	return nil
}
//任务保存
func  (jobMgr *JobMgr)SaveJob(job *common.Job)(oldJob *common.Job,err error)  {
	jobKey := "/cron/jobs"+job.Name
	jobValue ,err := json.Marshal(job)
	if err!=nil{
		fmt.Println("job序列化失败:",err.Error())
		return nil, err
	}
	//保存到etcd
	putResp,err:=jobMgr.kv.Put(context.Background(),jobKey,string(jobValue),client3.WithPrevKV())
	if err!=nil{
		fmt.Println("put etcd:",err.Error())
		return nil,err
	}
	if putResp.PrevKv!=nil{
		err:= json.Unmarshal(putResp.PrevKv.Value,&oldJob)
		if err!=nil{
			return nil,err
		}
		return  oldJob,nil
	}
	return
}