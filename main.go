package main

import (
	"MedalHelper/service"
	"MedalHelper/service/push"
	"MedalHelper/util"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron"
)

func usage() {
	fmt.Print(`Usage: main.go [command]

command:
    login   login bili account and get access key
`)
}

func initUsers() []service.User {
	users := make([]service.User, 0, 1)
	for _, userInfo := range util.GlobalConfig.UserList {
		if len(userInfo.AccessKey) == 0 {
			continue
		}
		banId := make([]int, 0)
		if userInfo.BannedUid != "" {
			banIdStr := strings.Split(userInfo.BannedUid, ",")
			for _, str := range banIdStr {
				id, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					continue
				}
				banId = append(banId, int(id))
			}
		}
		users = append(users, service.NewUser(userInfo.AccessKey, userInfo.PushName, banId))
	}
	return users
}

func exec() {
	users := initUsers()
	wg := sync.WaitGroup{}
	for _, user := range users {
		if status := user.Init(); status {
			wg.Add(1)
			go user.Start(&wg)
		}
	}
	wg.Wait()
	util.Info(" 今日任务已完成")
}

func main() {
	// Tool for login
	args := os.Args
	if len(args) > 1 {
		if args[1] == "login" {
			util.LoginBili()
		} else {
			usage()
		}
		return
	}

	// Init config file
	util.InitConfig()
	push.InitPush()
	// Start main block
	if len(util.GlobalConfig.Cron) == 0 {
		util.Info(" 外部调用,开启任务")
		exec()
	} else {
		// Block the process
		c := cron.New()
		c.AddFunc(util.GlobalConfig.Cron, exec)
		entry := c.Entries()
		timeNext := entry[0].Schedule.Next(time.Now()).Format(time.RFC3339)
		util.Info(" 使用内置定时器,开启定时任务,下次执行时间为%s", timeNext)
		c.Run()
	}
}
