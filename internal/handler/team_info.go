package handler

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/patrickmn/go-cache"
	"github.com/scutrobotlab/rm-schedule/internal/static"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
	"github.com/sirupsen/logrus"
)

const BilibiliOfficialKey = "bilibili_official"

// formerNameRegex 匹配"现名（原曾用名）"格式，如"广州城市理工学院（原华南理工大学广州学院）"
var formerNameRegex = regexp.MustCompile(`^(.+?)（原(.+?)）$`)

type TeamInfo struct {
	CollegeName string `json:"collegeName"`
	BilibiliUid int64  `json:"bilibiliUid"`
}

type BilibiliOfficial struct {
	School   string `json:"school"`
	Nickname string `json:"nickname"`
	Uid      int64  `json:"uid"`
}

func TeamInfoHandler(c iris.Context) {
	collegeName := c.URLParam("college_name")
	if collegeName == "" {
		c.StatusCode(400)
		c.JSON(iris.Map{"code": -1, "msg": "college_name is required"})
		return
	}

	var bilibiliOfficialMap map[string]BilibiliOfficial
	if cached, ok := svc.Cache.Get(BilibiliOfficialKey); ok {
		bilibiliOfficialMap = cached.(map[string]BilibiliOfficial)
	} else {
		var bilibiliOfficialList []BilibiliOfficial
		err := json.Unmarshal(static.BilibiliOfficialBytes, &bilibiliOfficialList)
		if err != nil {
			logrus.Errorf("Failed to parse Bilibili official data: %v", err)
			c.StatusCode(500)
			c.JSON(iris.Map{"code": -1, "msg": "Failed to parse Bilibili official data"})
			return
		}

		bilibiliOfficialMap = make(map[string]BilibiliOfficial)
		for _, official := range bilibiliOfficialList {
			school := strings.TrimSpace(official.School)
			if matches := formerNameRegex.FindStringSubmatch(school); matches != nil {
				// 学校有曾用名，eg. 广州城市理工学院（原华南理工大学广州学院）
				// matches[1]=现名, matches[2]=曾用名，均注册到 map
				bilibiliOfficialMap[strings.TrimSpace(matches[1])] = official
				bilibiliOfficialMap[strings.TrimSpace(matches[2])] = official
			} else {
				bilibiliOfficialMap[school] = official
			}
		}
		svc.Cache.Set(BilibiliOfficialKey, bilibiliOfficialMap, cache.NoExpiration)
	}
	bilibiliOfficial, ok := bilibiliOfficialMap[collegeName]
	if !ok {
		c.StatusCode(404)
		c.JSON(iris.Map{"code": -1, "msg": "School not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(TeamInfo{
		CollegeName: collegeName,
		BilibiliUid: bilibiliOfficial.Uid,
	})
}
