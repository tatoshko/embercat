package handlerDonate

import (
	"embercat/bot/types"
	redis2 "embercat/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func HandlerShow(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	redis := redis2.GetClient()
	if redis == nil {
		return
	}
	defer redis.Close()

	donates := redis.ZRangeWithScores(types.REDIS_SUPPORTERS_COLLECTION, 0, -1)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, getDonatesList(donates))
	if _, err := bot.Send(msg); err != nil {
		log.Printf("HandlerShow error %s\n", err.Error())
	}
}