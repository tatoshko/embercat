package handlerTurbo

import (
    "embercat/bot/helpers"
    redisServ "embercat/redis"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
    "strconv"
    "strings"
)

func HandlerWant(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    chatID := update.Message.Chat.ID
    args := update.Message.CommandArguments()

    validateArgs := regexp.MustCompile(`^\d{3}$`)

    if !validateArgs.MatchString(args) {
        msg := tgbotapi.NewMessage(chatID, "Неверный номер вкладыша, должно быть трицифры, например: 001")

        if _, err := bot.Send(msg); err != nil {
            logErr(err)
        }

        return
    }

    var liner Liner
    if liner, err = NewLinerFromString(args); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не могу понять, что ты хочешь.\n<code>%s</code>", err.Error()))
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logErr(err)
        }
    }

    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s хочет получить в дар вкладыш <b>%s</d>", helpers.GetUserName(update.Message.From, true), liner.ID))
    msg.ParseMode = tgbotapi.ModeHTML

    button := tgbotapi.NewInlineKeyboardButtonData("Подарить", fmt.Sprintf("/wantans %s %d", liner.ID, update.Message.From.ID))
    row := tgbotapi.NewInlineKeyboardRow(button)
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row)

    if _, err := bot.Send(msg); err != nil {
        logErr(err)
    }

}

func CallbackWant(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    redisInst := redisServ.GetClient()
    if redisInst == nil {
        return
    }
    defer redisInst.Close()

    query := update.CallbackQuery
    callback := tgbotapi.NewCallback(query.ID, query.Data)
    chatID := update.CallbackQuery.Message.Chat.ID

    if _, err := bot.AnswerCallbackQuery(callback); err != nil {
        logErr(err)
        return
    }

    data := strings.Split(strings.TrimLeft(query.Data, "/wantans "), " ")

    var liner Liner
    if liner, err = NewLinerFromString(data[0]); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Неверный номер вкладыша <b>%s</b>", data[0]))
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logErr(err)
        }
        return
    }

    var recipient int64
    if recipient, err = strconv.ParseInt(data[1], 10, 32); err != nil {
        logErr(err)
    }

    log.Printf("Trying to move liner '%s' from '%d' to '%d'", liner.ID, update.Message.From.ID, recipient)

    var giverCollection Collection
    if giverCollection, err = LoadCollection(redisInst, int64(update.Message.From.ID)); err != nil {
        msg := tgbotapi.NewMessage(chatID, "У тебя нет вкладышей, жмакай\n<code>/turbo@embercatbot</code>")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logErr(err)
        }
    }

    var recipientCollection Collection
    if recipientCollection, err = LoadCollection(redisInst, recipient); err != nil {
        msg := tgbotapi.NewMessage(chatID, "Что-то сломалось")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logErr(err)
        }
    }

    if err = giverCollection.MoveTo(recipientCollection, liner); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не получилось передать вкладыш, сорян.\n%s", err.Error()))
        msg.ParseMode = tgbotapi.ModeHTML
        if _, err = bot.Send(msg); err != nil {
            logErr(err)
        }
    }

    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Вкладыш <b>%s</b> подарен!", liner.ID))
    if _, err := bot.Send(msg); err != nil {
        logErr(err)
    }

}

func logErr(err error) {
    log.Printf("turbo.HandlerWant error %s", err.Error())
}