package handlerWednesday

import (
    "embercat/bot/core"
    "embercat/pgsql"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Save(bot *tba.BotAPI, update tba.Update) {
    var err error
    chatID := update.Message.Chat.ID
    logger := getLogger("SAVE")

    // Check. Message is replay with photo
    replay := update.Message.ReplyToMessage
    if replay == nil || replay.Photo == nil {
        msg := tba.NewMessage(chatID, "Вызов этой команды возможен только в ответ на сообщение с картинкой")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var member tba.ChatMember
    if member, err = core.GetChatMember(bot, chatID, update.Message.From.ID); err != nil {
        logger(err.Error())
        return
    }

    if !canSave(member) {
        msg := tba.NewMessage(chatID, "Слыш, пёс. Не только лишь все могут это делать.")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    // Get photoID
    var photoID string
    for _, pic := range replay.Photo {
        photoID = pic.FileID
        break
    }

    // Connect to DB and save
    pg := pgsql.GetClient()

    q := `insert into frog (photoId) VALUES ($1)`
    if _, err = pg.Exec(q, photoID); err != nil {
        logger("unable to save: ", err.Error())

        msg := tba.NewMessage(chatID, fmt.Sprintf("Скорее всего такая жаба уже есть или что-то пошло не так"))

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    msg := tba.NewMessage(chatID, "Схоронил, под сердечком прям")
    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}

func canSave(member tba.ChatMember) bool {
    if member.IsCreator() {
        return true
    }

    if member.IsAdministrator() {
        return member.CanPromoteMembers
    }

    return false
}
