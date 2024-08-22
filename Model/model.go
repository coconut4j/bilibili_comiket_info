package Model

import (
	"strings"
	"time"
)

type RespAll struct {
	Errno  int    `json:"errno"`
	Errtag int    `json:"errtag"`
	Msg    string `json:"msg"`
	Data   Data   `json:"data"`
}

type Data struct {
	Total      int           `json:"total"`
	NumResults int           `json:"numResults"`
	NumPages   int           `json:"numPages"`
	Page       int           `json:"page"`
	Pagesize   int           `json:"pagesize"`
	Seid       string        `json:"seid"`
	Result     []SingeResult `json:"result"`
}

func NewCompareResult(projectId int, projectName string, saleFlag string, saleEndTime int64, saleStartTime int64, venueId string, venueName string, cover string, url string, staff string, starttime string, endtime string) *CompareResult {
	return &CompareResult{ProjectId: projectId, ProjectName: projectName, SaleFlag: saleFlag, SaleEndTime: saleEndTime, SaleStartTime: saleStartTime, VenueId: venueId, VenueName: venueName, Cover: cover, Url: url, Staff: staff, StartTime: starttime, EndTime: endtime}
}

// sale_flag_number 2为正常卖中 1是没开票
type SingeResult struct {
	Banner         string      `json:"banner,omitempty"`
	City           string      `json:"city"`
	CityId         int         `json:"city_id"`
	CityName       string      `json:"city_name"`
	Countdown      string      `json:"countdown"`
	Cover          string      `json:"cover"`
	EndTime        string      `json:"end_time"`
	Id             int         `json:"id"`
	IsCommend      bool        `json:"is_commend,omitempty"`
	IsFree         bool        `json:"is_free"`
	IsPrice        bool        `json:"is_price"`
	IsRebate       bool        `json:"is_rebate"`
	IsSale         bool        `json:"is_sale,omitempty"`
	NeedUp         int         `json:"need_up,omitempty"`
	PickSeat       bool        `json:"pick_seat"`
	PriceHigh      int         `json:"price_high"`
	PriceLow       int         `json:"price_low"`
	ProjectName    string      `json:"project_name"`
	ProjectType    string      `json:"project_type,omitempty"`
	RankIndex      int         `json:"rank_index,omitempty"`
	RankOffset     int         `json:"rank_offset,omitempty"`
	RequiredNumber int         `json:"required_number,omitempty"`
	SaleEndTime    int64       `json:"sale_end_time"`
	SaleStartTime  int64       `json:"sale_start_time"`
	SaleFlag       string      `json:"sale_flag"`
	SaleFlagNumber int         `json:"sale_flag_number"`
	SalePoint      string      `json:"sale_point"`
	ShowTime       string      `json:"show_time,omitempty"`
	StartTime      string      `json:"start_time"`
	StartUnix      int         `json:"start_unix"`
	Tlabel         string      `json:"tlabel"`
	Type           interface{} `json:"type"`
	Url            string      `json:"url"`
	VenueId        string      `json:"venue_id"`
	VenueName      string      `json:"venue_name"`
	Wish           int         `json:"wish"`
	Guests         []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"guests,omitempty"`
	IsExclusive        bool          `json:"is_exclusive"`
	Promo              interface{}   `json:"promo"`
	Label              string        `json:"label,omitempty"`
	Mask               interface{}   `json:"mask"`
	Distance           int           `json:"distance"`
	DistrictName       string        `json:"district_name"`
	Tags               []interface{} `json:"tags"`
	RemindStatus       bool          `json:"remind_status"`
	ShowRemindBtn      bool          `json:"show_remind_btn"`
	CountdownSec       int           `json:"countdown_sec"`
	ProjectQuality     int           `json:"project_quality"`
	ProjectQualityDesc string        `json:"project_quality_desc"`
	ProjectId          int           `json:"project_id,omitempty"`
	Score              int           `json:"score,omitempty"`
	ExtraInfo          string        `json:"extra_info,omitempty"`
	Status             int           `json:"status,omitempty"`
	SaleStart          int           `json:"sale_start,omitempty"`
	Country            int           `json:"country,omitempty"`
	SaleEnd            int           `json:"sale_end,omitempty"`
	Staff              string        `json:"staff,omitempty"`
	Areas              string        `json:"areas,omitempty"`
	CityIdStr          string        `json:"city_id_str,omitempty"`
	UserCount          int           `json:"user_count,omitempty"`
	Province           string        `json:"province,omitempty"`
	Group              string        `json:"group,omitempty"`
	Span               string        `json:"span,omitempty"`
	Ord                float64       `json:"ord,omitempty"`
	StrategyOrd        float64       `json:"strategy_ord,omitempty"`
	Order              int           `json:"order,omitempty"`
	MaxOrder           int           `json:"max_order,omitempty"`
	Gmv                int           `json:"gmv,omitempty"`
	MaxGmv             int           `json:"max_gmv,omitempty"`
	Uv                 int           `json:"uv,omitempty"`
	MaxUv              int           `json:"max_uv,omitempty"`
	Wtg                int           `json:"wtg,omitempty"`
	MaxWtg             int           `json:"max_wtg,omitempty"`
	Comment            int           `json:"comment,omitempty"`
	MaxComment         int           `json:"max_comment,omitempty"`
	PubTime            int           `json:"pub_time,omitempty"`
	HasAct             int           `json:"has_act,omitempty"`
	Commend            int           `json:"commend,omitempty"`
	GrossMargin        int           `json:"gross_margin,omitempty"`
	IsSeckill          bool          `json:"is_seckill,omitempty"`
}

type BaseInfo struct {
	Total      int    `json:"total"`
	NumResults int    `json:"numResults"`
	NumPages   int    `json:"numPages"`
	Page       int    `json:"page"`
	Pagesize   int    `json:"pagesize"`
	Seid       string `json:"seid"`
}

type CompareResult struct {
	ProjectId     int    `json:"project_id,omitempty"`
	ProjectName   string `json:"project_name"`
	SaleFlag      string `json:"sale_flag"`
	SaleEndTime   int64  `json:"sale_end_time"`
	SaleStartTime int64  `json:"sale_start_time"`
	VenueId       string `json:"venue_id"`
	VenueName     string `json:"venue_name"`
	Cover         string `json:"cover"`
	Url           string `json:"url"`
	Staff         string `json:"staff,omitempty"`
	StartTime     string `json:"start_time`
	EndTime       string `json:"end_time"`
}

