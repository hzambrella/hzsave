package wxconfig

import (
	"fmt"
	"testing"
)

var (
	cfg = &Config{
		AppId:  "wxe668d7df1101ee6a",
		Secret: "66e7ee651a90a9653cb19725b3a7dd66",
	}
)

func TestGetAccessToken(t *testing.T) {
	res, err := GetAccessToken(cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetJsApiTicket(t *testing.T) {
	res, err := GetJsApiTicket(cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestNewJsApiSign(t *testing.T) {
	currUrl := "http://192.168.0.134:9500/wxapp/v1/redp/activity/1"
	res, err := NewJsApiSign(cfg, currUrl)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
