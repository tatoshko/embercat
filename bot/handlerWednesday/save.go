package handlerWednesday

import (
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

    // Connect to DB and save
    pg := pgsql.GetClient()
    if pg == nil {
        return
    }
    defer pg.Close()

    q := `insert into frog (photoId) VALUES ($1)`
    if _, err = pg.Exec(q, photoID); err != nil {
        logger("unable to save")
        return
    }

    msg := tgbotapi.NewMessage(chatID, "Схоронил, под сердечком прям")
    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}

func checkAdmin(admins []tgbotapi.ChatMember, user *tgbotapi.User) bool {
    for _, admin := range admins {
        if admin.IsCreator() {
            return true
        }

        if admin.User.ID != user.ID {
            continue
        }

        if !admin.CanPromoteMembers {
            continue
        }

        return true
    }

    return false
}
