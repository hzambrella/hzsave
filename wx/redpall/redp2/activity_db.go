package redp

import (
	"golang.org/x/net/context"

	"ndp/lib/etc"
	"ndp/lib/log"
	"ndp/lib/sql"
	"ndp/lib/utils"
)

var db = New()

type Redp struct {
	*sql.DB
}

func New() *Redp {
	db, err := sql.Default(etc.Str("dsn", "n3d_wxapp"))
	if err != nil {
		log.Fatal("err", err.Error())
	}
	return &Redp{db}
}

// 创建一个红包活动
func (db *Redp) CreateAct(ctx context.Context, a *Activity) (string, error) {
	id := utils.NewRandom()
	return id, db.Exec(
		ctx, actCreateSql, 1, id, a.Status, a.Pid,
		a.AgentId, a.Title, a.HuafeiL, a.HuafeiR, a.Jifen,
		a.Mb, a.RegImgUrl, a.PkgImgUrl, a.QrcodeImgUrl, a.ShareImgUrl,
		a.ShareTitle, a.ShareContent, a.RetMax, a.RetCount, a.RetHuafei,
		a.RetJifen, a.RetMb, a.TlHuafei, a.TlJifen, a.TlMb,
		a.ActiveFrom, a.ActiveTo, a.DrawMax, a.DrawCount, a.Checkbox,
		a.Footer)
}

const actCreateSql = `
INSERT INTO 
	redp_activities (
		id, create_at, update_at, status, delete_at,
        pid, agent_id, title, huafei_l, huafei_r, 
		jifen, mb, reg_img_url, pkg_img_url, qrcode_img_url,
        share_img_url, share_title, share_content,ret_max, ret_count, 
		ret_huafei, ret_jifen, ret_mb, tl_huafei, tl_jifen, 
		tl_mb, active_from, active_to, draw_max, draw_count,
        checkbox, footer
    ) 
VALUES (
        ?, NOW(), NOW(), ?, NOW(),
        ?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?,
        ?, ?
    )
`

// 根据红包活动id得到红包活动
func (db *Redp) ActById(ctx context.Context, actId string) (*Activity, error) {
	var act Activity
	if err := db.QueryRow(ctx, actByIdSql, actId).Scan(
		&act.Id, &act.CreateAt, &act.UpdateAt, &act.Status, &act.DeleteAt,
		&act.Pid, &act.AgentId, &act.Title, &act.HuafeiL, &act.HuafeiR,
		&act.Jifen, &act.Mb, &act.RegImgUrl, &act.PkgImgUrl, &act.QrcodeImgUrl,
		&act.ShareImgUrl, &act.ShareTitle, &act.ShareContent, &act.RetMax, &act.RetCount,
		&act.RetHuafei, &act.RetJifen, &act.RetMb, &act.TlHuafei, &act.TlJifen,
		&act.TlMb, &act.ActiveFrom, &act.ActiveTo, &act.DrawMax, &act.DrawCount,
		&act.Checkbox, &act.Footer, &act.ShareLinkUrl,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &act, nil
}

const actByIdSql = `
SELECT 
	* 
FROM 
	redp_activities
WHERE id = ? 
    LIMIT 1    
`

// 获取活动领取人数
func (db *Redp) GetDrawCount(ctx context.Context, actId string) (int64, error) {
	var num int64
	if err := db.QueryRow(ctx, getDrawCountSql, actId).Scan(&num); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, nil
}

const (
	getDrawCountSql = `
SELECT
	draw_count 
FROM 
	redp_activities	
WHERE
	id=?
`
)

// 更新活动参与人数
func (db *Redp) UpdateDrawCount(ctx context.Context, actId string, drawCount int64) error {
	return db.Exec(
		ctx, updateDrawCountSql, 0, drawCount, actId)
}

const (
	updateDrawCountSql = `
UPDATE
	redp_activities
SET
	draw_count=?
WHERE
	id=?
`
)

// 增加返还数量
func (db *Redp) IncrActRetCount(ctx context.Context, value int, actId string) error {
	return db.Exec(ctx, "UPDATE redp_activities SET ret_count = ret_count + ? WHERE id = ?", 0, value, actId)
}

func (db *Redp) GetShareInfo(ctx context.Context, actId string) {

}
