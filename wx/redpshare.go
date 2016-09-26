package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"ndp/lib/mux"
)

const (
	wechatJsapiTicket = `https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi`
)

type WxConfig struct {
	AppId     string // 必填，公众号的唯一标识
	Timestamp string // 必填，生成签名的时间戳
	NonceStr  string // 必填，生成签名的随机串
	Signature string // 必填，签名
}

// @handler GET /redp/share
func redpshare(ctx *mux.Ctx) {
	//TODO:appId,accessToken
	appId := getAppID(ctx)
	accessToken, err := ctx.FormString("access_token")
	if err != nil {
		ctx.Error(err)
		return
	}

	jsapiTicketUrl := fmt.Sprintf(wechatJsapiTicket, accessToken)
	//TODO:get jsticket by url
	jsapiTicket := jsapiTicketUrl

	noncestr :=strrand()
	//TODO:get url
	url := "url"

	timeNowUnix := time.Now().Unix()
	timestamp := strconv.FormatInt(timeNowUnix, 10)
	string1 := fmt.Sprint("jsapi_ticket=" + jsapiTicket + "&noncestr=" + noncestr + "&timestamp=" + timestamp + "&url=" + url)
	signature := toSha1(string1)

	config := &WxConfig{AppId: appId, Timestamp: timestamp, NonceStr: noncestr, Signature: signature}

	userTemplate(ctx, "redp/redp_index", map[string]interface{}{
		"config":config ,
	})

}

func toSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串生成
func strrand() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := r.Intn(4)
	switch x {
	case 0:
		return string(krand(16, KC_RAND_KIND_NUM))
	case 1:
		return string(krand(16, KC_RAND_KIND_LOWER))
	case 2:
		return string(krand(16, KC_RAND_KIND_UPPER))
	case 3:
		return string(krand(16, KC_RAND_KIND_ALL))
	default:
		return "default"
	}
}

// 随机字符串
func krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
