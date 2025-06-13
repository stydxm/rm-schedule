package job

import (
	"github.com/patrickmn/go-cache"
	"github.com/scutrobotlab/rm-schedule/internal/svc"
	"io"
	"log"
	"net/http"
)

const RobotDataUrl = "https://pro-robomasters-hz-n5i3.oss-cn-hangzhou.aliyuncs.com/live_json/robot_data.json"

func UpdateRobotData() {
	resp, err := http.Get(RobotDataUrl)
	if err != nil {
		log.Printf("Failed to get robot_data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get robot_data: status code %d\n", resp.StatusCode)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read robot_data: %v\n", err)
		return
	}
	bytes = replaceRMStatic(bytes)
	svc.Cache.Set("robot_data", bytes, cache.DefaultExpiration)

	log.Println("RobotData updated")
}
