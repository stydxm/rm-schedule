package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12/x/errors"
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

type CompleteFormRank struct {
	Rank   int    `json:"rank"`
	School string `json:"school"`
	Team   string `json:"team"`
}

var SeasonCompleteFormMap = map[string][]byte{
	"2024": static.CompleteFormBytes2024,
}

var SeasonCompleteFormRankMap = map[string][]byte{
	"2024": {},
}

var SeasonRankScoreMap = map[string][]byte{
	"2024": static.RankScoreBytes2024,
}

func RankListHandler(c iris.Context) {
	season := c.URLParam("season")
	rankScoreKey := fmt.Sprintf("rank_score_%s", season)
	var rankScoreBytes []byte
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

	completeFormMap, err := GetCompleteFormMap(season)
	if err != nil {
		c.StatusCode(500)
		c.JSON(iris.Map{"code": -1, "msg": "Failed to get complete form"})
		return
	}
	completeForm, ok := completeFormMap[schoolName]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "School not found"})
		return
	}

	rankScoreMap, ok := svc.Cache.Get(rankScoreKey)
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
		svc.Cache.Set(rankScoreKey, rankScoreMap, cache.NoExpiration)
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
		CompleteForm:  completeForm,
	})
}

// GetCompleteFormMap 获取完整形态
func GetCompleteFormMap(season string) (map[string]CompleteForm, error) {
	completeFormKey := fmt.Sprintf("complete_form_%s", season)
	ret, ok := svc.Cache.Get(completeFormKey)
	if ok {
		return ret.(map[string]CompleteForm), nil
	}

	var completeFormBytes []byte
	var completeFormRankBytes []byte
	if data, ok := SeasonCompleteFormMap[season]; ok {
		completeFormBytes = data
	} else {
		completeFormBytes = static.CompleteFormBytes
	}
	if data, ok := SeasonCompleteFormRankMap[season]; ok {
		completeFormRankBytes = data
	} else {
		completeFormRankBytes = static.CompleteFormRankBytes
	}

	completeFormJson := make([]CompleteForm, 0)
	err := json.Unmarshal(completeFormBytes, &completeFormJson)
	if err != nil {
		log.Printf("Failed to parse complete form: %v\n", err)
		return nil, errors.New("Failed to parse complete form")
	}

	if len(completeFormRankBytes) != 0 {
		// 有完整形态排名
		completeFormRankJson := make([]CompleteFormRank, 0)
		err = json.Unmarshal(completeFormRankBytes, &completeFormRankJson)
		if err != nil {
			log.Printf("Failed to parse complete form rank: %v\n", err)
			return nil, errors.New("Failed to parse complete form rank")
		}
		completeFormRankMap := make(map[string]CompleteFormRank)
		for _, item := range completeFormRankJson {
			completeFormRankMap[item.School] = item
		}
		for i, item := range completeFormJson {
			if rankItem, ok := completeFormRankMap[item.School]; ok {
				completeFormJson[i].Rank = rankItem.Rank
			}
		}
	} else {
		// 无完整形态排名 按照金币数量计算
		// 并列名次处理
		var rank int
		var lastCoinTotal int
		for i := range completeFormJson {
			if completeFormJson[i].InitialCoinTotal != lastCoinTotal {
				rank = i + 1
			}
			completeFormJson[i].Rank = rank
			lastCoinTotal = completeFormJson[i].InitialCoinTotal
		}
	}
	completeFormMap := lo.SliceToMap(completeFormJson, func(item CompleteForm) (string, CompleteForm) { return item.School, item })
	svc.Cache.Set(completeFormKey, completeFormMap, cache.NoExpiration)

	return completeFormMap, nil
}
