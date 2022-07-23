package service

import (
	"time"

	"github.com/ThreeCatsLoveFish/medalhelper/dto"
	"github.com/ThreeCatsLoveFish/medalhelper/manager"
	"github.com/ThreeCatsLoveFish/medalhelper/util"
)

// Like implement IExec, sync like 3 times
type Like struct {
	SyncAction
}

func (Like) Do(user User, medal dto.MedalInfo) bool {
	if util.GlobalConfig.CD.Like == 0 {
		return true
	}
	times := 1
	ticker := time.NewTicker(time.Duration(util.GlobalConfig.CD.Like) * time.Second)
	for i := 0; i < times; i++ {
		if ok := manager.LikeInteract(user.accessKey, medal.RoomInfo.RoomID); !ok {
			return false
		}
		<-ticker.C
	}
	return true
}

func (Like) Finish(user User, medal []dto.MedalInfo) {
	if util.GlobalConfig.CD.Like == 0 {
		user.info("跳过点赞")
		return
	}
	if len(medal) == 0 {
		user.info("点赞完成")
	} else {
		user.info("点赞未完成,剩余(%d/%d)", len(medal), len(user.medalsLow))
	}
}

// Like implement IExec, async like 3 times
type ALike struct {
	AsyncAction
}

func (ALike) Do(user User, medal dto.MedalInfo) bool {
	if util.GlobalConfig.CD.Like == 0 {
		return true
	}
	times := 1
	for i := 0; i < times; i++ {
		if ok := manager.LikeInteract(user.accessKey, medal.RoomInfo.RoomID); !ok {
			return false
		}
	}
	return true
}

func (ALike) Finish(user User, medal []dto.MedalInfo) {
	if util.GlobalConfig.CD.Like == 0 {
		user.info("跳过点赞")
		return
	}
	if len(medal) == 0 {
		user.info("点赞完成")
	} else {
		user.info("点赞未完成,剩余(%d/%d)", len(medal), len(user.medalsLow))
	}
}

// Share implement IExec, sync share 5 times
type Share struct {
	SyncAction
}

func (Share) Do(user User, medal dto.MedalInfo) bool {
	if util.GlobalConfig.CD.Share == 0 {
		return true
	}
	times := 1
	ticker := time.NewTicker(time.Duration(util.GlobalConfig.CD.Share) * time.Second)
	for i := 0; i < times; i++ {
		if ok := manager.ShareRoom(user.accessKey, medal.RoomInfo.RoomID); !ok {
			return false
		}
		<-ticker.C
	}
	return true
}

func (Share) Finish(user User, medal []dto.MedalInfo) {
	if util.GlobalConfig.CD.Share == 0 {
		user.info("跳过分享")
		return
	}
	if len(medal) == 0 {
		user.info("分享完成")
	} else {
		user.info("分享未完成,剩余(%d/%d)", len(medal), len(user.medalsLow))
	}
}

// Share implement IExec, async share 5 times
type AShare struct {
	AsyncAction
}

func (AShare) Do(user User, medal dto.MedalInfo) bool {
	if util.GlobalConfig.CD.Share == 0 {
		return true
	}
	times := 1
	ticker := time.NewTicker(time.Duration(util.GlobalConfig.CD.Share) * time.Second)
	for i := 0; i < times; i++ {
		if ok := manager.ShareRoom(user.accessKey, medal.RoomInfo.RoomID); !ok {
			return false
		}
		<-ticker.C
	}
	return true
}

func (AShare) Finish(user User, medal []dto.MedalInfo) {
	if util.GlobalConfig.CD.Share == 0 {
		user.info("跳过分享")
		return
	}
	if len(medal) == 0 {
		user.info("分享完成")
	} else {
		user.info("分享未完成,剩余(%d/%d)", len(medal), len(user.medalsLow))
	}
}

// Danmaku implement IExec, default sync, include sending daily danmu
type Danmaku struct {
	SyncAction
}

func (Danmaku) Do(user User, medal dto.MedalInfo) bool {
	if util.GlobalConfig.CD.Danmu == 0 {
		return true
	}
	if ok := manager.WearMedal(user.accessKey, medal.Medal.MedalID); !ok {
		return false
	}
	if ok := manager.SendDanmaku(user.accessKey, medal.RoomInfo.RoomID); !ok {
		return false
	}
	time.Sleep(time.Duration(util.GlobalConfig.CD.Danmu) * time.Second)
	user.info("%s 房间弹幕打卡完成", medal.AnchorInfo.NickName)
	return true
}

func (Danmaku) Finish(user User, medal []dto.MedalInfo) {
	if util.GlobalConfig.CD.Danmu == 0 {
		user.info("跳过弹幕打卡")
		return
	}
	if len(medal) == 0 {
		user.info("弹幕打卡完成")
	} else {
		user.info("弹幕打卡未完成,剩余(%d/%d)", len(medal), len(user.medalsLow))
	}
	if user.wearMedal == dto.DefaultMedal {
		manager.TakeoffMedal(user.accessKey)
		user.info("脱下勋章恢复原样")
	} else {
		manager.WearMedal(user.accessKey, user.wearMedal.Medal.MedalID)
		user.info("重新佩戴勋章 %s", user.wearMedal.Medal.MedalName)
	}
}

// WatchLive implement IExec, default async, include sending heartbeat
type WatchLive struct {
	AsyncAction
}

func (WatchLive) Do(user User, medal dto.MedalInfo) bool {
	times := 80
	for i := 0; i < times; i++ {
		if ok := manager.Heartbeat(
			user.accessKey,
			user.uuid,
			medal.RoomInfo.RoomID,
			medal.Medal.TargetID,
		); !ok {
			return false
		}
		user.info("%s 房间心跳包已发送(%d/%d)", medal.AnchorInfo.NickName, i+1, times)
		time.Sleep(1 * time.Minute)
	}
	return true
}

func (WatchLive) Finish(user User, medal []dto.MedalInfo) {
	if len(medal) == 0 {
		user.info("每日30分钟完成")
	} else {
		user.info("每日30分钟未完成,剩余(%d/%d)", len(medal), len(user.medalsLow))
	}
}
