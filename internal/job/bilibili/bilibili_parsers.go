package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/scutrobotlab/rm-schedule/internal/common"
	"github.com/scutrobotlab/rm-schedule/internal/router"
	"github.com/scutrobotlab/rm-schedule/internal/static"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
	"github.com/scutrobotlab/rm-schedule/internal/types"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var schedules = map[string][]byte{
	"2024": static.ScheduleBytes2024,
	"2025": router.RedirectParams[common.UpstreamNameSchedule].Data,
}

type Matches map[string]map[string][]types.MatchNode

func getMatches() Matches {
	//返回赛季、赛区、比赛三级嵌套的map,列出从日程中获取的所有比赛信息
	matches := Matches{}

	if !router.RedirectParams[common.UpstreamNameSchedule].Static {
		liveScheduleData, ok := svc.Cache.Get(router.RedirectParams[common.UpstreamNameSchedule].Name)
		if ok {
			schedules["2025"] = liveScheduleData.([]byte)
		}
	}

	for season, scheduleBytes := range schedules {
		var scheduleData types.ScheduleResp
		err := json.Unmarshal(scheduleBytes, &scheduleData)
		if err != nil {
			logrus.Errorf("failed to unmarshal schedule: %v\n", err)
			return nil
		}

		//遍历整理比赛
		zones := make(map[string][]types.MatchNode)
		for _, currentZone := range scheduleData.Data.Event.Zones.Nodes {
			var currentMatches []types.MatchNode
			for _, match := range currentZone.GroupMatches.Nodes {
				currentMatches = append(currentMatches, match)
			}
			for _, match := range currentZone.KnockoutMatches.Nodes {
				currentMatches = append(currentMatches, match)
			}
			zones[currentZone.Name] = currentMatches
		}
		matches[season] = zones
	}

	return matches
}

func findCollection(season string, zone string, collectionList *[]types.BiliBiliCollectionMetaData) (types.BiliBiliCollectionMetaData, bool) {
	for _, collection := range *collectionList {
		collectionName := collection.Name
		isReplay := strings.Contains(collectionName, "比赛回放") || strings.Contains(collectionName, "直播回放")
		isRMUC := strings.Contains(collectionName, "RMUC") || strings.Contains(collectionName, "超级对抗赛")
		isCorrectSeason := strings.Contains(collectionName, season)
		var isCorrectZone bool
		if len([]rune(zone)) > 3 {
			//港澳台及海外赛区&复活赛在赛程中被拆分为两段，回放中属于同一合集
			isCorrectZone = strings.Contains(collectionName, string([]rune(zone)[:3]))
		} else {
			isCorrectZone = strings.Contains(collectionName, zone)
		}
		if isReplay && isRMUC && isCorrectSeason && isCorrectZone {
			return collection, true
		}
	}
	return types.BiliBiliCollectionMetaData{Name: "not found", CollectionId: 0}, false
}

func checkStringInclude(fullStr string, subStr string) bool {
	var isInclusion bool
	if len([]rune(subStr)) > 3 {
		isInclusion = strings.Contains(fullStr, string([]rune(subStr)[:3])) || strings.Contains(fullStr, string([]rune(subStr)[3:]))
	} else {
		isInclusion = strings.Contains(fullStr, subStr)
	}
	return isInclusion
}

func findMatchVideo(match *types.MatchNode, collection *types.BiliBiliCollectionInfo) (types.BiliBiliVideoMetaData, bool) {
	for _, video := range collection.Data.Archives {
		videoTitle := video.Title
		isCorrectOrderNum := strings.Contains(videoTitle, fmt.Sprintf("第%d场", match.OrderNumber)) ||
			strings.Contains(videoTitle, fmt.Sprintf("第 %d 场", match.OrderNumber))
		isCorrectBlueTeam := checkStringInclude(videoTitle, match.BlueSide.Player.Team.Name)
		isCorrectBlueSchool := checkStringInclude(videoTitle, match.BlueSide.Player.Team.CollegeName)
		isCorrectRedTeam := checkStringInclude(videoTitle, match.RedSide.Player.Team.Name)
		isCorrectRedSchool := checkStringInclude(videoTitle, match.RedSide.Player.Team.CollegeName)
		if isCorrectOrderNum && isCorrectBlueTeam && isCorrectBlueSchool && isCorrectRedTeam && isCorrectRedSchool {
			return video, true
		}
	}
	return types.BiliBiliVideoMetaData{}, false
}
func stringKeyToIntKey[T any](stringKeyMap map[string]T) map[int]T {
	//int作为键的数据存入json时会被变为string，该函数将键转换回去
	intKeyMap := make(map[int]T)
	for key, value := range stringKeyMap {
		intKey, err := strconv.Atoi(key)
		if err == nil {
			intKeyMap[intKey] = value
		}
	}
	return intKeyMap
}
