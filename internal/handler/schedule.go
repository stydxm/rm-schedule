package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/job"
	"github.com/scutrobotlab/rm-schedule/internal/static"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
)

const (
	ScheduleStatic       = false
	ScheduleCacheControl = "public, max-age=5"
)

// SeasonScheduleMap 赛季赛程映射
var SeasonScheduleMap = map[string][]byte{
	"2024": static.ScheduleBytes2024,
}

func ScheduleHandler(c iris.Context) {
	season := c.URLParam("season")
	if data, ok := SeasonScheduleMap[season]; ok {
		c.Header("Cache-Control", "public, max-age=60")
		c.ContentType("application/json")
		c.Write(data)
		return
	}

	if ScheduleStatic {
		c.Header("Cache-Control", "public, max-age=60")
		c.ContentType("application/json")
		c.Write(static.ScheduleBytes)
		return
	}

	// 是否存在 Tencent-Acceleration-Domain-Name
	if c.GetHeader("Tencent-Acceleration-Domain-Name") != "" {
		c.Header("Cache-Control", ScheduleCacheControl)
		c.Redirect(job.ScheduleUrl, 301)
		return
	}

	if cached, b := svc.Cache.Get("schedule"); b {
		c.Header("Cache-Control", ScheduleCacheControl)
		c.ContentType("application/json")
		c.Write(cached.([]byte))
		return
	}

	c.Header("Cache-Control", ScheduleCacheControl)
	c.StatusCode(500)
	c.JSON(iris.Map{"code": -1, "msg": "Failed to get schedule"})
}
