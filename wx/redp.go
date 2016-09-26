package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"ndp/core/account"
	"ndp/core/dtplan"
	"ndp/core/token"
	"ndp/lib/mux"
	"ndp/n3d/wxapp/redp"
	"strconv"
)

// 红包活动入口
// @handler GET /redp/activity/:actId
func redpActivity(ctx *mux.Ctx, actId int64) {
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
		var parentId int64
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
func redpActivityShare(ctx *mux.Ctx, actId, parentid int64) {
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
		// 分享者Id 如何获取?
		var parentId int64
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
func redpRecDraw(ctx *mux.Ctx, recId int64, mobile string) {
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

	// 页面类型
	// 1--话费
	// 2--流量
	// 3--无
	redpType := rand.Intn(4-1) + 1
	fmt.Println(redpType)
	redpAmount := 100
	ctx.WriteKeyvals(
		"redpType", redpType,
		"redpAmount", redpAmount,
	)

	// 领取红包后加钱或流量
	// draw_number +1 <draw max
	switch redpType {

	case 2:
		// TODO:加钱

		err = redp.UpdateDrawCount(ctx, rec.ActId, int64(1))
		if err != nil {
			ctx.Error(err)
			return
		}

		fmt.Println(2)

	case 3:
		err := dtplan.Adjust(ctx, uid, strconv.Itoa(redpAmount), "wechat_test", "100", "test")
		if err != nil {
			ctx.Error(err)
			return
		}

		fmt.Println(3)
		// draw_count加1
		err = redp.UpdateDrawCount(ctx, rec.ActId, int64(1))
		if err != nil {
			ctx.Error(err)
			return
		}
	default:
		fmt.Println(1)
	}
}

// 分享页面
// @handler GET /redp/share
func redpShare(ctx *mux.Ctx, recId int64) {
	// load rec
	rec, err := redp.RecById(ctx, recId)
	if err != nil {
		ctx.Error(err)
		return
	}

	// load act
	act, err := redp.ActById(ctx, rec.ActId)
	if err != nil {
		ctx.Error(err)
		return
	}

	userTemplate(ctx, "redp/share", map[string]interface{}{
		"qrcodeImgUrl": act.QrcodeImgUrl,
		"headTitle":    "恭喜您获得30元话费",
	})
}

// 活动规则页面
// @handler GET /redp/rule
func redpRule(ctx *mux.Ctx) {
	userTemplate(ctx, "redp/redp_rules", nil)
}

type ContactList struct {
	UserId string
	Mobile string
}
type FriendRecordList struct {
	RedpTime   string
	HuaFei     int
	NickName   string
	HeadImgUrl string
}

// 点击查看好友
// @handler GET /redp/friendlist/:recId
func redpFriendList(ctx *mux.Ctx, recId int64) {
	uid := token.CurrentUser(ctx)
	user, err := account.UserByID(ctx, uid)
	if err != nil {
		ctx.Error(err)
		return
	}
	userMobile := user.Mobile
	record, err := redp.RecById(ctx, recId)
	if err != nil {
		ctx.Error(err)
		return
	}
	actId := record.ActId

	userList, err := redp.RedpDrawUserList(ctx, actId)
	if err != nil {
		ctx.Error(err)
		return
	}

	contactMobile := []ContactList{}

	for _, v := range userList {
		var c ContactList
		c.UserId = v
		u, err := account.UserByID(ctx, v)
		if err != nil {
			ctx.Error(err)
			return
		}
		c.Mobile = u.Mobile
		contactMobile = append(contactMobile, c)

	}

	contactList, err := redp.GetContactMobile(ctx, userMobile)
	if err != nil {
		ctx.Error(err)
		return
	}

	friendList := []ContactList{}
	for _, v := range contactMobile {
		for _, w := range contactList {
			if v.Mobile == w {
				friendList = append(friendList, v)
			}
		}
	}
	friendRecord := []FriendRecordList{}
	for _, v := range friendList {
		var fre FriendRecordList
		userInt, err := strconv.ParseInt(v.UserId, 10, 64)
		if err != nil {
			ctx.Error(err)
			return
		}
		friend, err := redp.FriendListOfRecord(ctx, userInt, actId)
		if err != nil {
			ctx.Error(err)
			return
		}
		userWechat, err := account.WechatInfo(ctx, v.UserId)
		if err != nil {
			ctx.Error(err)
			return
		}
		fre.NickName = userWechat.NickName
		fre.HeadImgUrl = userWechat.HeadImgurl
		for _, t := range friend {
			fre.RedpTime = t.RedpTime
			fre.HuaFei = t.HuaFei
		}
		friendRecord = append(friendRecord, fre)
	}

	userTemplate(ctx, "redp/redp_list", map[string]interface{}{
		"friendRecord": friendRecord,
	})
}
