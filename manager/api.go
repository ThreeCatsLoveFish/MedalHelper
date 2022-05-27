package manager

import (
	"MedalHelper/dto"
	"MedalHelper/util"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func LoginVerify(accessKey string) (dto.BiliAccountResp, error) {
	rawUrl := "http://app.bilibili.com/x/v2/account/mine"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	var resp dto.BiliAccountResp
	body, err := Get(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("LoginVerify error: %v, data: %v", err, data)
		return resp, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliAccountResp error: %v, raw data: %v", err, body)
		return resp, err
	}
	return resp, nil
}

func SignIn(accessKey string) (string, error) {
	rawUrl := "http://api.live.bilibili.com/rc/v1/Sign/doSign"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	body, err := Get(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("SignIn error: %v, data: %v", err, data)
		return "", err
	}
	return string(body), nil
}

func GetUserInfo(accessKey string) (dto.BiliLiveUserInfo, error) {
	rawUrl := "http://api.live.bilibili.com/xlive/app-ucenter/v1/user/get_user_info"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	body, err := Get(rawUrl, util.Map2Params(data))
	var resp dto.BiliLiveUserInfo
	if err != nil {
		util.Error("GetUserInfo error: %v, data: %v", err, data)
		return resp, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliLiveUserInfo error: %v, raw data: %v", err, body)
		return resp, err
	}
	return resp, nil
}

func GetFansMedalAndRoomID(accessKey string) []dto.MedalInfo {
	medals := make([]dto.MedalInfo, 0, 20)
	page := 1
	for {
		rawUrl := "http://api.live.bilibili.com/xlive/app-ucenter/v1/fansMedal/panel"
		data := map[string]string{
			"access_key": accessKey,
			"actionKey":  "appkey",
			"appkey":     util.AppKey,
			"ts":         util.GetTimestamp(),
			"page":       fmt.Sprint(page),
			"page_size":  "100",
		}
		util.Signature(&data)
		body, err := Get(rawUrl, util.Map2Params(data))
		if err != nil {
			util.Error("GetFansMedalAndRoomID error: %v, data: %v", err, data)
			return medals
		}
		var resp dto.BiliMedalResp
		if err = json.Unmarshal(body, &resp); err != nil {
			util.Error("Unmarshal BiliMedalResp error: %v, raw data: %v", err, body)
			return medals
		}
		medals = append(medals, resp.Data.SpecialList...)
		medals = append(medals, resp.Data.List...)
		if len(resp.Data.List) == 0 {
			break
		}
		page++
	}
	return medals
}

func LikeInteract(accessKey string, roomId int) bool {
	rawUrl := "http://api.live.bilibili.com/xlive/web-ucenter/v1/interact/likeInteract"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
		"roomid":     fmt.Sprint(roomId),
	}
	util.Signature(&data)
	body, err := Post(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("LikeInteract error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

func ShareRoom(accessKey string, roomId int) bool {
	rawUrl := "http://api.live.bilibili.com/xlive/app-room/v1/index/TrigerInteract"
	data := map[string]string{
		"access_key":    accessKey,
		"actionKey":     "appkey",
		"appkey":        util.AppKey,
		"interact_type": "3",
		"ts":            util.GetTimestamp(),
		"roomid":        fmt.Sprint(roomId),
	}
	util.Signature(&data)
	body, err := Post(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("ShareRoom error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

func SendDanmaku(accessKey string, roomId int) bool {
	rawUrl := "http://api.live.bilibili.com/xlive/app-room/v1/dM/sendmsg"
	params := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	data := map[string]string{
		"cid":      fmt.Sprint(roomId),
		"msg":      util.GlobalConfig.Danmuku[rand.Intn(len(util.GlobalConfig.Danmuku))],
		"rnd":      util.GetTimestamp(),
		"color":    "16777215",
		"fontsize": "25",
	}
	util.Signature(&params)
	body, err := PostWithParam(rawUrl, util.Map2Params(params), util.Map2Params(data))
	if err != nil {
		util.Error("GetFansMedalAndRoomID error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

func Heartbeat(accessKey string, uuids []string, roomId, upId int) bool {
	rawUrl := "https://live-trace.bilibili.com/xlive/data-interface/v1/heartbeat/mobileHeartBeat"
	data := map[string]string{
		"platform":         "android",
		"uuid":             uuids[0],
		"buvid":            strings.ToUpper(util.RandomString(37)),
		"seq_id":           "1",
		"room_id":          fmt.Sprint(roomId),
		"parent_id":        "6",
		"area_id":          "283",
		"timestamp":        fmt.Sprintf("%d", time.Now().Unix()-60),
		"secret_key":       "axoaadsffcazxksectbbb",
		"watch_time":       "60",
		"up_id":            fmt.Sprint(upId),
		"up_level":         "40",
		"jump_from":        "30000",
		"gu_id":            strings.ToUpper(util.RandomString(43)),
		"play_type":        "0",
		"play_url":         "",
		"s_time":           "0",
		"data_behavior_id": "",
		"data_source_id":   "",
		"up_session":       fmt.Sprintf("l:one:live:record:%d:%d", roomId, time.Now().Unix()-88888),
		"visit_id":         strings.ToUpper(util.RandomString(32)),
		"watch_status":     "%7B%22pk_id%22%3A0%2C%22screen_status%22%3A1%7D",
		"click_id":         uuids[1],
		"session_id":       "",
		"player_type":      "0",
		"client_ts":        util.GetTimestamp(),
	}
	data["client_sign"] = util.ClientSign(data)
	data["access_key"] = accessKey
	data["actionKey"] = "appkey"
	data["appkey"] = util.AppKey
	data["ts"] = util.GetTimestamp()
	util.Signature(&data)
	body, err := Post(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("Heartbeat error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}
