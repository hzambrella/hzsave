package redp

import (
	"golang.org/x/net/context"

	"ndp/lib/etc"
	"ndp/lib/log"
	"ndp/lib/sql"
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
func (db *Redp) CreateAct(ctx context.Context, a *Activity) (int64, error) {
	return db.ExecX(
		ctx, actCreateSql, 1, a.Status, a.Pid,
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
        NULL, NOW(), NOW(), ?, NOW(),
        ?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?,
        ?, ?
    )
`

// 根据红包活动id得到红包活动
func (db *Redp) ActById(ctx context.Context, actId int64) (*Activity, error) {
	var act Activity
	if err := db.QueryRow(ctx, actByIdSql, actId).Scan(
		&act.Id, &act.CreateAt, &act.UpdateAt, &act.Status, &act.DeleteAt,
		&act.Pid, &act.AgentId, &act.Title, &act.HuafeiL, &act.HuafeiR,
		&act.Jifen, &act.Mb, &act.RegImgUrl, &act.PkgImgUrl, &act.QrcodeImgUrl,
		&act.ShareImgUrl, &act.ShareTitle, &act.ShareContent, &act.RetMax, &act.RetCount,
		&act.RetHuafei, &act.RetJifen, &act.RetMb, &act.TlHuafei, &act.TlJifen,
		&act.TlMb, &act.ActiveFrom, &act.ActiveTo, &act.DrawMax, &act.DrawCount,
		&act.Checkbox, &act.Footer,
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
func (db *Redp) GetDrawCount(ctx context.Context, actId int64) (int64, error) {
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
func (db *Redp) UpdateDrawCount(ctx context.Context, actId, drawCount int64) error {
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
