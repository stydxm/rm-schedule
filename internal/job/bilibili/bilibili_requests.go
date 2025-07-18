package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/scutrobotlab/rm-schedule/internal/common"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
	"github.com/scutrobotlab/rm-schedule/internal/types"
	"io"
	"log"
	"net/http"
)

func processRequest(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:139.0) Gecko/20100101 Firefox/139.0")
	req.Header.Set("Host", "api.bilibili.com")
	req.Header.Set("Referer", "https://space.bilibili.com/")
	req.Header.Set("Origin", "https://space.bilibili.com/")
}

func getCollectionList() ([]types.BiliBiliCollectionMetaData, error) {
	pageNum := 1
	moreCollections := true
	var collections []types.BiliBiliCollectionMetaData
	for pageNum <= 10 && moreCollections {
		req, err := http.NewRequest("GET", fmt.Sprintf(common.UpstreamUrlBilibiliCollections, pageNum), nil)
		if err != nil {
			return []types.BiliBiliCollectionMetaData{}, fmt.Errorf("Failed to get collections page %d: %v\n", pageNum, err)
		}
		processRequest(req)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return []types.BiliBiliCollectionMetaData{}, fmt.Errorf("Failed to get collections page %d: %v\n", pageNum, err)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return []types.BiliBiliCollectionMetaData{}, fmt.Errorf("Failed to read collections page %d: %v\n", pageNum, err)
		}

		var data types.BiliBiliCollectionResponse
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			return []types.BiliBiliCollectionMetaData{}, fmt.Errorf("Failed to read collections page %d: %v\n", pageNum, err)
		}
		if data.Code != 0 || data.Message != "0" {
			return []types.BiliBiliCollectionMetaData{}, fmt.Errorf("Failed to read collections page %d: %s\n", pageNum, data.Message)
		}

		if data.Data.ItemsLists.Page.Total <= data.Data.ItemsLists.Page.PageSize*pageNum {
			moreCollections = false
		} else {
			pageNum += 1
		}

		for _, collectionMeta := range data.Data.ItemsLists.CollectionList {
			collections = append(collections, collectionMeta.Meta)
		}

		resp.Body.Close()
	}
	return collections, nil
}

func getCollectionInfo(id int) (types.BiliBiliCollectionInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.bilibili.com/x/polymer/web-space/seasons_archives_list?season_id=%d&page_size=100&page_num=1", id), nil)
	if err != nil {
		return types.BiliBiliCollectionInfo{}, fmt.Errorf("Failed to get collection %d: %v\n", id, err)
	}
	processRequest(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.BiliBiliCollectionInfo{}, fmt.Errorf("Failed to get collection %d: %v\n", id, err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.BiliBiliCollectionInfo{}, fmt.Errorf("Failed to get collection %d: %v\n", id, err)
	}

	var data types.BiliBiliCollectionInfo
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return types.BiliBiliCollectionInfo{}, fmt.Errorf("Failed to get collection %d: %v\n", id, err)
	}
	if data.Code != 0 || data.Message != "0" {
		return types.BiliBiliCollectionInfo{}, fmt.Errorf("Failed to get collection %d: %v\n", id, data.Message)
	}
	return data, nil
}

func FetchBiliBiliReplayVideos() {
	collectionList, err := getCollectionList()
	if err != nil {
		log.Println(err)
		return
	}
	matchIDToVideo := make(map[string]types.BiliBiliVideoMetaData)
	matchOrderToVideoAll := make(types.MatchOrderToVideoType)
	matchesList := getMatches()
	for season, zones := range matchesList {
		matchOrderToVideoSeason := make(map[string]map[int]types.BiliBiliVideoMetaData)
		for zone, matches := range zones {
			matchOrderToVideoZone := make(map[int]types.BiliBiliVideoMetaData)
			targetCollection, ok := findCollection(season, zone, &collectionList)
			if ok {
				collectionInfo, err := getCollectionInfo(targetCollection.CollectionId)
				if err != nil {
					continue
				}
				for _, match := range matches {
					correspondingVideo, ok := findMatchVideo(&match, &collectionInfo)
					if !ok {
						continue
					}
					matchIDToVideo[match.ID] = correspondingVideo
					matchOrderToVideoZone[match.OrderNumber] = correspondingVideo
				}
			}
			if len(matchOrderToVideoZone) > 0 {
				matchOrderToVideoSeason[zone] = matchOrderToVideoZone
			}
		}
		if len(matchOrderToVideoSeason) > 0 {
			matchOrderToVideoAll[season] = matchOrderToVideoSeason
		}
	}
	svc.Cache.Set("match_id_to_video", matchIDToVideo, cache.DefaultExpiration)
	if len(matchOrderToVideoAll) > 0 {
		svc.Cache.Set("match_order_to_video", matchOrderToVideoAll, cache.DefaultExpiration)
	}
	log.Printf("Bilibili collections updated")
}
