package manager

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ThreeCatsLoveFish/medalhelper/dto"
	"github.com/ThreeCatsLoveFish/medalhelper/util"
)

func LoginVerify(accessKey string) (dto.BiliAccountResp, error) {
	rawUrl := "https://app.bilibili.com/x/v2/account/mine"
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
	rawUrl := "https://api.live.bilibili.com/rc/v1/Sign/doSign"
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
	rawUrl := "https://api.live.bilibili.com/xlive/app-ucenter/v1/user/get_user_info"
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

func GetMedal(accessKey string) ([]dto.MedalInfo, bool) {
	medals := make([]dto.MedalInfo, 0, 20)
	wear := false
	page := 1
	for {
		rawUrl := "https://api.live.bilibili.com/xlive/app-ucenter/v1/fansMedal/panel"
		data := map[string]string{
			"access_key": accessKey,
			"actionKey":  "appkey",
			"appkey":     util.AppKey,
			"ts":         util.GetTimestamp(),
			"page":       fmt.Sprint(page),
			"page_size":  "50",
		}
		util.Signature(&data)
		body, err := Get(rawUrl, util.Map2Params(data))
		if err != nil {
			util.Error("GetFansMedalAndRoomID error: %v, data: %v", err, data)
			return medals, wear
		}
		var resp dto.BiliMedalResp
		if err = json.Unmarshal(body, &resp); err != nil {
			util.Error("Unmarshal BiliMedalResp error: %v, raw data: %v", err, body)
			return medals, wear
		}
		if len(resp.Data.SpecialList) > 0 {
			wear = true
		}
		medals = append(medals, resp.Data.SpecialList...)
		medals = append(medals, resp.Data.List...)
		if len(resp.Data.List) == 0 {
			break
		}
		page++
	}
	return medals, wear
}

func WearMedal(accessKey string, medalId int) bool {
	rawUrl := "https://api.live.bilibili.com/xlive/app-ucenter/v1/fansMedal/wear"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
		"medal_id":   fmt.Sprint(medalId),
		"platform":   "android",
		"type":       "1",
		"version":    "0",
	}
	util.Signature(&data)
	body, err := Post(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("WearMedal error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

func TakeoffMedal(accessKey string) bool {
	rawUrl := "https://api.live.bilibili.com/xlive/app-ucenter/v1/fansMedal/take_off"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
		"platform":   "android",
		"type":       "1",
		"version":    "0",
	}
	util.Signature(&data)
	body, err := Post(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("TakeoffMedal error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

func LikeInteract(accessKey string, roomId int) bool {
	rawUrl := "https://api.live.bilibili.com/xlive/web-ucenter/v1/interact/likeInteract"
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
	rawUrl := "https://api.live.bilibili.com/xlive/app-room/v1/index/TrigerInteract"
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
	rawUrl := "https://api.live.bilibili.com/xlive/app-room/v1/dM/sendmsg"
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
	dataStr := fmt.Sprintf(`{"platform":"%s","uuid":"%s","buvid":"%s","seq_id":"%s","room_id":"%s","parent_id":"%s","area_id":"%s","timestamp":"%s","secret_key":"%s","watch_time":"%s","up_id":"%s","up_level":"%s","jump_from":"%s","gu_id":"%s","play_type":"%s","play_url":"%s","s_time":"%s","data_behavior_id":"%s","data_source_id":"%s","up_session":"%s","visit_id":"%s","watch_status":"%s","click_id":"%s","session_id":"%s","player_type":"%s","client_ts":"%s"}`,
		data["platform"],
		data["uuid"],
		data["buvid"],
		data["seq_id"],
		data["room_id"],
		data["parent_id"],
		data["area_id"],
		data["timestamp"],
		data["secret_key"],
		data["watch_time"],
		data["up_id"],
		data["up_level"],
		data["jump_from"],
		data["gu_id"],
		data["play_type"],
		data["play_url"],
		data["s_time"],
		data["data_behavior_id"],
		data["data_source_id"],
		data["up_session"],
		data["visit_id"],
		data["watch_status"],
		data["click_id"],
		data["session_id"],
		data["player_type"],
		data["client_ts"],
	)
	data["client_sign"] = util.ClientSign(dataStr)
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
