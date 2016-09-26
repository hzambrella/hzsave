package redp

import "golang.org/x/net/context"

//go:generate ndp_tools_gen -pkg .

// @trace
func createAct(ctx context.Context, a *Activity) (int64, error) {
	return db.CreateAct(ctx, a)
}

// @trace
func actById(ctx context.Context, actId int64) (*Activity, error) {
	return db.ActById(ctx, actId)
}

// @trace
func createRec(ctx context.Context, r *Record) (int64, error) {
	return db.CreateRec(ctx, r)
}

// @trace
func recById(ctx context.Context, recId int64) (*Record, error) {
	return db.RecById(ctx, recId)
}

// @trace
func recByUser(ctx context.Context, actId, uid int64) (*Record, error) {
	return db.RecByUser(ctx, actId, uid)
}

// @trace
func markDrawed(ctx context.Context, r *Record) error {
	return db.MarkDrawed(ctx, r)
}

// @trace
func drawNum(ctx context.Context, actId int64) (int64, error) {
	return db.DrawNum(ctx, actId)
}

// @trace
func getDrawCount(ctx context.Context, actId int64) (int64, error) {
	return db.GetDrawCount(ctx, actId)
}

// @trace
func updateDrawCount(ctx context.Context, actId, drawCount int64) error {
	num, err := db.GetDrawCount(ctx, actId)
	if err != nil {
		return err
	}
	return db.UpdateDrawCount(ctx, actId, num+drawCount)
}

// 检查活动人数是否已满
func (act *Activity) IsOverDraw() bool {
	if act.DrawCount >= act.DrawMax {
		return true
	}
	return false
}

// @trace
func redpDrawUserList(ctx context.Context, actId  int64) ([]string, error) {
	return db.DrawUserList(ctx, actId )
}

// @trace
func getContactMobile(ctx context.Context, userMobile string) ([]string, error) {
	return db.GetContactMobile(ctx, userMobile)
}

// @trace
func friendListOfRecord(ctx context.Context, userId, actId int64) ([]Friend, error) {
	return db.FriendListOfRecord(ctx, userId, actId)
}
