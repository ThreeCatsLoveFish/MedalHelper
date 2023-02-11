# 使用说明

### 使用 Github Actions

参考 [run.yml](https://github.com/ThreeCatsLoveFish/MedalHelper/blob/master/.github/workflows/run.yml)

### Go 直接安装

```bash
go install github.com/ThreeCatsLoveFish/medalhelper@latest
```

### Windows 用户

> Windows10以上用户请直接前往[免配置高速通道](https://github.com/ThreeCatsLoveFish/MedalHelper/releases/tag/v1.4)

### Docker 用户

> 克隆本项目，构建镜像

```shell
git clone https://github.com/ThreeCatsLoveFish/MedalHelper.git
cd MedalHelper
docker build -t medalhelper .
```

> 获取 B 站账号的 access_key

```shell
docker run --rm -ti -v $(pwd)/users.yaml:/config/users.yaml medalhelper -login
```

按提示回车并扫码（或访问 URL），得到 `access_key`。

> 填写配置文件 users.yaml

[参考下文](#配置文件)编辑 `users.yaml`。注意 `CRON` 中填写的时间，时区需要使用环境变量 `TZ` 指定。

> 运行

```shell
docker run -d \
    -e TZ=Asia/Shanghai \
    -v $(pwd)/users.yaml:/config/users.yaml \
    --restart unless-stopped \
    --name medalhelper \
    medalhelper
```

> 无视定时任务立刻运行一次

```shell
docker run -d \
    -e TZ=Asia/Shanghai \
    -v $(pwd)/users.yaml:/config/users.yaml \
    --restart unless-stopped \
    --name medalhelper \
    medalhelper -start
```

> 查看日志

```shell
docker logs medalhelper
```

### 自行编译

环境需求：Go 1.16

> 克隆本项目 安装依赖

```shell
git clone https://github.com/ThreeCatsLoveFish/MedalHelper.git
cd MedalHelper
```

> 获取 B 站账号的 access_key

```shell
go run main.go -login
```
扫码登录，会得到 `access_key` 即可

> 填写配置文件 users.yaml

```shell
vim users.yaml
```
[参考下文](#配置文件)编辑 `users.yaml`

> 运行主程序

```shell
go run main.go
```

> 无视定时任务立刻运行一次

```shell
go run main.go -start
```

### 配置文件

```yaml
USERS:
  - access_key: XXXXXX # 注意冒号后的空格 否则会读取失败 英文冒号
    allowed_uid: #123,666 # 白名单UID,填了后将覆盖配置只打卡，点赞，分享这些用户的勋章 用英文逗号分隔 不填则不限制
    banned_uid: 789,100 # 黑名单UID,填了后将不会打卡，点赞，分享 用英文逗号分隔 不填则不限制
    push_name: PUSH_DEER_SAMPLE # 推送服务，留空表示不需要推送
  - access_key:
    allowed_uid:
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
  share: 5 # 【已废弃】分享间隔时间，单位秒，设置为0不分享
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
  - name: "TELEGRAM_SAMPLE"
    token: "<YOUR-TELEGRAM-CHATID>"
    type: "telegram"
    url: "https://api.telegram.org/bot<YOUR-BOT-TOKEN-HERE>/sendMessage"
  # 推送服务，每日打卡成功或报错日志推送
  # 目前支持PushDeer, PushPlus, Telegram
```

请务必严格填写，否则程序将读取失败，可以在这里 [YAML、YML 在线编辑器(格式化校验)-BeJSON.com](https://www.bejson.com/validators/yaml_editor/) 验证你填的 yaml 是否正确
