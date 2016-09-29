package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"ndp/core/account"
	"ndp/core/dtplan"
	"ndp/core/token"
	"ndp/lib/mux"
	"ndp/n3d/wxapp/redp"
	"ndp/n3d/wxapp/wxconfig"
)

// 红包活动入口
// @handler GET /redp/activity/:actId
func redpActivity(ctx *mux.Ctx, actId string) {
	//actId = 1
	//得到用户数据
	userid := token.CurrentUser(ctx)
	uid, err := strconv.Atoi(userid)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 加载红包活动
	act, err := redp.ActById(ctx, actId)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 创建视图数据
	data := make(map[string]interface{})
	data["pkgImgUrl"] = act.PkgImgUrl
	data["actId"] = act.Id
	var parentExsist = false
	data["p_exsist"] = parentExsist
	// 红包活动当前是否进行着
	var timeFlag = false
	now := time.Now()
	switch {
	case now.Before(act.ActiveFrom):
		timeFlag = true
		data["active_warn"] = timeFlag
	case now.After(act.ActiveTo):
		timeFlag = true
		data["active_warn"] = timeFlag
	}

	// 检查用户是否已经点击过该红包活动
	rec, err := redp.RecByUser(ctx, act.Id, int64(uid))
	if err != nil {
		ctx.Error(err)
		return
	}
	// 如果没有则创建红包记录
	if rec == nil {
		// 分享者Id 如何获取?
		var parentId string
		// 插入数据库
		rec = &redp.Record{
			ActId:    act.Id,
			ParentId: parentId,
			UserId:   int64(uid),
		}
		id, err := redp.CreateRec(ctx, rec)
		if err != nil {
			ctx.Error(err)
			return
		}
		rec.Id = id
	}

	// 检查活动是否用完
	act, err = redp.ActById(ctx, rec.ActId)
	if err != nil {
		ctx.Error(err)
		return
	}
	// 红包参与人数是否达到上限
	if act.IsOverDraw() {
		ctx.Error(err)
		data["recOverdrawed"] = true
	}

	data["recDrawPath"] = fmt.Sprintf("/redp/record/%d/draw", rec.Id)
	if rec.DrawStatus == redp.REC_DRAW_YES {
		data["rec_drawed"] = true
	}

	num, err := redp.DrawNum(ctx, actId)
	if err != nil {
		ctx.Error(err)
		return
	}
	data["drawNums"] = fmt.Sprint(num)

	data["recId"] = rec.Id

	userTemplate(ctx, "redp/redp_index", data)
}

