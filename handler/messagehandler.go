package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func SendMsg(bot *tgbotapi.BotAPI,msg *tgbotapi.MessageConfig,update *tgbotapi.Update)(*tgbotapi.Message,error){
	res, err := bot.Send(msg)
	if err != nil {
		return &res,err
	}
	return &res,nil
}

func DelMsgWithTimeOut(duration time.Duration,bot *tgbotapi.BotAPI,msg *tgbotapi.Message)(*tgbotapi.APIResponse,error){
	time.Sleep(duration)
	conf := tgbotapi.NewDeleteMessage(msg.Chat.ID,msg.MessageID)
	res, err := bot.Request(conf)
	if err != nil {
		return res,err
	}
	return res,nil
}

func EditMsg(text string,parsemode string,bot *tgbotapi.BotAPI,msg *tgbotapi.Message)(*tgbotapi.APIResponse,error){
	conf := tgbotapi.NewEditMessageText(msg.Chat.ID,msg.MessageID,text)
	conf.ParseMode = parsemode
	res, err := bot.Request(conf)
	if err != nil {
		return res,err
	}
	return res,nil
}