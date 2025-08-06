package analyze

import (
	"encoding/csv"
	"encoding/json"
	"github.com/kataras/iris/v12/x/errors"
	"github.com/scutrobotlab/rm-schedule/internal/common"
	"github.com/scutrobotlab/rm-schedule/internal/types"
	"io"
	"net/http"
	"os"
	"strconv"
)

// GetScheduleData 获取赛程数据
func GetScheduleData() ([]byte, error) {
	resp, err := http.Get(common.UpstreamUrlSchedule)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// GetSchedule 获取赛程数据并解析为 ScheduleResp 结构体
func GetSchedule() (*types.ScheduleResp, error) {
	data, err := GetScheduleData()
	if err != nil {
		return nil, err
	}

	var schedule types.ScheduleResp
	if err := json.Unmarshal(data, &schedule); err != nil {
		return nil, err
	}

	return &schedule, nil
}

// ExportScheduleToFile 导出赛程数据到文件
func ExportScheduleToFile(filename string) error {
	schedule, err := GetSchedule()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"序号",
		"赛季",
		"赛区",
		"赛次",
		"比赛阶段",
		"红方学校",
		"红队名",
		"蓝方学校",
		"蓝队名",
		"红方比分",
		"蓝方比分",
		"备注",
	})

	writeMatch := func(zone types.ZoneNode, match types.MatchNode) {
		writer.Write([]string{
			match.ID,
			"2025",
			zone.Name,
			strconv.Itoa(match.OrderNumber),
			match.MatchType,
			match.RedSide.Player.Team.CollegeName,
			match.RedSide.Player.Team.Name,
			match.BlueSide.Player.Team.CollegeName,
			match.BlueSide.Player.Team.Name,
			strconv.Itoa(match.RedSideWinGameCount),
			strconv.Itoa(match.BlueSideWinGameCount),
			"",
		})
	}

	event := schedule.Data.Event
	for _, zone := range event.Zones.Nodes {
		for _, match := range zone.GroupMatches.Nodes {
			writeMatch(zone, match)
		}
		for _, match := range zone.KnockoutMatches.Nodes {
			writeMatch(zone, match)
		}
	}

	return nil
}
