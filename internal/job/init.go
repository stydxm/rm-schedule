package job

import (
	"github.com/scutrobotlab/rm-schedule/internal/common"
	"log"

	"github.com/robfig/cron/v3"
)

type CronJobParam struct {
	Name            string
	Url             string
	ReplaceRMStatic bool
}

var Params = []CronJobParam{
	{
		Name:            common.UpstreamNameGroupRankInfo,
		Url:             common.UpstreamUrlGroupRankInfo,
		ReplaceRMStatic: false,
	},
	{
		Name:            common.UpstreamNameRobotData,
		Url:             common.UpstreamUrlRobotData,
		ReplaceRMStatic: false,
	},
	{
		Name:            common.UpstreamNameSchedule,
		Url:             common.UpstreamUrlSchedule,
		ReplaceRMStatic: true,
	},
}

func InitCronJob() *cron.Cron {
	var jobFuncArray []func()
	for _, param := range Params {
		jobFuncArray = append(jobFuncArray, CronJobFactory(param))
	}

	c := cron.New()

	for _, jobFunc := range jobFuncArray {
		_, err := c.AddFunc("@every 5s", jobFunc)
		if err != nil {
			log.Fatalf("cron add func failed: %v", err)
		}

		jobFunc()
	}

	return c
}
