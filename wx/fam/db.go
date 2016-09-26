package family

import (
	"golang.org/x/net/context"

	"ndp/lib/etc"
	"ndp/lib/log"
	"ndp/lib/sql"
)

var dc = New()

type BigNum struct {
	*sql.DB
}

func New() *BigNum {
	db, err := sql.Default(etc.Str("dsn", "n2d_center"))
	if err != nil {
		log.Fatal("err", err.Error())
	}
	return &BigNum{db}
}

// 给出用户的绑定资源数
func (c *BigNum) getBindNumOfUser(ctx context.Context, uid uint64, pid string) (int, error) {
	var bindNumCount int
	if err := c.QueryRow(ctx, getBindNumOfCountSql, uid, pid).Scan(&bindNumCount); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return bindNumCount, nil
}

const (
	getBindNumOfCountSql = `
SELECT
   count(tb2.id)
FROM
   call_bignum_bind tb1
INNER JOIN
   call_bignum tb2
ON
   tb1.bignum_id = tb2.id
AND
   tb2.status = 0
WHERE
   tb1.userid = ?
AND
   tb1.pid=?
AND
   tb1.status = 0
	`
)

// 获取给当前用户分配的大号ID
func (c *BigNum) getCurrentBigNum(ctx context.Context, caller string) (int, error) {
	var list []int

	rows, err := c.Query(ctx, getBigNumIdSql, caller)
	if err != nil {
		return 0, err
	}
	defer sql.Close(rows)

	var id int
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		list = append(list, id)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	// 如果所有大号都已经被该用户绑定了
	if len(list) == 0 {
		return 0, nil
	}
	return list[0], nil
}

// 获取给当前用户分配的大号ID
const (
	getBigNumIdSql = `
SELECT
     tb1.id
FROM
     call_bignum tb1
LEFT JOIN
     call_bignum_bind tb2
ON
     tb1.id=tb2.bignum_id
AND
     caller=?
AND
     tb2.status=0
WHERE
     caller is null
AND
     tb1.status=0
ORDER BY
     tb1.id
LIMIT 1
	`
)

// 更新绑定的大号Id(更新bignum_id为新绑定的bignum_id)
func (c *BigNum) updateBindBigNumId(ctx context.Context, uid, version uint64, pid, callee string, bigNumId int) error {
	return c.Exec(ctx, updateBindBigNumIdSql, 1, bigNumId, version, pid, uid, callee)
}

// 更新绑定的大号Id(更新bignum_id为新绑定的bignum_id)
const (
	updateBindBigNumIdSql = `
UPDATE
   call_bignum_bind
SET
   bignum_id=?,status=0,version=?
WHERE
   pid=?
AND
   userid=?
AND
   callee=?
   `
)

// 新的绑定被叫-大号数据
func (c *BigNum) bindData(ctx context.Context, uid, version uint64, pid, caller, callee string, bindBigNumId int) error {
	return c.Exec(ctx, bindDataSql, 1,
		pid,
		uid,
		caller,
		callee,
		bindBigNumId,
		version,
	)
}

// 新的绑定被叫-大号数据
const (
	bindDataSql = `
INSERT INTO
    call_bignum_bind(
	pid, userid, caller, callee, bignum_id,
	call_uri, status, version, create_time
    )
VALUES
    (
	?,?,?,?,?,"",0,?,CURRENT_TIMESTAMP()
    )
	`
)

// 在已绑定过大号的情况下修改亲情号码
func (c *BigNum) updateBigNumId(ctx context.Context, uid, version uint64, pid, callee string, bigNumId int) error {
	return c.Exec(ctx, updateBigNumIdSql, 1, callee, version, uid, bigNumId, pid)
}

const (
	updateBigNumIdSql = `
UPDATE
	 call_bignum_bind
SET
	 callee=?,version=?
WHERE
	 userid=?
AND
     bignum_id=?
AND
     pid=?
AND
     status=0
	`
)

