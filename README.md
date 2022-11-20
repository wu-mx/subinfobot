# Subinfo Bot

[![Action-Release](https://github.com/wu-mx/subinfobot/actions/workflows/Build.yml/badge.svg)](https://github.com/wu-mx/subinfobot/actions/workflows/Build.yml/)
[![Download](https://img.shields.io/github/downloads/wu-mx/subinfobot/total.svg)](https://github.com/wu-mx/subinfobot/releases)
#### 一个由Go编写的开源轻量订阅查询Telegram Bot。

### 快速开始
您可以通过直接下载编译好的二进制文件或自行编译来使用Subinfo Bot。
<br>[Release下载](https://github.com/wu-mx/subinfobot/releases/tag/v0.0.1)

#### 自行编译
```shell
git clone https://github.com/wu-mx/subinfobot/
go get 
go build
```
####快速启动
在启动程序参数中加入bot token即可。<br>
Subinfo Bot时区由系统时区决定，错误的时区将导致获取结果出现偏差，更改方法请自行Google。
```shell
./subinfobot 1234567890:AABCDEFGHIJKLMNOPQRSTUVWXYZ-abcdefg #你的bot token
```

### 使用
您可以在私聊中对bot直接发送带有订阅链接的消息，bot将会自动提取文本中的订阅链接。
```
10PB不限时机场订阅链接：https://subapi.paimon.gq/api/v1/client/subscribe?token=1145141919810subinfobotyyds

✅该订阅有效
订阅链接:https://subapi.paimon.gq/api/v1/client/subscribe?token=1145141919810subinfobotyyds
总流量:10.00PB
剩余流量:10.00PB
已上传:11.45GB
已下载:19.19GB
该订阅将于2023-11-04 05:01:04 +0800 CST过期,距离到期还有191天09小时08分10秒
```

### Credits
[go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) <br>
[Dreamacro/clash](https://github.com/Dreamacro/clash)

### Licence
[MIT](https://github.com/wu-mx/subinfobot/LICENSE.txt)