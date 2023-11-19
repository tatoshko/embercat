package handlerDonate

import (
	"embercat/bot/types"
	redis2 "embercat/redis"
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Add(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	userID := update.Message.From.UserName

	if userID != "tatoshko" {
		msg.Text = getUserDonateMessage()
	} else {
		if info, err := addDonater(update.Message.CommandArguments()); err != nil {
			log.Printf("HandlerAdd error %s\n", err.Error())
			return
		} else {
			msg.Text = getDonateMessage(info)
		}
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("HandlerAdd error %s\n", err.Error())
	}
}

func addDonater(args string) (redis2.Z, error) {
	redis := redis2.GetClient()
	if redis == nil {
		return redis2.Z{}, errors.New("failed to get the redis client")
	}
	defer redis.Close()

	if info, err := newDonateInfo(args); err != nil {
		return info, err
	} else {
		redis.ZAdd(types.REDIS_SUPPORTERS_COLLECTION, info)
		return info, nil
	}
}
