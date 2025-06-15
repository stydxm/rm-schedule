package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
)

// RedirectRouteHandlerParam 定义重定向路由处理器的参数
type RedirectRouteHandlerParam struct {
	Name         string
	CacheControl string
	OriginalUrl  string
	Static       bool
	SeasonMap    map[string][]byte
	Data         []byte
}

// RedirectRouteHandlerFactory 处理重定向路由的工厂函数
func RedirectRouteHandlerFactory(param RedirectRouteHandlerParam) func(c iris.Context) {
	return func(c iris.Context) {
		if param.SeasonMap != nil {
			season := c.URLParam("season")
			if data, ok := param.SeasonMap[season]; ok {
				c.Header("Cache-Control", "public, max-age=60")
				c.ContentType("application/json")
				_, err := c.Write(data)
				if err != nil {
					_ = c.JSON(iris.Map{"code": -1, "msg": "Failed to get " + param.Name})
				}
				return
			}
		}

		if param.Static {
			c.Header("Cache-Control", "public, max-age=60")
			c.ContentType("application/json")
			_, err := c.Write(param.Data)
			if err != nil {
				_ = c.JSON(iris.Map{"code": -1, "msg": "Failed to get " + param.Name})
			}
			return
		}

		// 是否存在 Tencent-Acceleration-Domain-Name
		if c.GetHeader("Tencent-Acceleration-Domain-Name") != "" {
			c.Header("Cache-Control", param.CacheControl)
			c.Redirect(param.OriginalUrl, 301)
			return
		}

		if cached, b := svc.Cache.Get(param.Name); b {
			c.Header("Cache-Control", param.CacheControl)
			c.ContentType("application/json")
			_, err := c.Write(cached.([]byte))
			if err != nil {
				_ = c.JSON(iris.Map{"code": -1, "msg": "Failed to get " + param.Name})
			}
			return
		}

		c.Header("Cache-Control", param.CacheControl)
		c.StatusCode(500)
		_ = c.JSON(iris.Map{"code": -1, "msg": "Failed to get " + param.Name})
	}
}
