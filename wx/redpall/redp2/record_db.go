package redp

import (
	"time"

	"ndp/lib/sql"
	"ndp/lib/utils"

	"golang.org/x/net/context"
)

// 创建一条记录
func (db *Redp) CreateRec(ctx context.Context, r *Record) (string, error) {
	id := utils.NewRandom()
	return id, db.Exec(
		ctx, recCreateSql, 1, id, r.ActId, r.ParentId,
		r.UserId, r.UtmSrc, r.UtmId, r.DrawStatus, r.TlStatus,
		r.Huafei, r.RetCount)
}

const recCreateSql = `
INSERT INTO 
	redp_records (
        id, create_at, update_at, act_id, parent_id, 
		user_id, utm_src, utm_id, draw_status, tl_status,  
		huafei, ret_count
    ) 
VALUES (
        ?, NOW(), NOW(), ?, ?,
        ?, ?, ?, ?, ?,
        ?, ?
    )
`

// 根据红包活动id得到活动记录
func (db *Redp) RecById(ctx context.Context, actId string) (*Record, error) {
	var rec Record
	if err := db.QueryRow(ctx, recByIdSql, actId).Scan(
		&rec.Id, &rec.CreateAt, &rec.UpdateAt, &rec.ActId, &rec.ParentId,
		&rec.UserId, &rec.UtmSrc, &rec.UtmId, &rec.DrawStatus, &rec.TlStatus,
		&rec.Huafei, &rec.RetCount,&rec.LiuLiang,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

const recByIdSql = `
SELECT 
	* 
FROM 
	redp_records 
WHERE 
	id = ? 
LIMIT 1    
`

// 根据用户得到红包记录
func (db *Redp) RecByUser(ctx context.Context, actId string, uid int64) (*Record, error) {
	var recId string
	if err := db.QueryRow(ctx, recByUserSql, actId, uid).Scan(&recId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return db.RecById(ctx, recId)
}

const recByUserSql = `
SELECT 
	id
FROM
	redp_records
WHERE 
	act_id = ?
AND
	user_id = ?
LIMIT 1
`

// 标记已领取
func (db *Redp) MarkDrawed(ctx context.Context, r *Record) error {
	return db.Exec(
		ctx, markDrawedSql, 0, REC_DRAW_YES, r.Id,
	)
}

const markDrawedSql = `
UPDATE
	redp_records
SET
	draw_status = ?
WHERE
	id = ?
`

// 计算活动领取人数
func (db *Redp) DrawNum(ctx context.Context, actId string) (int64, error) {
	var num int64
	if err := db.QueryRow(ctx, findDrawNumSql, actId).Scan(&num); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, nil
}

const (
	findDrawNumSql = `
SELECT
	COUNT(id) 
FROM 
	redp_records	
WHERE
	act_id=?
AND
	draw_status = 1
	`
)

// 点击查看好友
func (db *Redp) FriendListOfRecord(ctx context.Context, recId, actId string) ([]Friend, error) {
	rows, err := db.Query(ctx, friendListOfRecordSql, recId, actId)
	if err != nil {
		return nil, err
	}
	defer sql.Close(rows)

	results := []Friend{}
	var times time.Time
	for rows.Next() {
		result := Friend{}
		err := rows.Scan(
			&result.HuaFei,
			&result.LiuLiang,
			&times,
			&result.Uid,
		)
		if err != nil {
			return nil, err
		}

		result.RedpTime = times.Format("01-02 15:04")
		results = append(results, result)
	}
	return results, nil
}

const (
	friendListOfRecordSql = `
SELECT
	huafei,liuliang,update_at,user_id
FROM
	redp_records
WHERE
	parent_id=?
AND
    act_id=?
AND
	draw_status = 1
AND
(
   liuliang!=0
OR
   huafei!=0
)
`
)

// 增加返还数量
func (db *Redp) IncrRecRetCount(ctx context.Context, value int, recId string) error {
	return db.Exec(ctx, "UPDATE redp_records SET ret_count = ret_count + ? WHERE id = ?", 0, value, recId)
}

// 写入红包话费金额
func (db *Redp) IncrRecHuaFei(ctx context.Context, HuaFei int, recId string) error {
	return db.Exec(ctx, "UPDATE redp_records SET huafei=huafei+? where id= ?", 0, HuaFei, recId)
}

// 写入红包流量
func (db *Redp) IncrRecLiuLiang(ctx context.Context, LiuLiang int, recId string) error {
	return db.Exec(ctx, "UPDATE redp_records SET liuliang=liuliang+? where id= ?", 0, LiuLiang, recId)
}

// 标记朋友圈分享状态
func (db *Redp) MarkTlShared(ctx context.Context, recId string) error {
	return db.Exec(ctx, "UPDATE redp_records SET tl_status = ? WHERE id = ?", 0, REC_TL_YES, recId)
}

// 如果有记录，更新parentid,抢完红包后,点击其它分享者的链接,就不能更换parentId
func (db *Redp)UpdateParentId(ctx context.Context,parentId,recId,actId string)error{
	return db.Exec(ctx, "UPDATE redp_records SET parent_id= ? WHERE act_id = ? AND id=? AND draw_status=0",0,parentId,actId, recId)
}
