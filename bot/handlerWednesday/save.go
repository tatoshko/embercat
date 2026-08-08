package handlerWednesday

import (
    "embercat/bot/core"
    "embercat/bot/core/loader"
    "fmt"
    "github.com/corona10/goimagehash"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "image"
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

    // Check permissions
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

    // Get file and compute hash
    var fileURL string
    if fileURL, err = bot.GetFileDirectURL(photoID); err != nil {
        msg := tba.NewMessage(chatID, "Не смог получить картинку")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var img *image.RGBA
    if img, err = loader.LoadPicByURL(fileURL); err != nil {
        logger("LoadPicByURL", err.Error())
        return
    }

    var hash *goimagehash.ImageHash
    if hash, err = goimagehash.PerceptionHash(img); err != nil {
        logger("PerceptionHash", err.Error())
        return
    }
    vector := core.HashToVector(hash.GetHash())

    frogService := NewFrogService()

    // Check exists
    if exists, frog, err := frogService.IsExists(vector); err != nil {
        msg := tba.NewMessage(chatID, fmt.Sprintf("При сравнении жаб что-то пошло не так"))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    } else if exists {
        msg := tba.NewPhoto(chatID, tba.FileID(frog.PhotoId))
        msg.Caption = fmt.Sprintf("Похожая жаба уже есть")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    // Save new
    if err = frogService.Add(photoID, vector); err != nil {
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
