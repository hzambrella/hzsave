package wxconfig

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// js api ticket

var (
	jsapiTicketMu sync.Mutex
	jsapiTickets  = map[string]*ExpireTicket{}
)

type ExpireTicket struct {
	Ticket   string
	ExpireAt time.Time
}

// 如果获取缓存的，那么只要AppId；如果要重新获取，就要Secret
func GetJsApiTicket(cfg *Config) (string, error) {
	jsapiTicketMu.Lock()
	defer jsapiTicketMu.Unlock()

	// 检查是否存在，并且是否有效
	ticket, ok := jsapiTickets[cfg.AppId]
	if ok {
		if time.Now().Before(ticket.ExpireAt) {
			return ticket.Ticket, nil
		}
	}

	// token
	token, err := GetAccessToken(cfg)
	if err != nil {
		return "", err
	}

	// get jsApiTicket
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	cfgUrl := fmt.Sprintf(url, token)

	// response
	resp, err := http.Get(cfgUrl)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	result := make(Data)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if int(result.Int("errcode")) != 0 {
		err := errors.New(fmt.Sprint(result["errmsg"]))
		return "", err
	}

	// ticket
	expiresIn := time.Duration(result.Int64("expires_in")-200) * time.Second
	if expiresIn <= 0 {
		log.Println("wxmp: [WARN] expires_in should > 0")
		expiresIn = 3600 * time.Second
	}
	ticket = &ExpireTicket{
		ExpireAt: time.Now().Add(expiresIn),
		Ticket:   fmt.Sprint(result["ticket"]),
	}
	jsapiTickets[cfg.AppId] = ticket

	// ok
	return ticket.Ticket, nil
}

// ------------------------------------------------------------

// 签名

const (
	jsapiSignTpl = "jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s"
)

type JsApiSign struct {
	Timestamp int64
	NonceStr  string
	Signature string
}

func NewJsApiSign(cfg *Config, currUrl string) (*JsApiSign, error) {

	// ticket
	ticket, err := GetJsApiTicket(cfg)
	if err != nil {
		return nil, err
	}

	// sign
	timestamp := time.Now().Unix()
	// 没找到库，这里先写死
	//	nonceStr := randx.String(12)
	nonceStr := "Wm3WZYTPz0wzccnW"
	t := sha1.New()
	io.WriteString(t, fmt.Sprintf(jsapiSignTpl, ticket, nonceStr, timestamp, currUrl))
	signature := fmt.Sprintf("%x", t.Sum(nil))

	// return
	return &JsApiSign{
		Timestamp: timestamp,
		NonceStr:  nonceStr,
		Signature: signature,
	}, nil
}
