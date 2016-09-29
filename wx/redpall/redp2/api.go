package redp

import "golang.org/x/net/context"

//go:generate ndp_tools_gen -pkg .

// @trace
func createAct(ctx context.Context, a *Activity) (string, error) {
	return db.CreateAct(ctx, a)
}

// @trace
func actById(ctx context.Context, actId string) (*Activity, error) {
	return db.ActById(ctx, actId)
}

// @trace
func createRec(ctx context.Context, r *Record) (string, error) {
	return db.CreateRec(ctx, r)
}

// @trace
func recById(ctx context.Context, recId string) (*Record, error) {
	return db.RecById(ctx, recId)
}

// @trace
func recByUser(ctx context.Context, actId string, uid int64) (*Record, error) {
	return db.RecByUser(ctx, actId, uid)
}

// @trace
func markDrawed(ctx context.Context, r *Record) error {
	return db.MarkDrawed(ctx, r)
}

// @trace
func drawNum(ctx context.Context, actId string) (int64, error) {
	return db.DrawNum(ctx, actId)
}

// @trace
func getDrawCount(ctx context.Context, actId string) (int64, error) {
	return db.GetDrawCount(ctx, actId)
}

// @trace
func updateDrawCount(ctx context.Context, actId string, drawCount int64) error {
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

// 点击查看好友
// @trace
func friendListOfRecord(ctx context.Context, recId, actId string) ([]Friend, error) {
	return db.FriendListOfRecord(ctx, recId, actId)
}

// @trace
func incrActRetCount(ctx context.Context, value int, actId string) error {
	return db.IncrActRetCount(ctx, value, actId)
}

// @trace
func incrRecRetCount(ctx context.Context, value int, recId string) error {
	return db.IncrRecRetCount(ctx, value, recId)
}

// @trace
func incrRecHuaFei(ctx context.Context, huaFei int, recId string) error {
	return db.IncrRecHuaFei(ctx, huaFei, recId)
}

// @trace
func incrRecLiuLiang(ctx context.Context, liuLiang int, recId string) error {
	return db.IncrRecLiuLiang(ctx,liuLiang, recId)
}

// @trace
func MarkTlShared(ctx context.Context, recId string) error {
	return db.MarkTlShared(ctx, recId)
}

// @trace
func UpdateParentId(ctx context.Context,parentId ,recId,actId string)error{
	return db.UpdateParentId(ctx,parentId,recId,actId)
}
