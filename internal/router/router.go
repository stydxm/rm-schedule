package router

import (
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/common"
	"github.com/scutrobotlab/rm-schedule/internal/handler"
)

// Router defines the router for this service
func Router(r *iris.Application, frontend string) {
	api := r.Party("/api")
	api.Get("/static/*path", handler.RMStaticHandler)
	api.Get("/mp/match", handler.MpMatchHandler)
	api.Get("/rank", handler.RankListHandler)
	api.Get("/group_rank_info", handler.RedirectRouteHandlerFactory(RedirectParams[common.UpstreamNameGroupRankInfo]))
	api.Get("/robot_data", handler.RedirectRouteHandlerFactory(RedirectParams[common.UpstreamNameRobotData]))
	api.Get("/schedule", handler.RedirectRouteHandlerFactory(RedirectParams[common.UpstreamNameSchedule]))

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
