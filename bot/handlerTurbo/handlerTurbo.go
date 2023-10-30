package handlerTurbo

import (
	"embercat/assets"
	redis2 "embercat/redis"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"time"
)

func HandlerTurbo(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	redis := redis2.GetClient()
	if redis == nil {
		return
	}
	defer redis.Close()

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	localCollectionKey := makeLocalKey(REDIS_KEY_TURBO_COLLECTION, userID)
	localDayKey := makeLocalKey(REDIS_KEY_TURBO_DAY, userID)

	currentDateString := time.Now().Format("2006-Jan-02")

	var isTodayGot bool
	var err error
	if isTodayGot, err = redis.SIsMember(localDayKey, currentDateString).Result(); err != nil {
		log.Printf("HandlerTurbo SIsMember error %s", err.Error())
		return
	}

	if isTodayGot {
		msg := tgbotapi.NewMessage(chatID, "Можно только одну жвачку в день")
		if _, err = api.Send(msg); err != nil {
			log.Printf("HandlerTurbo api.Send error %s", err.Error())
		}

		return
	} else {
		if _, err = redis.SAdd(localDayKey, currentDateString).Result(); err != nil {
			log.Printf("HandlerTurbo SAdd error %s", err.Error())
			return
		}
	}

	box := assets.GetBox()

	rand.Seed(time.Now().UnixNano())
	n := fmt.Sprintf("%03d", rand.Intn(TOTAL_PICTURES))
	filename := fmt.Sprintf(TURBO_FILENAME_KEY, n)

	if _, err := redis.ZIncrBy(localCollectionKey, 1, n).Result(); err != nil {
		log.Printf("HandlerTurbo ZIncrBy error %s", err.Error())
		return
	}

	if b, err := box.Bytes(filename); err != nil {
		log.Printf("HandlerTurbo box.Bytes error %s", err.Error())
		return
	} else {
		msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: filename, Bytes: b})

		if _, err = api.Send(msg); err != nil {
			log.Printf("HandlerTurbo api.Send error %s", err.Error())
		}
	}

}
