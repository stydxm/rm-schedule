package static

import _ "embed"

// 公共的静态文件

//go:embed complete_form.json
var CompleteFormBytes []byte

//go:embed complete_form_rank.json
var CompleteFormRankBytes []byte

//go:embed rank_score.json
var RankScoreBytes []byte

//go:embed schedule.json
var ScheduleBytes []byte

//go:embed group_rank_info.json
var GroupRankInfoBytes []byte

//go:embed robot_data.json
var RobotDataBytes []byte

//go:embed bilibili_official.json
var BilibiliOfficialBytes []byte

//go:embed history_match.json
var HistoryMatchBytes []byte

//go:embed bilibili_videos.json
var BilibiliVideosBytes []byte

// 2024 赛季的静态文件

//go:embed season_2024/complete_form.json
var CompleteFormBytes2024 []byte

//go:embed season_2024/rank_score.json
var RankScoreBytes2024 []byte

//go:embed season_2024/schedule.json
var ScheduleBytes2024 []byte

//go:embed season_2024/group_rank_info.json
var GroupRankInfoBytes2024 []byte

//go:embed season_2024/bilibili_videos.json
var BilibiliVideosBytes2024 []byte