// 查询某大号是否被该用户使用
func (c *BigNum) checkBigNumId(ctx context.Context, pid string, uid uint64, bigNumId int) ([]string, error) {
	var list []string

	rows, err := c.Query(ctx, checkBigNumIdSql, pid, uid, bigNumId)
	if err != nil {
		return nil, err
	}
	defer sql.Close(rows)
	var callee string
	for rows.Next() {
		if err := rows.Scan(&callee); err != nil {
			return nil, err
		}
		list = append(list, callee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

const (
	checkBigNumIdSql = `
SELECT
    callee
FROM
    call_bignum_bind
WHERE
    pid=?
AND
    userid=?
AND
    bignum_id=?
AND
   status=0
	`
)

// 查询某被叫号码的大号ID
func (c *BigNum) getNumId(ctx context.Context, pid string, uid uint64, callee string) ([]int, error) {
	var list []int

	rows, err := c.Query(ctx, getNumIdSql, pid, uid, callee)
	if err != nil {
		return nil, err
	}
	defer sql.Close(rows)

	var bignumId int
	for rows.Next() {
		if err := rows.Scan(&bignumId); err != nil {
			return nil, err
		}
		list = append(list, bignumId)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// 查询某被叫号码的大号ID
const (
	getNumIdSql = `
SELECT
    bignum_id
FROM
    call_bignum_bind
WHERE
    pid=?
AND
    userid=?
AND
    callee=?
	`
)

// 查看大号对应的大号ID
func (c *BigNum) bigNumOfId(ctx context.Context, num string) (int, error) {
	var bigNumId int
	if err := c.QueryRow(ctx, bigNumOfIdSql, num).Scan(&bigNumId); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return bigNumId, nil
}

const (
	bigNumOfIdSql = `
SELECT
	id
FROM
	call_bignum
WHERE
    num=?
AND
	status=0
`
)

// 查看有多少个大号资源
func (c *BigNum) getMaxBigNum(ctx context.Context) (string, error) {
	var maxnum string
	if err := c.QueryRow(ctx, maxBigNumSql).Scan(&maxnum); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return maxnum, nil
}

// 查看有多少个大号资源
const (
	maxBigNumSql = `
SELECT
	count(*)
FROM
	call_bignum
WHERE
	status=0
`
)

// 查询亲情号码
func (c *BigNum) GetCallee(ctx context.Context, uid uint64, pid string) ([]string, error) {
	rows, err := c.Query(ctx, getCalleeSql, uid, pid)
	if err != nil {
		return nil, err
	}
	defer sql.Close(rows)

	result := []string{}

	for rows.Next() {
		var callee string
		if err := rows.Scan(&callee); err != nil {
			return nil, err
		}
		result = append(result, callee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// 查询用户绑定的亲情号码
const (
	getCalleeSql = `
SELECT
	tb1.callee
FROM
    call_bignum_bind tb1
INNER JOIN
    call_bignum tb2
ON
    tb1.bignum_id = tb2.id
AND
    tb2.status = 0
WHERE
	tb1.userid = ?
AND
    tb1.pid = ?
AND
    tb1.status = 0
ORDER BY
    tb1.version
`
)

// 取消亲情号码
func (c *BigNum) deleteCallee(ctx context.Context, uid, version uint64, pid, callee string) error {
	return c.Exec(ctx, updateBindStatusSql, 0, version, uid, callee, pid)
}

// 将绑定号码的状态变为1
const (
	updateBindStatusSql = `
UPDATE
	call_bignum_bind
SET
    status=1,version=?
WHERE
    userid=?
AND
    callee=?
AND
    pid=?
	`
)

// 用大号取消亲情号码
func (c *BigNum) deleteByBigNum(ctx context.Context, uid, version uint64, pid string, bigNumId int) error {
	return c.Exec(ctx, deleteByBigNumSql, 1, version, uid, bigNumId, pid)
}

// 将绑定号码的状态变为1
const (
	deleteByBigNumSql = `
UPDATE
	call_bignum_bind
SET
    status=1,version=?
WHERE
    userid=?
AND
    bignum_id=?
AND
    pid=?
AND
   status=0
	`
)
