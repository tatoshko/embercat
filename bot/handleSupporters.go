package bot

import (
	"fmt"
	"log"

	"embercat/bot/types"
	redis2 "embercat/redis"
)

func handleSupporters(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	redis := redis2.GetClient()
	if redis == nil {
		return
	}
	defer redis.Close()

	if supporters, err := redis.SMembers(types.REDIS_SUPPORTERS_COLLECTION).Result(); err != nil {
		log.Printf("handleSupporters error %s\n", err.Error())
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, getMessage(supporters))
		if _, err := bot.Send(msg); err != nil {
			log.Printf("handleSupporters error %s\n", err.Error())
		}
	}
}

func getMessage(supporters []string) string {
	msg := "Список донатеров:\n\n"
	for i, s := range supporters {
		msg += fmt.Sprintf("%d. %s\n", i+1, s)
	}
	return msg
}
