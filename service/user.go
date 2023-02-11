package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ThreeCatsLoveFish/medalhelper/dto"
	"github.com/ThreeCatsLoveFish/medalhelper/manager"
	"github.com/ThreeCatsLoveFish/medalhelper/service/push"
	"github.com/ThreeCatsLoveFish/medalhelper/util"

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
	// 白名单的房间ID
	allowedUIDs []int
	// 被禁止的房间ID
	bannedUIDs []int
	// 推送服务
	pushName string

	// 用户佩戴的勋章
	wearMedal dto.MedalInfo
	// 用户等级小于20的勋章
	medalsLow []dto.MedalInfo
	// 今日亲密度没满的勋章
	remainMedals []dto.MedalInfo

	// 日志信息
	message string
}

func NewUser(accessKey, pushName string, allowUIDs, banUIDs []int) User {
	return User{
		accessKey:   accessKey,
		allowedUIDs: allowUIDs,
		bannedUIDs:  banUIDs,
		pushName:    pushName,
		wearMedal:   dto.DefaultMedal,
		uuid:        []string{uuid.NewString(), uuid.NewString()},
		message:     "",
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

func (user *User) signIn() {
	signInfo, err := manager.SignIn(user.accessKey)
	if err != nil {
		return
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
		return
	}
	level := userInfo.Data.Exp.UserLevel
	unext := userInfo.Data.Exp.Unext
	user.info("当前用户UL等级: %d, 还差 %d 经验升级", level, unext)
}

func (user *User) setMedals() {
	// Clean medals storage
	user.medalsLow = make([]dto.MedalInfo, 0, 10)
	user.remainMedals = make([]dto.MedalInfo, 0, 10)
	// Fetch and update medals
	medals, wearMedal := manager.GetMedal(user.accessKey)
	if wearMedal {
		user.wearMedal = medals[0]
	}
	// Whitelist
	if len(user.allowedUIDs) > 0 {
		for _, medal := range medals {
			if util.IntContain(user.allowedUIDs, medal.Medal.TargetID) != -1 {
				user.medalsLow = append(user.medalsLow, medal)
				if medal.Medal.TodayFeed < 1300 {
					user.remainMedals = append(user.remainMedals, medal)
				}
			}
		}
		return
	}
	// Default blacklist
	for _, medal := range medals {
		if util.IntContain(user.bannedUIDs, medal.Medal.TargetID) != -1 {
			continue
		}
		if medal.RoomInfo.RoomID == 0 {
			continue
		}
		if medal.Medal.Level <= 20 {
			user.medalsLow = append(user.medalsLow, medal)
			if medal.Medal.TodayFeed < 1500 {
				user.remainMedals = append(user.remainMedals, medal)
			}
		}
	}
}

func (user *User) checkMedals() bool {
	user.setMedals()
	fullMedalList := make([]string, 0, len(user.medalsLow))
	failMedalList := make([]string, 0)
	for _, medal := range user.medalsLow {
		if medal.Medal.TodayFeed == 1500 {
			fullMedalList = append(fullMedalList, medal.AnchorInfo.NickName)
		} else {
			failMedalList = append(failMedalList, medal.AnchorInfo.NickName)
		}
	}
	user.message = fmt.Sprintf(
		"20级以下牌子共 %d 个\n【1500】%d个\n【1500以下】 %v等 %d个\n",
		len(user.medalsLow), len(fullMedalList),
		failMedalList, len(failMedalList),
	)
	user.info(user.message)
	return len(fullMedalList) == len(user.medalsLow)
}

// Send daily report notification
func (user *User) report() {
	if len(user.pushName) != 0 {
		pushEnd := push.NewPush(user.pushName)
		_ = pushEnd.Submit(push.Data{
			Title:   "# 今日亲密度获取情况如下",
			Content: fmt.Sprintf("用户%s，%s", user.Name, user.message),
		})
	}
}

// Send expire notification
func (user *User) expire() {
	if len(user.pushName) != 0 {
		pushEnd := push.NewPush(user.pushName)
		_ = pushEnd.Submit(push.Data{
			Title:   "# AccessKey 过期",
			Content: fmt.Sprintf("用户未登录, accessKey: %s", user.accessKey),
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
		user.expire()
		return false
	}
}

func (user *User) RunOnce() bool {
	switch util.GlobalConfig.CD.Async {
	case 0: // Sync
		task := NewTask(*user, []IAction{
			&Like{},
			&Danmaku{},
			&WatchLive{},
		})
		task.Start()
	case 1: // Async
		task := NewTask(*user, []IAction{
			&ALike{},
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
		_ = retry.Do(context.Background(), backOff, func(ctx context.Context) error {
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