// 红包收到分享入口
// @handler GET /redp/activity/:actId/:parentid
func redpActivityShare(ctx *mux.Ctx, actId string, parentid string) {
	//actId = 1
	//得到用户数据
	parentRec, err := redp.RecById(ctx, parentid)
	if err != nil {
		ctx.Error(err)
		return
	}
	parent, err := account.WechatInfo(ctx, strconv.FormatInt(parentRec.UserId, 10))
	if err != nil {
		ctx.Error(err)
		return
	}

	userid := token.CurrentUser(ctx)
	uid, err := strconv.Atoi(userid)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 加载红包活动
	act, err := redp.ActById(ctx, actId)
	if err != nil {
		ctx.Error(err)
		return
	}

	var parentExsist = true
	// 创建视图数据
	data := make(map[string]interface{})
	data["pkgImgUrl"] = act.PkgImgUrl
	data["parent_name"] = parent.NickName
	data["parent_img"] = parent.HeadImgurl
	data["p_exsist"] = parentExsist
	data["actId"] = act.Id
	// 红包活动当前是否进行着
	var timeFlag = false
	now := time.Now()
	switch {
	case now.Before(act.ActiveFrom):
		timeFlag = true
		data["active_warn"] = timeFlag
	case now.After(act.ActiveTo):
		timeFlag = true
		data["active_warn"] = timeFlag
	}

	// 检查用户是否已经点击过该红包活动
	rec, err := redp.RecByUser(ctx, act.Id, int64(uid))
	if err != nil {
		ctx.Error(err)
		return
	}
	// 如果没有则创建红包记录
	if rec == nil {
		// 插入数据库
		rec = &redp.Record{
			ActId:    act.Id,
			ParentId: parentid,
			UserId:   int64(uid),
		}
		id, err := redp.CreateRec(ctx, rec)
		if err != nil {
			ctx.Error(err)
			return
		}
		rec.Id = id
	} else {
		// 如果有记录，更新parentid
		err = redp.UpdateParentId(ctx, parentid, rec.Id, act.Id)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
	act, err = redp.ActById(ctx, rec.ActId)
	if err != nil {
		ctx.Error(err)
		return
	}
	// 红包参与人数是否达到上限
	if act.IsOverDraw() {
		ctx.Error(err)
		data["recOverdrawed"] = true
	}

	data["recDrawPath"] = fmt.Sprintf("/redp/record/%d/draw", rec.Id)
	if rec.DrawStatus == redp.REC_DRAW_YES {
		data["rec_drawed"] = true
	}

	num, err := redp.DrawNum(ctx, actId)
	if err != nil {
		ctx.Error(err)
		return
	}
	data["drawNums"] = fmt.Sprint(num)

	data["recId"] = rec.Id

	userTemplate(ctx, "redp/redp_index", data)
}

// 抢红包
// @handler GET /redp/record/draw/:recId/:mobile
func redpRecDraw(ctx *mux.Ctx, recId, mobile string) {
	//得到用户数据
	uid := token.CurrentUser(ctx)
	user, err := account.UserByID(ctx, uid)
	if err != nil {
		ctx.Error(err)
		return
	}

	if user.Mobile == "" {
		user.Mobile = mobile
	}

	if user.Mobile != mobile {
		err := errors.New("请输入自己的手机号")
		ctx.Error(err)
		return
	}

	// 加载record
	rec, err := redp.RecById(ctx, recId)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 加载act
	act, err := redp.ActById(ctx, rec.ActId)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 检查红包是否已经领取过
	if rec.DrawStatus == redp.REC_DRAW_YES {
		ctx.Error(err)
		return
	}

	// 领取红包
	if err := redp.MarkDrawed(ctx, rec); err != nil {
		ctx.Error(err)
		return
	}

	tagId := "1"
	pid := "testing"

	// 页面类型
	// 1--话费
	// 2--流量
	// 3--无
	redpType := rand.Intn(2) + 2
	redpAmount := rand.Intn(50-1) + 1
	cfg := &wxconfig.Config{
		AppId:  "wxe668d7df1101ee6a",
		Secret: "66e7ee651a90a9653cb19725b3a7dd66",
	}

	currUid, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	currRec, err := redp.RecByUser(ctx, act.Id, currUid)
	if err != nil {
		ctx.Error(err)
		return
	}

	shareUrlTmp := act.ShareLinkUrl
	shareUrl := shareUrlTmp + "/" + currRec.Id + "?app_id=" + cfg.AppId + "&tag=" + tagId + "&pid=" + pid
	//currUrl := "http://" + shareUrlTmp
	JsApiSign, err := wxconfig.NewJsApiSign(cfg, string(ctx.Referer()))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.WriteKeyvals(
		"redpType", redpType,
		"redpAmount", redpAmount,
		"appId", cfg.AppId,
		"timestamp", JsApiSign.Timestamp,
		"nonceStr", JsApiSign.NonceStr,
		"Signature", JsApiSign.Signature,
		"shareImgUrl", act.ShareImgUrl,
		"shareTitle", act.ShareTitle,
		"shareContent", act.ShareContent,
		"shareLink", shareUrl,
		"qrCodeUrl", act.QrcodeImgUrl,
	)

	// 领取红包后加钱或流量
	switch redpType {

	case 2:
		// TODO:加钱
		if err := redp.IncrRecHuaFei(ctx, redpAmount, recId); err != nil {
			ctx.Error(err)
			return
		}

		err = redp.UpdateDrawCount(ctx, rec.ActId, int64(1))
		if err != nil {
			ctx.Error(err)
			return
		}

	case 3:
		if err := redp.IncrRecLiuLiang(ctx, redpAmount, recId); err != nil {
			ctx.Error(err)
			return
		}

		err := dtplan.Adjust(ctx, uid, strconv.Itoa(redpAmount), "微信红包", "100", "test")
		if err != nil {
			ctx.Error(err)
			return
		}
		// draw_count加1
		err = redp.UpdateDrawCount(ctx, rec.ActId, int64(1))
		if err != nil {
			ctx.Error(err)
			return
		}
	default:

	}
	// 如果是从分享链接进来的，那么parentId不为空
	if rec.ParentId != "" {
		// 由parentId得到分享者
		prec, err := redp.RecById(ctx, rec.ParentId)
		if err != nil {
			ctx.Error(err)
		}

		// 得到抽取的红包
		act, err := redp.ActById(ctx, rec.ActId)
		if err != nil {
			ctx.Error(err)
		}

		// 增加红包总返还计数
		if err := redp.IncrActRetCount(ctx, 1, act.Id); err != nil {
			ctx.Error(err)
			return
		}

		// 增加分享者的返还计数
		if err := redp.IncrRecRetCount(ctx, 1, prec.Id); err != nil {
			ctx.Error(err)
			return
		}
		// 如果分享次数少于限定的，则返回给分享者
		if prec.RetCount < act.RetMax {
			// 如何操作 使查询的时候会有记录
			// TODO 奖励用户
		}
	}
}

// 活动规则页面
// @handler GET /redp/rule
func redpRule(ctx *mux.Ctx) {
	userTemplate(ctx, "redp/redp_rules", nil)
}

type FriendList struct {
	RedpTime   string
	Prize      int
	Describe   string
	NickName   string
	HeadImgUrl string
}

// 点击查看好友
// @handler GET /redp/friendlist/:recId
func redpFriendList(ctx *mux.Ctx, recId string) {
	record, err := redp.RecById(ctx, recId)
	if err != nil {
		ctx.Error(err)
		return
	}
	actId := record.ActId

	friendRecord := []FriendList{}
	var fre FriendList
	friend, err := redp.FriendListOfRecord(ctx, recId, actId)
	if err != nil {
		ctx.Error(err)
		return
	}
	for _, v := range friend {
		userWechat, err := account.WechatInfo(ctx, v.Uid)
		if err != nil {
			ctx.Error(err)
			return
		}
		fre.NickName = userWechat.NickName
		fre.HeadImgUrl = userWechat.HeadImgurl
		if v.LiuLiang != 0 && v.HuaFei == 0 {
			fre.Prize = v.LiuLiang
			fre.Describe = "M流量"
		} else {
			fre.Prize = v.HuaFei
			fre.Describe = "元话费"
		}

		friendRecord = append(friendRecord, fre)
	}
	userTemplate(ctx, "redp/redp_list", map[string]interface{}{
		"friendRecord": friendRecord,
	})
}

// 分享到朋友圈处理
// @handler GET /redp/share/timeline/:recId
func shareTimeline(ctx *mux.Ctx, recId string) {
	//rec, err := redp.RecById(ctx, recId)
	//if err != nil {
	//	ctx.Error(err)
	//	return
	//}

	// 分享朋友圈状态
	//ShareStatus := false
	//if rec.TlStatus == redp.REC_TL_YES {
	//	ShareStatus = true
	//}

	// TODO 给用户分享到朋友圈的奖励

	// 标记分享到朋友圈
	if err := redp.MarkTlShared(ctx, recId); err != nil {
		ctx.Error(err)
		return
	}

	//	ctx.WriteKeyvals(
	//		"shareStatus", ShareStatus,
	//	)
	ctx.Ok()
}
