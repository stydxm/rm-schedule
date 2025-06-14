package job

import (
	"log"

	"github.com/robfig/cron/v3"
)

type CronJobParam struct {
	Name, Url       string
	ReplaceRMStatic bool
}

var Params = map[string]CronJobParam{
	"group_rank_info": {
		Name:            "group_rank_info",
		Url:             "https://pro-robomasters-hz-n5i3.oss-cn-hangzhou.aliyuncs.com/live_json/group_rank_info.json",
		ReplaceRMStatic: false,
	},
	"robot_data": {
		Name:            "robot_data",
		Url:             "https://pro-robomasters-hz-n5i3.oss-cn-hangzhou.aliyuncs.com/live_json/robot_data.json",
		ReplaceRMStatic: false,
	},
	"schedule": {
		Name:            "schedule",
		Url:             "https://pro-robomasters-hz-n5i3.oss-cn-hangzhou.aliyuncs.com/live_json/schedule.json",
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
