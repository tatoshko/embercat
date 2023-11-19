package wednesday

import (
    "embercat/assets"
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "time"
)

func Check(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    chatID := update.Message.Chat.ID
    logger := getLogger("CHECK")

    if time.Now().Weekday() == time.Wednesday {
        redis := redis2.GetClient()
        if redis == nil {
            return
        }
        defer redis.Close()

        var frog string
        if frog, err = redis.SRandMember(REDIS_KEY).Result(); err != nil {
            logger(err.Error())
        }

        msg := tgbotapi.NewPhotoShare(chatID, frog)
        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    } else {
        box := assets.GetBox()
        if b, err := box.Bytes(NO_WEDNESDAY); err != nil {
            logger(err.Error())
        } else {
            msg := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: NO_WEDNESDAY, Bytes: b})
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
        }
    }
}
