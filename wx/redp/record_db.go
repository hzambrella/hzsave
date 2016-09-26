package redp

import (
	"time"

	"ndp/lib/sql"

	"golang.org/x/net/context"
)

// 创建一条记录
func (db *Redp) CreateRec(ctx context.Context, r *Record) (int64, error) {
	return db.ExecX(
		ctx, recCreateSql, 1, r.ActId, r.ParentId,
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
        NULL, NOW(), NOW(), ?, ?,
        ?, ?, ?, ?, ?,
        ?, ?
    )
`

// 根据红包活动id得到活动记录
func (db *Redp) RecById(ctx context.Context, actId int64) (*Record, error) {
	var rec Record
	if err := db.QueryRow(ctx, recByIdSql, actId).Scan(
		&rec.Id, &rec.CreateAt, &rec.UpdateAt, &rec.ActId, &rec.ParentId,
		&rec.UserId, &rec.UtmSrc, &rec.UtmId, &rec.DrawStatus, &rec.TlStatus,
		&rec.Huafei, &rec.RetCount,
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
func (db *Redp) RecByUser(ctx context.Context, actId, uid int64) (*Record, error) {
	var recId int64
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
func (db *Redp) DrawNum(ctx context.Context, actId int64) (int64, error) {
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

//
func (db *Redp) DrawUserList(ctx context.Context, actId int64) ([]string, error) {
	rows, err := db.Query(ctx, findDrawUserListSql, actId)
	if err != nil {
		return nil, err
	}

	defer sql.Close(rows)
	var res string
	userList := []string{}
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			return nil, err
		}
		userList = append(userList, res)
	}
	return userList, nil
}

const (
	findDrawUserListSql = `
SELECT
    user_id
FROM
	redp_records
WHERE
    act_id=?
AND
	draw_status = 1
`
)

func (db *Redp) GetContactMobile(ctx context.Context, userMobile string) ([]string, error) {
	rows, err := db.Query(ctx, getContactMobileSql, userMobile)
	if err != nil {
		return nil, err
	}

	defer sql.Close(rows)
	var res string
	contactList := []string{}
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			return nil, err
		}
		contactList = append(contactList, res)
	}
	return contactList, nil

}

const (
	getContactMobileSql = `
SELECT
   contact_mobile
FROM
   call_contact
WHERE
   user_mobile=?
`
)

func (db *Redp) FriendListOfRecord(ctx context.Context, userId, actId int64) ([]Friend, error) {
	rows, err := db.Query(ctx, friendListOfRecordSql, userId, actId)
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
			&times,
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
	huafei,create_at
FROM
	redp_records
WHERE
	user_id=?
AND
    act_id=?
`
)
