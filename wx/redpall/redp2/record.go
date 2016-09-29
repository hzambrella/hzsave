package redp

import (
	"errors"
	"time"
)

// 领取状态
const (
	// 未领取
	REC_DRAW_NO = 0

	// 已领取
	REC_DRAW_YES = 1
)

// 分享朋友圈状态
const (
	// 未分享
	REC_TL_NO = 0

	// 已分享
	REC_TL_YES = 1
)

// 错误定义
var (
	// 红包不存在
	ErrRecNotFound = errors.New("redp: record not found")
)

//
// 用户的红包记录
//
type Record struct {
	// 红包id
	Id string `db:"id" json:"id"`

	// 创建时间
	CreateAt time.Time `db:"create_at" json:"create_at"`
	// 更新时间
	UpdateAt time.Time `db:"update_at" json:"update_at"`

	// 所属商户，方便查询，而且也不会变更
	ActId string `db:"act_id" json:"act_id"`
	// 源分享RecId, 不是uid
	ParentId string `db:"parent_id" json:"parent_id"`
	// 用户ID
	UserId int64 `db:"user_id" json:"user_id"`

	// 来源
	UtmSrc string `db:"utm_src" json:"utm_src"`
	UtmId  string `db:"utm_id" json:"utm_id"`

	// 领取状态
	DrawStatus int `db:"draw_status" json:"draw_status"`
	// 分享到朋友圈状态
	TlStatus int `db:"tl_status" json:"tl_status"`

	// 已领取话费
	Huafei int `db:"huafei" json:"huafei"`
	// 返还计数
	RetCount int `db:"ret_count" json:"ret_count"`
	LiuLiang int`db:"liuliang" json:"liuliang"`
}

type Friend struct {
	HuaFei   int
	LiuLiang int
	RedpTime string
	Uid      string
}
