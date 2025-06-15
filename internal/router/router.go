package router

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/common"
	"github.com/scutrobotlab/rm-schedule/internal/handler"
	"github.com/scutrobotlab/rm-schedule/internal/static"
)

var Params = map[string]handler.RouteHandlerParam{
	common.UpstreamNameGroupRankInfo: {
		Name:         common.UpstreamNameGroupRankInfo,
		Static:       false,
		CacheControl: "public, max-age=5",
		OriginalUrl:  common.UpstreamUrlGroupRankInfo,
		Data:         static.GroupRankInfoBytes,
		SeasonMap: map[string][]byte{
			"2024": static.GroupRankInfoBytes2024,
		},
	},
	common.UpstreamNameRobotData: {
		Name:         common.UpstreamNameRobotData,
		Static:       false,
		CacheControl: "public, max-age=5",
		OriginalUrl:  common.UpstreamUrlRobotData,
		Data:         static.RobotDataBytes,
		SeasonMap:    nil,
	},
	common.UpstreamNameSchedule: {
		Name:         common.UpstreamNameSchedule,
		Static:       false,
		CacheControl: "public, max-age=5",
		OriginalUrl:  common.UpstreamUrlSchedule,
		Data:         static.ScheduleBytes,
		SeasonMap: map[string][]byte{
			"2024": static.ScheduleBytes2024,
		},
	},
}

// Router defines the router for this service
func Router(r *iris.Application, frontend string) {
	api := r.Party("/api")
	api.Get("/static/*path", handler.RMStaticHandler)
	api.Get("/mp/match", handler.MpMatchHandler)
	api.Get("/rank", handler.RankListHandler)

	for name, param := range Params {
		api.Get(fmt.Sprintf("/%s", name), handler.RouteHandlerFactory(param))
	}

	r.HandleDir("/", iris.Dir(frontend), iris.DirOptions{
		IndexName: "index.html",
		ShowList:  false,
		Compress:  true,
	})

	// on 404, redirect to the index.html
	r.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
	})
}
