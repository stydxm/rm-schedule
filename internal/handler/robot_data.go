package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/job"
	"github.com/scutrobotlab/rm-schedule/internal/static"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
)

const (
	RobotDataStatic       = false
	RobotDataCacheControl = "public, max-age=5"
)

func RobotDataHandler(c iris.Context) {
	if RobotDataStatic {
		c.Header("Cache-Control", "public, max-age=60")
		c.ContentType("application/json")
		c.Write(static.RobotDataBytes)
		return
	}

	// 是否存在 Tencent-Acceleration-Domain-Name
	if c.GetHeader("Tencent-Acceleration-Domain-Name") != "" {
		c.Header("Cache-Control", RobotDataCacheControl)
		c.Redirect(job.RobotDataUrl, 301)
		return
	}

	if cached, b := svc.Cache.Get("robot_data"); b {
		c.Header("Cache-Control", RobotDataCacheControl)
		c.ContentType("application/json")
		c.Write(cached.([]byte))
		return
	}

	c.Header("Cache-Control", RobotDataCacheControl)
	c.StatusCode(500)
	c.JSON(iris.Map{"code": -1, "msg": "Failed to get robot_data"})
}
