package types

type ScheduleResp struct {
	Data ScheduleData `json:"data"`
}

type ScheduleData struct {
	Event     Event     `json:"event"`
	LastEvent LastEvent `json:"last_event"`
}

type Event struct {
	Title string `json:"title"`
	Zones Zones  `json:"zones"`
}

type Zones struct {
	Nodes []ZoneNode `json:"nodes"`
}

type ZoneNode struct {
	ID              string   `json:"id"`
	MatchDates      []string `json:"matchDates"`
	Name            string   `json:"name"`
	ZoneType        string   `json:"zoneType"`
	Groups          Groups   `json:"groups"`
	GroupMatches    Matches  `json:"groupMatches"`
	KnockoutMatches Matches  `json:"knockoutMatches"`
}

type Groups struct {
	Nodes []GroupNode `json:"nodes"`
}

type GroupNode struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Players Players `json:"players"`
}

type Players struct {
	Nodes []PlayerNode `json:"nodes"`
}

type PlayerNode struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Rank   int     `json:"rank"`
	Score  int     `json:"score"`
	TeamID *string `json:"teamId,omitempty"`
	Team   *Team   `json:"team,omitempty"`
}

type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CollegeLogo string `json:"collegeLogo"`
	CollegeName string `json:"collegeName"`
}

type Matches struct {
	Nodes []MatchNode `json:"nodes"`
}

type MatchNode struct {
	ID                   string      `json:"id"`
	GroupID              string      `json:"groupId"`
	MatchType            string      `json:"matchType"`
	OrderNumber          int         `json:"orderNumber"`
	PlanGameCount        int         `json:"planGameCount"`
	PlanStartedAt        string      `json:"planStartedAt"`
	Result               string      `json:"result"`
	Slug                 interface{} `json:"slug"`
	SlugName             string      `json:"slugName"`
	Status               string      `json:"status"`
	WinnerPlaceholdName  interface{} `json:"winnerPlaceholdName"`
	LoserPlaceholdName   interface{} `json:"loserPlaceholdName"`
	BlueSideID           string      `json:"blueSideId"`
	BlueSideScore        int         `json:"blueSideScore"`
	BlueSideWinGameCount int         `json:"blueSideWinGameCount"`
	BlueSide             BlueSide    `json:"blueSide"`
	RedSideID            string      `json:"redSideId"`
	RedSideScore         int         `json:"redSideScore"`
	RedSideWinGameCount  int         `json:"redSideWinGameCount"`
	RedSide              RedSide     `json:"redSide"`
}

type BlueSide struct {
	ID               string  `json:"id"`
	PreparedStatus   string  `json:"preparedStatus"`
	FillSourceID     *string `json:"fillSourceId,omitempty"`
	FillSourceType   *string `json:"fillSourceType,omitempty"`
	FillSourceNumber *int    `json:"fillSourceNumber,omitempty"`
	FillStatus       string  `json:"fillStatus"`
	PlayerID         *string `json:"playerId,omitempty"`
	Player           *Player `json:"player,omitempty"`
	UpdatedAt        string  `json:"updatedAt"`
}

type Player struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Rank   int     `json:"rank"`
	Score  int     `json:"score"`
	TeamID *string `json:"teamId,omitempty"`
	Team   *Team   `json:"team,omitempty"`
}

type PlayerWithMatch struct {
	Player Player    `json:"player"`
	Match  MatchNode `json:"match"`
}

type RedSide struct {
	ID               string  `json:"id"`
	PreparedStatus   string  `json:"preparedStatus"`
	FillSourceID     *string `json:"fillSourceId,omitempty"`
	FillSourceType   *string `json:"fillSourceType,omitempty"`
	FillSourceNumber *int    `json:"fillSourceNumber,omitempty"`
	FillStatus       string  `json:"fillStatus"`
	PlayerID         *string `json:"playerId,omitempty"`
	Player           *Player `json:"player,omitempty"`
	UpdatedAt        string  `json:"updatedAt"`
}

type LastEvent struct {
	Title string `json:"title"`
	Zones Zones  `json:"zones"`
}
