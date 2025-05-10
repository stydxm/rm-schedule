package static

import _ "embed"

// 公共的静态文件

//go:embed complete_form.json
var CompleteFormBytes []byte

//go:embed rank_score.json
var RankScoreBytes []byte

//go:embed schedule.json
var ScheduleBytes []byte

//go:embed group_rank_info.json
var GroupRankInfoBytes []byte

// 2024 赛季的静态文件

//go:embed season_2024/complete_form.json
var CompleteFormBytes2024 []byte

//go:embed season_2024/rank_score.json
var RankScoreBytes2024 []byte

//go:embed season_2024/schedule.json
var ScheduleBytes2024 []byte

//go:embed season_2024/group_rank_info.json
var GroupRankInfoBytes2024 []byte
