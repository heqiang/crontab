package master

import (
	"context"
	"crontab/common"
	"encoding/json"
	"errors"
	"fmt"
	client3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	G_jobMgr *JobMgr
)

var JobNotExist error = errors.New("不存在该任务")

type JobMgr struct {
	Client *client3.Client
	kv     client3.KV
	lease  client3.Lease
}

func InitJobMag(conf Config) error {
	config := client3.Config{
		Endpoints:   conf.EtcdEndpoints,
		DialTimeout: time.Duration(conf.EtcdDialTimeout) * time.Millisecond,
	}
	//
	client, err := client3.New(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	G_jobMgr = &JobMgr{
		Client: client,
		kv:     client3.NewKV(client),
		lease:  client3.NewLease(client),
	}
	return nil
}

//任务保存
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	jobKey := common.JobKeyPrefix + job.Name
	jobValue, err := json.Marshal(job)
	if err != nil {
		fmt.Println("job序列化失败:", err.Error())
		return nil, err
	}
	//保存到etcd
	putResp, err := jobMgr.kv.Put(context.Background(), jobKey, string(jobValue), client3.WithPrevKV())
	if err != nil {
		fmt.Println("put etcd:", err.Error())
		return nil, err
	}
	if putResp.PrevKv != nil {
		fmt.Println(string(putResp.PrevKv.Key))
		err := json.Unmarshal(putResp.PrevKv.Value, &oldJob)
		if err != nil {
			return nil, err
		}
		return oldJob, nil
	}
	return
}

func (jobMgr *JobMgr) DeleteJob(jobName string) (oldJob *common.Job, err error) {
	jobKey := common.JobKeyPrefix + jobName
	deleteResp, err := jobMgr.kv.Delete(context.Background(), jobKey, client3.WithPrevKV())
	if err != nil {
		fmt.Println("job 删除失败")
		return
	}
	if len(deleteResp.PrevKvs) != 0 {
		err := json.Unmarshal(deleteResp.PrevKvs[0].Value, &oldJob)
		if err != nil {
			return nil, err
		}
		return oldJob, nil
	}
	return nil, JobNotExist
}

func (jobMgr *JobMgr) GetJobList() (jobList []*common.Job, err error) {
	getResp, err := jobMgr.kv.Get(context.Background(), common.JobKeyPrefix, client3.WithPrefix())
	if err != nil {
		return nil, err
	}
	jobList = make([]*common.Job, 0)
	for _, job := range getResp.Kvs {
		var job1 common.Job
		json.Unmarshal(job.Value, &job1)
		jobList = append(jobList, &job1)
	}
	return jobList, nil
}
func (jobMgr *JobMgr) KillJob(name string) error {
	killKey := common.JobKillDir + name
	grant, err := jobMgr.lease.Grant(context.Background(), 1)
	if err != nil {
		return err
	}

	_, err = jobMgr.kv.Put(context.Background(), killKey, "", client3.WithLease(grant.ID))
	if err != nil {
		return err
	}

	return nil

}
