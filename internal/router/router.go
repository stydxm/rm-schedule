package router

import (
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/handler"
	"github.com/scutrobotlab/rm-schedule/internal/job"
	"github.com/scutrobotlab/rm-schedule/internal/static"
)

var Params = map[string]handler.RouteHandlerParam{
	"group_rank_info": {
		Name:         "group_rank_info",
		Static:       false,
		CacheControl: "public, max-age=5",
		OriginalUrl:  job.Params["group_rank_info"].Url,
		Data:         static.GroupRankInfoBytes,
		SeasonMap: map[string][]byte{
			"2024": static.GroupRankInfoBytes2024,
		},
	},
	"robot_data": {
		Name:         "robot_data",
		Static:       false,
		CacheControl: "public, max-age=5",
		OriginalUrl:  job.Params["robot_data"].Url,
		Data:         static.RobotDataBytes,
		SeasonMap:    nil,
	},
	"schedule": {
		Name:         "schedule",
		Static:       false,
		CacheControl: "public, max-age=5",
		OriginalUrl:  job.Params["schedule"].Url,
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

	for handlerName, param := range Params {
		api.Get("/"+handlerName, handler.RouteHandlerFactory(param))
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
