```
         __  __          _       _   _   _      _                 
        |  \/  | ___  __| | __ _| | | | | | ___| |_ __   ___ _ __ 
        | |\/| |/ _ \/ _` |/ _` | | | |_| |/ _ \ | '_ \ / _ \ '__|
        | |  | |  __/ (_| | (_| | | |  _  |  __/ | |_) |  __/ |   
        |_|  |_|\___|\__,_|\__,_|_| |_| |_|\___|_| .__/ \___|_|   
                                                 |_|              
```

<div align="center">
  <h1> 最新 B 站粉丝牌助手</h1>
  <p>当前版本：1.0</p>
</div>

### 功能说明

- 每日直播区签到
- 每日观看 30 分钟
- 每日点赞 3 次直播间 （200\*3 亲密度）
- 每日分享 5 次直播间 （100\*5 亲密度）
- 每日自定义弹幕打卡 （100 亲密度）
- 多账号支持
- 本日亲密度已满徽章不重复打卡
- 可选需要的打卡类型
- 多种推送通知
- 同步异步配置

<small>ps: 新版 B 站粉丝牌的亲密度每一个牌子都将单独计算  </small>

### 使用说明

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ThreeCatsLoveFish/MedalHelper)
[![wakatime](https://wakatime.com/badge/github/ThreeCatsLoveFish/MedalHelper.svg)](https://wakatime.com/badge/github/ThreeCatsLoveFish/MedalHelper)

#### 环境需求：Go 1.16

> 克隆本项目 安装依赖

```shell
git clone https://github.com/ThreeCatsLoveFish/MedalHelper.git
cd fansMedalHelper
```

> 获取 B 站账号的 access_key

```shell
go run main.go login
```
扫码登录，会得到 `access_key` 即可

> 填写配置文件 users.yaml

```shell
vim users.yaml
```

```yaml
USERS:
  - access_key: XXXXXX # 注意冒号后的空格 否则会读取失败 英文冒号
    banned_uid: 789,100 # 黑名单UID 同上,填了后将不会打卡，点赞，分享 用英文逗号分隔 不填则不限制
    push_name: PUSH_DEER_SAMPLE # 推送服务，留空表示不需要推送
  - access_key:
    banned_uid:
    push_name:
  # 注意对齐
  # 多用户以上格式添加
  # 井号后为注释 井号前后必须有空格
DANMU:
  [
    "(⌒▽⌒).",
    "（￣▽￣）.",
    "(=・ω・=).",
    "(｀・ω・´).",
    "(〜￣△￣)〜.",
    "(･∀･).",
    "(°∀°)ﾉ.",
    "(￣3￣).",
    "╮(￣▽￣)╭.",
    "_(:3」∠)_.",
    "(^・ω・^ ).",
    "(●￣(ｴ)￣●).",
    "ε=ε=(ノ≧∇≦)ノ.",
    "⁄(⁄ ⁄•⁄ω⁄•⁄ ⁄)⁄.",
    "←◡←.",
  ]
  # 可在此处自定义打卡弹幕
CRON: #3 2 1 * *
# 这里是 cron 表达式, 从左到右参数为秒，分钟，小时，日期，月份
# Second | Minute | Hour | Dom | Month
# 例如每天凌晨01点02分03秒执行一次为 3 2 1 * *
# 如果不填,则不使用内置定时器,填写正确后要保持该进程一直运行
CD:
  async: 1 # 异步执行，默认为1表示异步，0表示同步
  retry: 1 # 任务失败重试时间，单位秒，设置为0不重试
  max_try: 10 # 任务失败最多重试次数，单位次，设置为0不重试
  like: 2 # 点赞间隔时间，单位秒，设置为0不点赞
  share: 5 # 分享间隔时间，单位秒，设置为0不分享
  danmu: 6 # 弹幕间隔时间，单位秒，设置为0不发送弹幕，只支持同步
PUSH:
  - name: "PUSH_PLUS_SAMPLE" # 推送名称，对应上面对应用户的推送，请保证名称唯一
    token: "<YOUR-TOKEN-HERE>"  # 推送服务TOKEN
    type: "push_plus" # 推送服务类型为 PushPlus
    url: "http://www.pushplus.plus/send" # 推送服务URL
  - name: "PUSH_DEER_SAMPLE" 
    token: "<YOUR-TOKEN-HERE>" 
    type: "push_deer" # 推送服务类型为 PushDeer
    url: "http://<pushdeer-url-or-ip>/message/push" 
# 推送服务，每日打卡成功或报错日志推送
# 目前仅支持PushDeer和PushPlus
```

请务必严格填写，否则程序将读取失败，可以在这里 [YAML、YML 在线编辑器(格式化校验)-BeJSON.com](https://www.bejson.com/validators/yaml_editor/) 验证你填的 yaml 是否正确

> 运行主程序

```shell
go run main.go
```

### 友情链接

Python 版本可前往 [新B站粉丝牌助手](https://github.com/XiaoMiku01/fansMedalHelper)