func (b *BaseInfo) String() string {

	return "stop"
}

func (s *SingeResult) Conv2Com() *CompareResult {
	return NewCompareResult(s.ProjectId, s.ProjectName, s.SaleFlag, s.SaleEndTime, s.SaleStartTime, s.VenueId, s.VenueName, "https:"+s.Cover, s.Url, s.Staff, s.StartTime, s.EndTime)
}

// 按 Unix 时间戳排序的接口实现
type ByUnixTimestamp []SingeResult

func (a ByUnixTimestamp) Len() int           { return len(a) }
func (a ByUnixTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUnixTimestamp) Less(i, j int) bool { return a[i].StartUnix < a[j].StartUnix }

func (v *CompareResult) String() string {
	startStr := time.Unix(v.SaleStartTime, 0).Format(time.DateTime)
	endStr := time.Unix(v.SaleEndTime, 0).Format(time.DateTime)
	var s strings.Builder
	s.WriteString("名称： " + v.ProjectName + "\r\n")
	s.WriteString("售票状态： " + v.SaleFlag + "\r\n")
	s.WriteString("售票开始时间： " + startStr + "\r\n")
	s.WriteString("售票结束时间： " + endStr + "\r\n")
	s.WriteString("地址 ： " + v.VenueName + "\r\n")
	s.WriteString("开始时间" + v.StartTime + "\r\n")
	s.WriteString("结束时间" + v.EndTime + "\r\n")

	return s.String()
}
