package handler

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/patrickmn/go-cache"
	"github.com/scutrobotlab/rm-schedule/internal/static"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
	"strings"
)

const BilibiliOfficialKey = "bilibili_official"

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
			c.StatusCode(500)
			c.JSON(iris.Map{"code": -1, "msg": "Failed to parse Bilibili official data"})
			return
		}

		bilibiliOfficialMap = make(map[string]BilibiliOfficial)
		for _, official := range bilibiliOfficialList {
			school := strings.TrimSpace(official.School)
			if strings.Contains(school, "（原") {
				// 学校有曾用名
				// eg. 广州城市理工学院（原华南理工大学广州学院）
				for _, split := range strings.Split(school, "（原") {
					// 去掉末尾的括号和内容
					split = strings.TrimSuffix(split, "）")
					// 去掉前后的空格
					split = strings.TrimSpace(split)
					if split == "" {
						continue
					}
					bilibiliOfficialMap[split] = official
				}
			} else {
				// 学校没有曾用名
				bilibiliOfficialMap[official.School] = official
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
