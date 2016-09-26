package record

import (
	"time"

	"golang.org/x/net/context"

	"ndp/lib/etc"
	"ndp/lib/log"
	"ndp/lib/sql"
)

var (
	db = New()
)

type Record struct {
	*sql.DB
}

type RecordList struct {
	Id        int
	Caller    string
	Callee    string
	NickName  string
	BeginTime string
}

func New() *Record {
	db, err := sql.Default(etc.Str("dsn", "n3d_wxapp"))
	if err != nil {
		log.Fatal("err", err.Error())
	}
	return &Record{db}
}

func (db *Record) AddRecord(ctx context.Context, caller, callee string) error {
	beginTime := time.Now()
	return db.Exec(ctx, addRecordSql, 1, caller, callee, beginTime)
}

const (
	addRecordSql = `
INSERT INTO
     call_records
     (caller,callee,begin_time)
VALUES
     (?,?,?)
	`
)

func (db *Record) GetRecordList(ctx context.Context, caller string) ([]RecordList, error) {
	rows, err := db.Query(ctx, getRecordListSql, caller)
	if err != nil {
		return nil, err
	}
	defer sql.Close(rows)

	results := []RecordList{}
	var times time.Time
	for rows.Next() {
		result := RecordList{}
		err := rows.Scan(
			&result.Id,
			&result.Callee,
			&times,
		)
		if err != nil {
			return nil, err
		}

		result.NickName, err = db.GetCalleeName(ctx, caller, result.Callee)
		if err != nil {
			return nil, err
		}
		result.BeginTime = times.Format("01-02 15:04")
		result.Caller = caller
		results = append(results, result)
	}
	return results, nil
}

const (
	getRecordListSql = `
SELECT
	id,callee,begin_time
FROM
	call_records
WHERE
	caller=?
ORDER BY
	begin_time
DESC
`
)

func (db *Record) GetCalleeName(ctx context.Context, caller, callee string) (string, error) {
	var nickName string
	if err := db.QueryRow(ctx, getCalleeName, caller, callee).Scan(&nickName); err != nil {
		if err == sql.ErrNoRows {
			return callee, nil
		}
		return "", err
	}
	return nickName, nil
}

const (
	getCalleeName = `
SELECT
   contact_nickname
FROM
   call_contact
WHERE
   user_mobile=?
AND
   contact_mobile=?
	`
)

func (db *Record) DeleteRecord(ctx context.Context, id int) error {
	return db.Exec(ctx, deleteRecordSql, 1, id)
}

const (
	deleteRecordSql = `
DELETE FROM
     call_records
WHERE
     id=?
`
)
