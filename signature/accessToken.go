package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	getAccessToken()
	//	getJsapiTicket()
	//	getSignature()
}

func getAccessToken() {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wxe668d7df1101ee6a&secret=66e7ee651a90a9653cb19725b3a7dd66"
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response", string(content))
	resp.Body.Close()

}

func getJsapiTicket() {
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=rAGwr3tY0swElHUyGYTP0ZeZ3cAGdN_GVIvuHNouS9ntu9wqD5gqvapEvyBW4LI07cxlPlh_tOKT9q3kDrZCJ4DhnDWrIFW9sAlScrGLxhvCCPDZ67q4e6bg1cAj-SSXGPUcADAZIN&type=jsapi"
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response", string(content))
	resp.Body.Close()
}

func getSignature() {
	string1 := "jsapi_ticket=sM4AOVdWfPE4DxkXGEs8VMropoM9K8moU2DOD_9f51sjB1qfRECOMaf5J7xaYhkS04B0pOT3W1iUx-AcFdOnNg&noncestr=Wm3WZYTPz0wzccnW&timestamp=1414587457&url=http://mp.weixin.qq.com?params=value"
	fmt.Printf("%x", sha1.Sum([]byte(string1)))
}
