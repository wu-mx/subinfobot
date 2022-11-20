package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"regexp"
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
		logger.Panic(fmt.Sprintf("Connect failed. %s", err))
	}
	bot.Debug = true
	logger.Printf("Connected with name %s.", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if !update.Message.IsCommand() {
				if update.Message.Chat.IsPrivate() {
					linkReg := regexp.MustCompile("(http|https){0,1}://[^\\x{4e00}-\\x{9fa5}\\n\\r\\s]{3,}")
					if linkReg.MatchString(update.Message.Text) {
						slice := linkReg.FindAllStringSubmatch(update.Message.Text, -1)
						subInfoMsg(slice[0][0], &update, bot, &msg)
					} else {
						msg.Text = "❌没有在你发送的内容中找到任何有效信息哦！"
						msg.ReplyToMessageID = update.Message.MessageID
						_, err := handler.SendMsg(bot, &msg)
						handler.HandleError(err)
					}
				}
			}
			switch update.Message.Command() {
			case "start":
				if update.Message.Chat.IsPrivate() {
					msg.ParseMode = "html"
					msg.Text = "🌈欢迎使用订阅信息查看bot！\n\n 📖命令列表: \n/start 开始\n/get 获取订阅链接的详细信息\n/about 关于\n/version 查看版本\n\n欢迎加入<a href=\"https://t.me/paimonnodegroup\">@paimonnodegroup</a>来改善此bot!\n"
					_, err := handler.SendMsg(bot, &msg)
					handler.HandleError(err)
				}
			case "version":
				if update.Message.Chat.IsPrivate() {
					msg.ParseMode = "html"
					msg.Text = fmt.Sprintf("<strong>Subinfo Bot</strong>\n\n<strong>版本:</strong><code>%s</code>\n<strong>Commit id:</strong><code>%s</code>", version, commit)
					_, err := handler.SendMsg(bot, &msg)
					handler.HandleError(err)
				}
			case "about":
				msg.ParseMode = "html"
				msg.Text = fmt.Sprintf("<strong>Subinfo Bot %s</strong>\n\nSubinfo Bot是一款由Golang编写的开源轻量订阅查询Bot。体积小巧，无需任何第三方运行时，即点即用。\n\n<strong>Github:<a href=\"https://github.com/wu-mx/subinfobot\">https://github.com/wu-mx/subinfobot</a></strong>\n<strong>讨论群组:<a href=\"https://t.me/paimonnodegroup\">@paimonnodegroup</a></strong>", version)
				_, err := handler.SendMsg(bot, &msg)
				handler.HandleError(err)
			case "get":
				msg.ParseMode = "html"
				commandSlice := strings.Split(update.Message.Text, " ")
				if len(commandSlice) < 2 {
					msg.Text = "❌参数错误，请检查后再试"
					msg.ReplyToMessageID = update.Message.MessageID
					res, err := handler.SendMsg(bot, &msg)
					handler.HandleError(err)
					if err == nil {
						if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
							_, _ = handler.DelMsgWithTimeOut(10*time.Second, bot, res)
						}
					}
				} else if strings.HasPrefix(commandSlice[1], "http://") || strings.HasPrefix(commandSlice[1], "https://") {
					subInfoMsg(commandSlice[1], &update, bot, &msg)
				} else {
					msg.Text = "❌链接错误，请检查后再试"
					msg.ReplyToMessageID = update.Message.MessageID
					res, err := handler.SendMsg(bot, &msg)
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
