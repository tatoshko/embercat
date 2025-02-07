package handlerQuote

import (
    "embercat/bot/core"
    "embercat/bot/handlerQuote/service"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Add(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("ADD")
    chatID := update.Message.Chat.ID

    var member tgbotapi.ChatMember
    if member, err = core.GetChatMember(bot, chatID, update.Message.From.ID); err != nil {
        logger(err.Error())
        return
    }

    if !canSave(member) {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Слыш, пёс. Не только лишь все могут это делать. %s", chatID))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    quoteService := service.NewService()
    msg := tgbotapi.NewMessage(chatID, "")

    if update.Message.ReplyToMessage == nil {
        msg.Text = "Нужно сделать реплай на сообщение"
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    quote := service.NewQuoteFromMessage(update.Message.ReplyToMessage)

    if err = quoteService.Add(quote); err != nil {
        logger(err.Error())
        msg.Text = fmt.Sprintf("Что-то не так: %s", err.Error())
    } else {
        msg.Text = fmt.Sprintf("Записал цитатку:\n%s", quote.ToString())
    }

    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}

func canSave(member tgbotapi.ChatMember) bool {
    if member.IsCreator() {
        return true
    }

    if member.IsAdministrator() {
        return member.CanPromoteMembers
    }

    return false
}
