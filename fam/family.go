package main

import (
	"fmt"
	"time"

	"ndp/core/account"
	"ndp/lib/e164"
	"ndp/lib/mux"
	"ndp/n3d/app/proto"
	"ndp/n3d/family"
)

// 亲情号设置
// @handler POST /app/family/set
func familySet(ctx *mux.Ctx, uid uint64, pid string) {
	// 版本号为时间戳
	version := time.Now().Unix()

	req := ctx.Req.(*proto.Request)

	user, err := account.UserByID(ctx, fmt.Sprint(uid))
	if err != nil {
		ctx.Error(err)
		return
	}

	// 检查号码格式，格式不对不予绑定
	for newNum, _ := range req.Param {
		if _, err := e164.Search(user.Ccode, newNum, true); err != nil {
			ctx.Error(err)
			return
		}
	}
	// 校验参数，分类数据
	addNum := make(map[string]string)
	setNum := make(map[string]string)
	for key, value := range req.Param {
		if value == "" {
			addNum[key] = value.(string)
		} else {
			setNum[key] = value.(string)
		}
	}

	// 校验用户的绑定上限
	bigNumCount, err := family.GetBindNumOfUser(ctx, uid, pid)
	if err != nil {
		ctx.Error(err)
		return
	}
	// 测试阶段为2，正式为10
	if bigNumCount+len(addNum) > 2 {
		ctx.WriteStatus(500, "设定的亲情号已达上限")
		return
	}

	// 若oldNum为零长度号码，服务器新增newNum的绑定
	for newNum, _ := range addNum {
		// 获得大号ID
		bigNumId, err := family.GetBigNumID(ctx, uid, user.Mobile)
		if err != nil {
			ctx.Error(err)
			return
		}
		if bigNumId == 0 {
			ctx.WriteStatus(500, "设定的亲情号已达上限")
			return
		}
		// 将大号和被叫号码绑定
		err = family.AlloctionNum(ctx, uid, uint64(version), pid, user.Mobile, newNum, bigNumId)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	// 若oldNum不为空值
	for newNum, oldNum := range setNum {
		err = family.UpdateBigNumId(ctx, uid, uint64(version), pid, user.Mobile, newNum, oldNum)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
	ctx.WriteStatus(200, "设定成功")
}

// 亲情号查询
// max为最大可设定值,测试时为2,正式为10
// callee为已设定的被叫号码列表
// @handler POST /app/family/num
func familyNum(ctx *mux.Ctx, uid uint64, pid string) {
	callee, max, err := family.GetFamilyNum(ctx, uid, pid)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.WriteKeyvals("max", max, "callee", callee)
}

// 亲情号取消
// @handler POST /app/family/delete
func familyDelete(ctx *mux.Ctx, uid uint64, pid string) {
	version := time.Now().Unix()
	req := ctx.Req.(*proto.Request)
	var nums []string
	for _, inums := range req.Param {
		for _, num := range inums.([]interface{}) {
			nums = append(nums, num.(string))
		}
	}
	for _, callee := range nums {
		err := family.DeleteFamilyNum(ctx, uid, uint64(version), pid, callee)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
	ctx.WriteStatus(200, "取消成功")
}
