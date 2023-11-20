package handlerWednesday

import (
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func Save(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    chatID := update.Message.Chat.ID
    logger := getLogger("SAVE")

    // Check. Message is replay with photo
    reply := update.Message.ReplyToMessage
    if reply == nil || reply.Photo == nil {
        msg := tgbotapi.NewMessage(chatID, "Вызов этой команды возможен только в ответ на сообщение с картинкой")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    // Check. Current user is Administrator
    var admins []tgbotapi.ChatMember
    config := tgbotapi.ChatConfig{ChatID: chatID}
    if admins, err = bot.GetChatAdministrators(config); err != nil {
        logger(err.Error())
    }

    if !checkAdmin(admins, update.Message.From) {
        msg := tgbotapi.NewMessage(chatID, "Слыш, пёс. Не только лишь все могут это делать.")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    // Get photoID
    var photoID string
    for _, pic := range *reply.Photo {
        photoID = pic.FileID
        break
    }

    // Connect to redis and save
    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    if _, err = redis.SAdd(REDIS_KEY, photoID).Result(); err != nil {
        logger(err.Error())
        return
    }

    msg := tgbotapi.NewMessage(chatID, "Схоронил, под сердечком прям")
    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}

func checkAdmin(admins []tgbotapi.ChatMember, user *tgbotapi.User) bool {
    for _, admin := range admins {
        log.Printf("DEBUG %d, %d", admin.User.ID, user.ID)

        if admin.User.ID != user.ID {
            continue
        }

        log.Printf("DEBUG CAN PROMOTE %q", admin)
        if !admin.CanPromoteMembers {
            continue
        }

        return true
    }

    return false
}
