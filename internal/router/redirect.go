package router

import (
	"github.com/scutrobotlab/rm-schedule/internal/common"
	"github.com/scutrobotlab/rm-schedule/internal/handler"
	"github.com/scutrobotlab/rm-schedule/internal/static"
)

// RedirectParams 定义重定向路由的参数
var RedirectParams = map[string]handler.RedirectRouteHandlerParam{
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
