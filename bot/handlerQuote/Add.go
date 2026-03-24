package handlerQuote

import (
    "embercat/bot/core"
    "embercat/bot/handlerQuote/service"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Add(bot *tba.BotAPI, update tba.Update) {
    var err error
    logger := getLogger("ADD")
    chatID := update.Message.Chat.ID

    var member tba.ChatMember
    if member, err = core.GetChatMember(bot, chatID, update.Message.From.ID); err != nil {
        logger(err.Error())
        return
    }

    if !canSave(member) {
        msg := tba.NewMessage(chatID, fmt.Sprintf("Слыш, пёс. Не только лишь все могут это делать. %s", chatID))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    quoteService := service.NewService(chatID)
    msg := tba.NewMessage(chatID, "")

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

func canSave(member tba.ChatMember) bool {
    if member.IsCreator() {
        return true
    }

    if member.IsAdministrator() {
        return member.CanPromoteMembers
    }

    return false
}
