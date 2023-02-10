package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	gjson "github.com/tidwall/gjson"
)

var accessKey string

func getQRcode() (string, string) {
	api := "http://passport.bilibili.com/x/passport-tv-login/qrcode/auth_code"
	data := map[string]string{
		"local_id": "0",
		"ts":       GetTimestamp(),
	}
	Signature(&data)
	dataString := strings.NewReader(Map2string(data))
	client := http.Client{}
	req, _ := http.NewRequest("POST", api, dataString)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	code := gjson.Parse(string(body)).Get("code").Int()
	if code == 0 {
		qrcodeUrl := gjson.Parse(string(body)).Get("data.url").String()
		authCode := gjson.Parse(string(body)).Get("data.auth_code").String()
		return qrcodeUrl, authCode
	} else {
		panic("getQRcode error")
	}
}

func verifyLogin(auth_code string) {
	for {
		time.Sleep(time.Second * 3)
		api := "http://passport.bilibili.com/x/passport-tv-login/qrcode/poll"
		data := map[string]string{
			"auth_code": auth_code,
			"local_id":  "0",
			"ts":        GetTimestamp(),
		}
		Signature(&data)
		dataString := strings.NewReader(Map2string(data))
		client := http.Client{}
		req, _ := http.NewRequest("POST", api, dataString)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		code := gjson.Parse(string(body)).Get("code").Int()
		accessKey = gjson.Parse(string(body)).Get("data.access_token").String()
		if code == 0 {
			fmt.Println("登录成功")
			fmt.Println("access_key:", string(accessKey))
			filename := "login_info.txt"
			err := ioutil.WriteFile(filename, []byte(string(accessKey) + "\n"), 0644)
			if err != nil {
				panic(err)
			}
			fmt.Printf("access_key 已保存在 %s\n", filename)
			break
		} else {
			fmt.Println(string(body))
		}
	}
}

func LoginBili() {
	fmt.Println("请最大化窗口，以确保二维码完整显示，回车继续")
	fmt.Scanf("%s", "")
	loginUrl, authCode := getQRcode()
	qrcode := qrcodeTerminal.New()
	qrcode.Get([]byte(loginUrl)).Print()
	fmt.Println("或将此链接复制到手机B站打开:", loginUrl)
	verifyLogin(authCode)
}
