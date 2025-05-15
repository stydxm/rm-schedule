package handler

import (
	"encoding/json"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"github.com/scutrobotlab/rm-schedule/internal/static"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
)

type RankListItem struct {
	RankScoreItem RankScoreItem `json:"rankScoreItem"`
	CompleteForm  CompleteForm  `json:"completeForm"`
}

type RankScoreItem struct {
	Rank          int     `json:"rank"`
	SchoolChinese string  `json:"schoolChinese"`
	SchoolEnglish string  `json:"schoolEnglish"`
	Score         float64 `json:"score"`
}

type CompleteForm struct {
	Rank                  int    `json:"rank"`
	School                string `json:"school"`
	Team                  string `json:"team"`
	Score                 int    `json:"score"`
	InitialCoinDocument   int    `json:"initialCoinDocument"`
	LevelDocument         string `json:"levelDocument"`
	InitialCoinTechnology int    `json:"initialCoinTechnology"`
	LevelTechnology       string `json:"levelTechnology"`
	InitialCoinTotal      int    `json:"initialCoinTotal"`
}

var SeasonCompleteFormMap = map[string][]byte{
	"2024": static.CompleteFormBytes2024,
}

var SeasonRankScoreMap = map[string][]byte{
	"2024": static.RankScoreBytes2024,
}

func RankListHandler(c iris.Context) {
	season := c.URLParam("season")
	var completeFormBytes []byte
	var rankScoreBytes []byte
	if data, ok := SeasonCompleteFormMap[season]; ok {
		completeFormBytes = data
	} else {
		completeFormBytes = static.CompleteFormBytes
	}
	if data, ok := SeasonRankScoreMap[season]; ok {
		rankScoreBytes = data
	} else {
		rankScoreBytes = static.RankScoreBytes
	}

	schoolName := c.URLParam("school_name")
	if schoolName == "" {
		c.StatusCode(400)
		c.JSON(iris.Map{"code": -1, "msg": "School name is empty"})
		return
	}

	completedFormMap, ok := svc.Cache.Get("completed_form")
	if !ok {
		completedFormJson := make([]CompleteForm, 0)
		err := json.Unmarshal(completeFormBytes, &completedFormJson)
		if err != nil {
			log.Printf("Failed to parse completed form: %v\n", err)
			c.StatusCode(500)
			c.JSON(iris.Map{"code": -1, "msg": "Failed to parse completed form"})
			return
		}

		// 并列名次处理
		var rank int
		var lastCoinTotal int
		for i := range completedFormJson {
			if completedFormJson[i].InitialCoinTotal != lastCoinTotal {
				rank = i + 1
			}
			completedFormJson[i].Rank = rank
			lastCoinTotal = completedFormJson[i].InitialCoinTotal
		}
		completedFormMap = lo.SliceToMap(completedFormJson, func(item CompleteForm) (string, CompleteForm) { return item.School, item })
		svc.Cache.Set("completed_form", completedFormMap, cache.NoExpiration)
	}

	completedForm, ok := completedFormMap.(map[string]CompleteForm)[schoolName]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "School not found"})
		return
	}

	rankScoreMap, ok := svc.Cache.Get("rank_score")
	if !ok {
		rankScoreJson := make([]RankScoreItem, 0)
		err := json.Unmarshal(rankScoreBytes, &rankScoreJson)
		if err != nil {
			log.Printf("Failed to parse rank list: %v\n", err)
			c.StatusCode(500)
			c.JSON(iris.Map{"code": -1, "msg": "Failed to parse rank list"})
			return
		}

		rankScoreMap = lo.SliceToMap(rankScoreJson, func(item RankScoreItem) (string, RankScoreItem) { return item.SchoolChinese, item })
		svc.Cache.Set("rank_list", rankScoreMap, cache.NoExpiration)
	}

	rankScore, ok := rankScoreMap.(map[string]RankScoreItem)[schoolName]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "School not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(RankListItem{
		RankScoreItem: rankScore,
		CompleteForm:  completedForm,
	})
}
