package family

import (
	"golang.org/x/net/context"
)

//go:generate ndp_tools_gen -pkg .

// 亲情号设置
// 给出用户已经绑定的资源数目
// @trace
func getBindNumOfUser(ctx context.Context, uid uint64, pid string) (int, error) {
	return dc.getBindNumOfUser(ctx, uid, pid)
}

// 若oldNum为零长度号码，给用户分配大号
// @trace
func getBigNumID(ctx context.Context, uid uint64, caller string) (int, error) {
	return dc.getCurrentBigNum(ctx, caller)
}

// 若oldNum为零长度号码，服务器新增newNum的绑定
// @trace
func alloctionNum(ctx context.Context, uid, version uint64, pid, caller, callee string, bindBigNumId int) error {
	var list []int

	// 得到该用户是否以前绑定过此被叫号码
	list, err := dc.getNumId(ctx, pid, uid, callee)
	if err != nil {
		return err
	}

	// 如果已经绑定过更新状态
	if len(list) == 1 {
		return dc.updateBindBigNumId(ctx, uid, version, pid, callee, bindBigNumId)
	}

	// 没被绑定过：建立新的绑定关系:被叫号码-无流量大号
	return dc.bindData(ctx, uid, version, pid, caller, callee, bindBigNumId)
}

// 若oldNum不为空值,修改数据
// @trace
func updateBigNumId(ctx context.Context, uid, version uint64, pid, caller, callee string, oldNum string) error {
	var list []string
	// 获得大号对应的大号ID
	bigNumId, err := dc.bigNumOfId(ctx, oldNum)
	if err != nil {
		return err
	}

	// 得到该用户是否以前绑定过此被叫号码
	calleeList, err := dc.getNumId(ctx, pid, uid, callee)
	if err != nil {
		return err
	}

	// 服务器检查大号是否被占用
	list, err = dc.checkBigNumId(ctx, pid, uid, bigNumId)
	if err != nil {
		return err
	}

	// 若被占用,修改该绑定为新号码
	if len(list) > 0 {
		// 如果以前绑定过该被叫号码，先删除占用大号的号码，再更新
		if len(calleeList) == 1 {
			err := dc.deleteByBigNum(ctx, uid, version, pid, bigNumId)
			if err != nil {
				return err
			}

			return dc.updateBindBigNumId(ctx, uid, version, pid, callee, bigNumId)
		}

		// 若以前没有此被叫号码的记录，直接修改
		return dc.updateBigNumId(ctx, uid, version, pid, callee, bigNumId)

	} else {

		// 如果以前绑定过
		if len(calleeList) == 1 {
			return dc.updateBindBigNumId(ctx, uid, version, pid, callee, bigNumId)
		}

		return dc.bindData(ctx, uid, version, pid, caller, callee, bigNumId)
	}
}

// 亲情号查询
// @trace
func getFamilyNum(ctx context.Context, uid uint64, pid string) ([]string, string, error) {
	// 获取亲情号码列表
	callee, err := dc.GetCallee(ctx, uid, pid)
	if err != nil {
		return nil, "", err
	}

	// 获取大号总数
	max, err := dc.getMaxBigNum(ctx)
	if err != nil {
		return nil, "", err
	}
	return callee, max, nil
}

// 亲情号取消
// @trace
func deleteFamilyNum(ctx context.Context, uid, version uint64, pid, callee string) error {
	return dc.deleteCallee(ctx, uid, version, pid, callee)
}
