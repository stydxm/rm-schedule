package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
	"github.com/scutrobotlab/rm-schedule/internal/types"
	"strconv"
)

func MatchIDHandler(c iris.Context) {
	matchId := c.URLParam("match_id")
	if matchId == "" {
		c.StatusCode(400)
		c.JSON(iris.Map{"error": "match_id is required"})
		return
	}

	_ret, ok := svc.Cache.Get("match_id_to_video")
	if !ok {
		c.StatusCode(500)
		c.JSON(iris.Map{"code": -1, "msg": "Failed to get videos"})
		return
	}

	ret := _ret.(map[string]types.BiliBiliVideoMetaData)
	video, ok := ret[matchId]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "Match not found"})
		return
	}
	c.JSON(video)
	return
}

func MatchOrderHandler(c iris.Context) {
	season := c.URLParam("season")
	zone := c.URLParam("zone")
	_matchOrder := c.URLParam("match_order")
	if season == "" || zone == "" || _matchOrder == "" {
		c.StatusCode(400)
		c.JSON(iris.Map{"error": "season & zone & match_order is required"})
		return
	}
	matchOrder, err := strconv.Atoi(_matchOrder)
	if err != nil {
		c.StatusCode(400)
		c.JSON(iris.Map{"error": "match_order should be int"})
		return
	}

	_ret, ok := svc.Cache.Get("match_order_to_video")
	if !ok {
		c.StatusCode(500)
		c.JSON(iris.Map{"code": -1, "msg": "Failed to get videos"})
		return
	}
	ret := _ret.(types.MatchOrderToVideoType)

	selectedSeason, ok := ret[season]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "Season not found"})
		return
	}
	selectedZone, ok := selectedSeason[zone]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "Zone not found"})
		return
	}
	video, ok := selectedZone[matchOrder]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "Match not found"})
		return
	}
	c.JSON(video)
	return
}
