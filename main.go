package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	"subinfobot/handler"
	"time"
)

var (
	version string
	commit  string
	logger  = log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime)
)

func main() {
	logger.Printf("Subbot %s start.", version)
	bot, err := tgbotapi.NewBotAPI(os.Args[1])
	if err != nil {
		logger.Panic(fmt.Sprintf("Connect failed. %s"), err)
	}
	bot.Debug = true
	logger.Printf("Connected with name %s.", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				if update.Message.Chat.IsPrivate() {
					msg.ParseMode = "html"
					msg.Text = "🌈欢迎使用订阅信息查看bot！\n\n 📖命令列表: \n/start 开始\n/get 获取订阅链接的详细信息\n/about 关于\n/version 查看版本\n\n欢迎加入<a href=\"https://t.me/paimonnodegroup\">@paimonnodegroup</a>来改善此bot!\n"
					_, err := handler.SendMsg(bot, &msg, &update)
					handler.HandleError(err)
				}
			case "version":
				if update.Message.Chat.IsPrivate() {
					msg.ParseMode = "html"
					msg.Text = fmt.Sprintf("<strong>Subinfo Bot</strong>\n\n<strong>版本:</strong><code>%s</code>\n<strong>Commit id:</strong><code>%s</code>", version, commit)
					_, err := handler.SendMsg(bot, &msg, &update)
					handler.HandleError(err)
				}
			case "get":
				msg.ParseMode = "html"
				commandSlice := strings.Split(update.Message.Text, " ")
				if len(commandSlice) < 2 {
					msg.Text = "❌参数错误，请检查后再试"
					msg.ReplyToMessageID = update.Message.MessageID
					res, err := handler.SendMsg(bot, &msg, &update)
					handler.HandleError(err)
					if err == nil {
						if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
							_, _ = handler.DelMsgWithTimeOut(10*time.Second, bot, res)
						}
					}
				} else if strings.HasPrefix(commandSlice[1], "http://") || strings.HasPrefix(commandSlice[1], "https://") {
					msg.Text = "🕰获取中..."
					msg.ReplyToMessageID = update.Message.MessageID
					sres, err := handler.SendMsg(bot, &msg, &update)
					handler.HandleError(err)
					if err == nil {
						err, sinf := getSinf(commandSlice[1])
						handler.HandleError(err)
						if err != nil {
							_, err := handler.EditMsg(fmt.Sprintf("<strong>❌获取失败</strong>\n\n获取订阅<code>%s</code>时发生错误:\n<code>%s</code>", sinf.Link, err), "html", bot, sres)
							handler.HandleError(err)
							if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
								_, _ = handler.DelMsgWithTimeOut(10*time.Second, bot, sres)
							}
						} else {
							var resMsg string
							if sinf.Expired == 0 && sinf.Available == 0 {
								resMsg = "✅该订阅有效"
							}
							if sinf.Expired == 2 || sinf.Available == 2 {
								resMsg = "❓该订阅状态未知"
							}
							if sinf.Expired == 1 || sinf.Available == 1 {
								resMsg = "❌该订阅不可用"
							}
							_, err = handler.EditMsg(fmt.Sprintf("<strong>%s</strong>\n<strong>订阅链接:</strong><code>%s</code>\n<strong>总流量:</strong><code>%s</code>\n<strong>剩余流量:</strong><code>%s</code>\n<strong>已上传:</strong><code>%s</code>\n<strong>已下载:</strong><code>%s</code>\n<strong>该订阅将于<code>%s</code>过期,%s</strong>", resMsg, sinf.Link, sinf.Total, sinf.DataRemain, sinf.Upload, sinf.Download, sinf.ExpireTime, sinf.TimeRemain), "html", bot, sres)
							handler.HandleError(err)
							if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
								_, _ = handler.DelMsgWithTimeOut(10*time.Second, bot, sres)
							}
						}
					}
				} else {
					msg.Text = "❌链接错误，请检查后再试"
					msg.ReplyToMessageID = update.Message.MessageID
					res, err := handler.SendMsg(bot, &msg, &update)
					handler.HandleError(err)
					if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
						_, _ = handler.DelMsgWithTimeOut(10*time.Second, bot, res)
					}
				}
			default:
			}
		}
	}
}
