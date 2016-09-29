// Code generated by "gen -pkg ."; DO NOT EDIT.

package redp

import (
	"golang.org/x/net/context"
	"ndp/lib/trace"
)

func CreateAct(ctx context.Context, a *Activity) (string, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.CreateAct")
	return createAct(ctx, a)
}

func ActById(ctx context.Context, actId string) (*Activity, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.ActById")
	return actById(ctx, actId)
}

func CreateRec(ctx context.Context, r *Record) (string, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.CreateRec")
	return createRec(ctx, r)
}

func RecById(ctx context.Context, recId string) (*Record, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.RecById")
	return recById(ctx, recId)
}

func RecByUser(ctx context.Context, actId string, uid int64) (*Record, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.RecByUser")
	return recByUser(ctx, actId, uid)
}

func MarkDrawed(ctx context.Context, r *Record) error {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.MarkDrawed")
	return markDrawed(ctx, r)
}

func DrawNum(ctx context.Context, actId string) (int64, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.DrawNum")
	return drawNum(ctx, actId)
}

func GetDrawCount(ctx context.Context, actId string) (int64, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.GetDrawCount")
	return getDrawCount(ctx, actId)
}

func UpdateDrawCount(ctx context.Context, actId string, drawCount int64) error {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.UpdateDrawCount")
	return updateDrawCount(ctx, actId, drawCount)
}

func FriendListOfRecord(ctx context.Context, recId, actId string) ([]Friend, error) {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.FriendListOfRecord")
	return friendListOfRecord(ctx, recId, actId)
}

func IncrActRetCount(ctx context.Context, value int, actId string) error {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.IncrActRetCount")
	return incrActRetCount(ctx, value, actId)
}

func IncrRecRetCount(ctx context.Context, value int, recId string) error {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.IncrRecRetCount")
	return incrRecRetCount(ctx, value, recId)
}

func IncrRecHuaFei(ctx context.Context, huaFei int, recId string) error {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.IncrRecHuaFei")
	return incrRecHuaFei(ctx, huaFei, recId)
}

func IncrRecLiuLiang(ctx context.Context, liuLiang int, recId string) error {
	ctx = trace.WithLables(ctx, "ndp/n3d/wxapp/redp.IncrRecLiuLiang")
	return incrRecLiuLiang(ctx, liuLiang, recId)
}