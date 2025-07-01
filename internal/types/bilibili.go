package types

// B站返回值中“合集”被翻译为`season`，这里为了避免与赛季混淆，均改为`collection`

// BiliBiliCollectionMetaData 合集列表中单个合集的元数据
type BiliBiliCollectionMetaData struct {
	Name         string `json:"name"`
	CollectionId int    `json:"season_id"`
	Category     int    `json:"category"`
	Cover        string `json:"cover"`
	VideoCount   int    `json:"total"`
	CreateTime   int    `json:"ptime"`
}

// BiliBiliCollectionResponse B站接口返回的某用户合集列表
type BiliBiliCollectionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		ItemsLists struct {
			Page struct {
				PageNum  int `json:"page_num"`
				PageSize int `json:"page_size"`
				Total    int `json:"total"`
			} `json:"page"`
			CollectionList []struct {
				Meta BiliBiliCollectionMetaData `json:"meta"`
			} `json:"seasons_list"`
		} `json:"items_lists"`
	} `json:"data"`
}

// BiliBiliVideoMetaData 合集信息中单个视频的元数据
type BiliBiliVideoMetaData struct {
	AID         int    `json:"aid"`
	BVID        string `json:"bvid"`
	CreateTime  int    `json:"ctime"`
	EnableVT    bool   `json:"enable_vt"`
	Interactive bool   `json:"interactive_video"`
	Cover       string `json:"pic"`
	PublishTime int    `json:"pubdate"`
	Statistics  struct {
		Views    int `json:"view"`
		ViewTime int `json:"vt"`
	} `json:"stat"`
	Title string `json:"title"`
}

// BiliBiliCollectionInfo B站接口返回的合集信息
type BiliBiliCollectionInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		AIDs     []int                      `json:"aids"`
		Archives []BiliBiliVideoMetaData    `json:"archives"`
		Meta     BiliBiliCollectionMetaData `json:"meta"`
		Page     struct {
			PageNum  int `json:"page_num"`
			PageSize int `json:"page_size"`
			Total    int `json:"total"`
		} `json:"page"`
	} `json:"data"`
}

type MatchOrderToVideoType map[string]map[string]map[int]BiliBiliVideoMetaData
