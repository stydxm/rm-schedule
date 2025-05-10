package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/job"
	"github.com/scutrobotlab/rm-schedule/internal/static"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
)

const (
	GroupRankInfoStatic       = false
	GroupRankInfoCacheControl = "public, max-age=5"
)

var SeasonGroupRankInfoMap = map[string][]byte{
	"2024": static.GroupRankInfoBytes2024,
}

func GroupRankInfoHandler(c iris.Context) {
	season := c.URLParam("season")
	if data, ok := SeasonGroupRankInfoMap[season]; ok {
		c.Header("Cache-Control", "public, max-age=60")
		c.ContentType("application/json")
		c.Write(data)
		return
	}

	if GroupRankInfoStatic {
		c.Header("Cache-Control", "public, max-age=60")
		c.ContentType("application/json")
		c.Write(static.GroupRankInfoBytes)
		return
	}

	if c.GetHeader("Tencent-Acceleration-Domain-Name") != "" {
		c.Header("Cache-Control", GroupRankInfoCacheControl)
		c.Redirect(job.GroupRankInfoUrl, 301)
		return
	}

	if cached, b := svc.Cache.Get("group_rank_info"); b {
		c.Header("Cache-Control", GroupRankInfoCacheControl)
		c.ContentType("application/json")
		c.Write(cached.([]byte))
		return
	}

	c.Header("Cache-Control", GroupRankInfoCacheControl)
	c.StatusCode(500)
	c.JSON(iris.Map{"code": -1, "msg": "Failed to get group rank info"})
}
