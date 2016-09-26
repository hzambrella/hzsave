package wxconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	accessTokenMu sync.Mutex
	accessTokens  = map[string]*ExpireToken{}
)

type ExpireToken struct {
	Token    string
	ExpireAt time.Time
}

// 如果获取缓存的，那么只要AppId；如果要重新获取，就要Secret
func GetAccessToken(cfg *Config) (string, error) {
	accessTokenMu.Lock()
	defer accessTokenMu.Unlock()

	// 如果仅仅是读取，只要AppId就可以了
	if len(cfg.AppId) == 0 {
		return "", errors.New("ctx no AppId")
	}
	// 检查是否存在，并且是否有效
	token, ok := accessTokens[cfg.AppId]
	if ok {
		if time.Now().Before(token.ExpireAt) {
			return token.Token, nil
		}
	}

	// 如果要重新获取，必须传递Secret
	if len(cfg.Secret) == 0 {
		return "", errors.New("ctx no Secret")
	}

	// get access_token
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

	cfgUrl := fmt.Sprintf(url, cfg.AppId, cfg.Secret)
	fmt.Println("----------1----------")

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
	fmt.Println("----------2----------")
	if int(result.Int("errcode")) != 0 {
		err := errors.New(fmt.Sprint(result["errmsg"]))
		return "", err
	}

	fmt.Println("----------3----------")
	// token
	expiresIn := time.Duration(result.Int64("expires_in")-200) * time.Second
	if expiresIn <= 0 {
		log.Println("weixin: [WARN] expires_in should > 0")
		expiresIn = 3600 * time.Second
	}
	token = &ExpireToken{
		ExpireAt: time.Now().Add(expiresIn),
		Token:    fmt.Sprint(result["access_token"]),
	}
	accessTokens[cfg.AppId] = token

	// ok
	return token.Token, nil

}
