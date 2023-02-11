package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ThreeCatsLoveFish/medalhelper/service"
	"github.com/ThreeCatsLoveFish/medalhelper/service/push"
	"github.com/ThreeCatsLoveFish/medalhelper/util"

	"github.com/robfig/cron"
)

var login = flag.Bool("login", false, "登录以获取 accesskey（先执行这个）")
var configPath = flag.String("config", "./users.yaml", "指定配置文件路径")
var startNow = flag.Bool("start", false, "无视定时任务立刻运行一次")

func logo() {
	fmt.Print(`
     __       __                __          __      __    __          __                            
    |  \     /  \              |  \        |  \    |  \  |  \        |  \                           
    | ▓▓\   /  ▓▓ ______   ____| ▓▓ ______ | ▓▓    | ▓▓  | ▓▓ ______ | ▓▓ ______   ______   ______  
    | ▓▓▓\ /  ▓▓▓/      \ /      ▓▓|      \| ▓▓    | ▓▓__| ▓▓/      \| ▓▓/      \ /      \ /      \ 
    | ▓▓▓▓\  ▓▓▓▓  ▓▓▓▓▓▓\  ▓▓▓▓▓▓▓ \▓▓▓▓▓▓\ ▓▓    | ▓▓    ▓▓  ▓▓▓▓▓▓\ ▓▓  ▓▓▓▓▓▓\  ▓▓▓▓▓▓\  ▓▓▓▓▓▓\
    | ▓▓\▓▓ ▓▓ ▓▓ ▓▓    ▓▓ ▓▓  | ▓▓/      ▓▓ ▓▓    | ▓▓▓▓▓▓▓▓ ▓▓    ▓▓ ▓▓ ▓▓  | ▓▓ ▓▓    ▓▓ ▓▓   \▓▓
    | ▓▓ \▓▓▓| ▓▓ ▓▓▓▓▓▓▓▓ ▓▓__| ▓▓  ▓▓▓▓▓▓▓ ▓▓    | ▓▓  | ▓▓ ▓▓▓▓▓▓▓▓ ▓▓ ▓▓__/ ▓▓ ▓▓▓▓▓▓▓▓ ▓▓      
    | ▓▓  \▓ | ▓▓\▓▓     \\▓▓    ▓▓\▓▓    ▓▓ ▓▓    | ▓▓  | ▓▓\▓▓     \ ▓▓ ▓▓    ▓▓\▓▓     \ ▓▓      
     \▓▓      \▓▓ \▓▓▓▓▓▓▓ \▓▓▓▓▓▓▓ \▓▓▓▓▓▓▓\▓▓     \▓▓   \▓▓ \▓▓▓▓▓▓▓\▓▓ ▓▓▓▓▓▓▓  \▓▓▓▓▓▓▓\▓▓      
                                                                        | ▓▓                        
                                                                        | ▓▓                        
                                                                         \▓▓                        

`)
}

func Usage() {
	fmt.Println("Usage of: ", os.Args[0])
	flag.PrintDefaults()
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
		allowId := make([]int, 0)
		if userInfo.AllowedUid != "" {
			allowIdStr := strings.Split(userInfo.AllowedUid, ",")
			for _, str := range allowIdStr {
				id, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					continue
				}
				allowId = append(allowId, int(id))
			}
		}
		users = append(users, service.NewUser(
			userInfo.AccessKey,
			userInfo.PushName,
			allowId,
			banId,
		))
	}
	return users
}

func exec() {
	users := initUsers()
	wg := sync.WaitGroup{}
	for _, user := range users {
		wg.Add(1)
		go func(user service.User, wg *sync.WaitGroup) {
			if status := user.Init(); status {
				user.Start(wg)
			}
		}(user, &wg)
	}
	wg.Wait()
	util.Info(" 今日任务已完成")
}

func main() {
	flag.Parse()

	// Init config file
	util.LoadConfig(*configPath)
	push.InitPush()

	// Tool for login
	if *login {
		util.LoginBili()
		return
	}

	logo()
	// Start right now
	if *startNow {
		exec()
	}

	// Start main block
	if len(util.GlobalConfig.Cron) == 0 {
		util.Info(" 外部调用,开启任务")
		exec()
	} else {
		c := cron.New()
		err := c.AddFunc(util.GlobalConfig.Cron, exec)
		if err != nil {
			panic(err)
		}
		entry := c.Entries()
		timeNext := entry[0].Schedule.Next(time.Now()).Format(time.RFC3339)
		util.Info(" 使用内置定时器,开启定时任务,下次执行时间为%s", timeNext)
		c.Run()
	}
}
