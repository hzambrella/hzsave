package redp

import "time"

// 红包活动配置
type Activity struct {
	// 活动id
	Id string `db:"id" json:"id"`

	// 时间
	CreateAt time.Time `db:"create_at" json:"create_at"`
	UpdateAt time.Time `db:"update_at" json:"update_at"`

	// 状态, 0正常，-1已删除
	Status   int       `db:"status" json:"status"`
	DeleteAt time.Time `db:"delete_at" json:"delete_at"`

	// 所属产品和商户，方便查询和使用
	Pid     string `db:"pid" json:"pid"`
	AgentId int64  `db:"agent_id" json:"agent_id"`

	// 页面上显示的标题
	Title string `db:"title" json:"title"`
	// 话费区间左值
	HuafeiL int `db:"huafei_l" json:"huafei_l"`
	// 话费区间右值
	HuafeiR int `db:"huafei_r" json:"huafei_r"`
	// 积分
	Jifen int `db:"jifen" json:"jifen"`
	// 流量，以M为单位
	Mb int `db:"mb" json:"mb"`

	// 注册页面图片
	RegImgUrl string `db:"reg_img_url" json:"reg_img_url"`
	// 红包页面图片
	PkgImgUrl string `db:"pkg_img_url" json:"pkg_img_url"`
	// 分享页面的二维码
	QrcodeImgUrl string `db:"qrcode_img_url" json:"qrcode_img_url"`

	// 分享图片，一般是logo图片
	ShareImgUrl string `db:"share_img_url" json:"share_img_url"`
	// 分享标题
	ShareTitle string `db:"share_title" json:"share_title"`
	// 分享内容
	ShareContent string `db:"share_content" json:"share_content"`
	// 分享链接
	ShareLinkUrl string `db:"share_link_url" json:"share_link_url"`

	// 返还最多次数
	RetMax int `db:"ret_max" json:"ret_max"`
	// 返还计数
	RetCount int `db:"ret_count" json:"ret_count"`
	// 每返还话费
	RetHuafei int `db:"ret_huafei" json:"ret_huafei"`
	// 每返还积分
	RetJifen int `db:"ret_jifen" json:"ret_jifen"`
	// 返还流量
	RetMb int `db:"ret_mb" json:"ret_mb"`

	// 分享到朋友圈话费
	TlHuafei int `db:"tl_huafei" json:"tl_huafei"`
	// 分享到朋友圈积分
	TlJifen int `db:"tl_jifen" json:"tl_jifen"`
	// 分享到朋友流量
	TlMb int `db:"tl_mb" json:"tl_mb"`

	// 有效时间范围
	ActiveFrom time.Time `db:"active_from" json:"active_from"`
	ActiveTo   time.Time `db:"active_to" json:"active_to"`

	// 领取份数控制
	DrawMax   int `db:"draw_max" json:"draw_max"`
	DrawCount int `db:"draw_count" json:"draw_count"`

	// 是否显示确认框
	Checkbox int `db:"checkbox" json:"checkbox"`
	// 页尾提示文字
	Footer string `db:"footer" json:"footer"`
}
