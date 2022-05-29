package service

import (
	"MedalHelper/dto"
	"MedalHelper/manager"
	"MedalHelper/service/push"
	"MedalHelper/util"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/TwiN/go-color"
	"github.com/google/uuid"
	"github.com/sethvargo/go-retry"
	"github.com/tidwall/gjson"
)

type User struct {
	// 用户ID
	Uid int
	// 用户名称
	Name string
	// 是否登录
	isLogin bool
	// UUID
	uuid []string

	// 登录凭证
	accessKey string
	// 被禁止的房间ID
	bannedUIDs []int
	// 推送服务
	pushName string

	// 用户所有勋章
	medals []dto.MedalInfo
	// 用户等级小于20的勋章
	medalsLow []dto.MedalInfo
	// 今日亲密度没满的勋章
	remainMedals []dto.MedalInfo

	// 日志信息
	message string
}

func NewUser(accessKey, pushName string, uids []int) User {
	return User{
		accessKey:  accessKey,
		bannedUIDs: uids,
		pushName:   pushName,
		uuid:       []string{uuid.NewString(), uuid.NewString()},
		message:    "",
	}
}

func (user User) info(format string, v ...interface{}) {
	format = color.Green + "[INFO] " + color.Reset + format
	format = color.Reset + color.Blue + user.Name + color.Reset + " " + format
	util.PrintColor(format, v...)
}

func (user *User) loginVerify() bool {
	resp, err := manager.LoginVerify(user.accessKey)
	if err != nil || resp.Data.Mid == 0 {
		user.isLogin = false
		return false
	}
	user.Uid = resp.Data.Mid
	user.Name = resp.Data.Name
	user.isLogin = true
	user.info("登录成功")
	return true
}

func (user *User) signIn() error {
	signInfo, err := manager.SignIn(user.accessKey)
	if err != nil {
		return nil
	}
	resp := gjson.Parse(signInfo)
	if resp.Get("code").Int() == 0 {
		signed := resp.Get("data.hadSignDays").String()
		all := resp.Get("data.allDays").String()
		user.info("签到成功, 本月签到次数: %s/%s", signed, all)
	} else {
		user.info("%s", resp.Get("message").String())
	}

	userInfo, err := manager.GetUserInfo(user.accessKey)
	if err != nil {
		return nil
	}
	level := userInfo.Data.Exp.UserLevel
	unext := userInfo.Data.Exp.Unext
	user.info("当前用户UL等级: %d, 还差 %d 经验升级", level, unext)
	return nil
}

func (user *User) setMedals() {
	// Clean medals storage
	user.medals = make([]dto.MedalInfo, 0, 10)
	user.medalsLow = make([]dto.MedalInfo, 0, 10)
	user.remainMedals = make([]dto.MedalInfo, 0, 10)
	// Fetch and update medals
	for _, medal := range manager.GetFansMedalAndRoomID(user.accessKey) {
		if util.IntContain(user.bannedUIDs, medal.Medal.TargetID) != -1 {
			continue
		}
		if medal.RoomInfo.RoomID == 0 {
			continue
		}
		user.medals = append(user.medals, medal)
		if medal.Medal.Level <= 20 {
			user.medalsLow = append(user.medalsLow, medal)
			if medal.Medal.TodayFeed < 1300 {
				user.remainMedals = append(user.remainMedals, medal)
			}
		}
	}
}

func (user *User) checkMedals() bool {
	user.setMedals()
	medalList1 := make([]string, 0, len(user.medalsLow))
	medalList2 := make([]string, 0)
	medalList3 := make([]string, 0)
	medalList4 := make([]string, 0)
	for _, medal := range user.medalsLow {
		if medal.Medal.TodayFeed >= 1300 {
			medalList1 = append(medalList1, medal.AnchorInfo.NickName)
		} else if medal.Medal.TodayFeed >= 1200 {
			medalList2 = append(medalList2, medal.AnchorInfo.NickName)
		} else if medal.Medal.TodayFeed >= 1100 {
			medalList3 = append(medalList3, medal.AnchorInfo.NickName)
		} else if medal.Medal.TodayFeed >= 1000 {
			medalList4 = append(medalList4, medal.AnchorInfo.NickName)
		}
	}
	user.message = fmt.Sprintf(
		"20级以下牌子共 %d 个\n【1300及以上】 %v等 %d个\n【1200至1300】 %v等 %d个\n【1100至1200】 %v等 %d个\n【1100以下】 %v等 %d个\n",
		len(user.medalsLow),
		medalList1, len(medalList1),
		medalList2, len(medalList2),
		medalList3, len(medalList3),
		medalList4, len(medalList4),
	)
	user.info(user.message)
	return len(medalList1) == len(user.medalsLow)
}

func (user *User) report() {
	if len(user.pushName) != 0 {
		pushEnd := push.NewPush(user.pushName)
		pushEnd.Submit(push.Data{
			Title:   "# 今日亲密度获取情况如下",
			Content: fmt.Sprintf("用户%s，%s", user.Name, user.message),
		})
	}
}

func (user *User) Init() bool {
	if user.loginVerify() {
		user.signIn()
		user.setMedals()
		return true
	} else {
		util.Error("用户登录失败, accessKey: %s", user.accessKey)
		return false
	}
}

func (user *User) RunOnce() bool {
	switch util.GlobalConfig.CD.Async {
	case 0: // Sync
		task := NewTask(*user, []IAction{
			&Like{},
			&Share{},
			&Danmaku{},
			&WatchLive{},
		})
		task.Start()
	case 1: // Async
		task := NewTask(*user, []IAction{
			&ALike{},
			&AShare{},
			&Danmaku{},
			&WatchLive{},
		})
		task.Start()
	}
	return user.checkMedals()
}

func (user *User) Start(wg *sync.WaitGroup) {
	if user.isLogin {
		backOff := retry.NewConstant(5 * time.Second)
		backOff = retry.WithMaxRetries(3, backOff)
		retry.Do(context.Background(), backOff, func(ctx context.Context) error {
			if ok := user.RunOnce(); !ok {
				return retry.RetryableError(errors.New("task not complete"))
			}
			return nil
		})
		user.report()
	} else {
		util.Error("用户未登录, accessKey: %s", user.accessKey)
	}
	wg.Done()
}
