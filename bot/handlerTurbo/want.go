package handlerTurbo

import (
    "embercat/bot/core"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "regexp"
    "strconv"
    "strings"
)

func Want(bot *tba.BotAPI, update tba.Update) {
    var err error

    logger := getLogger("WANT")

    chatID := update.Message.Chat.ID
    args := update.Message.CommandArguments()

    validateArgs := regexp.MustCompile(`^\d{3}$`)

    if !validateArgs.MatchString(args) {
        msg := tba.NewMessage(chatID, "Неверный номер вкладыша, должно быть три цифры, например: 001")

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    var liner Liner
    if liner, err = NewLinerFromString(args); err != nil {
        msg := tba.NewMessage(chatID, fmt.Sprintf("Не могу понять, что ты хочешь.\n<code>%s</code>", err.Error()))
        msg.ParseMode = tba.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

    msg := tba.NewMessage(chatID, fmt.Sprintf("%s хочет получить в дар вкладыш <b>%s</b>", core.GetUserName(update.Message.From, true), liner.ToString()))
    msg.ParseMode = tba.ModeHTML

    action := fmt.Sprintf("/wantans %d %d", liner.ID, update.Message.From.ID)
    button := tba.NewInlineKeyboardButtonData("Подарить", action)
    row := tba.NewInlineKeyboardRow(button)
    msg.ReplyMarkup = tba.NewInlineKeyboardMarkup(row)

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}

func CallbackWant(bot *tba.BotAPI, update tba.Update) {
    var err error

    logger := getLogger("WANTCALLBACK")

    query := update.CallbackQuery
    callback := tba.NewCallback(query.ID, query.Data)
    chatID := query.Message.Chat.ID

    if _, err := bot.Request(callback); err != nil {
        logger(err.Error())
        return
    }

    data := strings.Split(strings.TrimLeft(query.Data, "/wantans "), " ")

    var liner Liner
    if liner, err = NewLinerFromString(data[0]); err != nil {
        logger(err.Error())
        msg := tba.NewMessage(chatID, fmt.Sprintf("Неверный номер вкладыша <b>%s</b>", data[0]))
        msg.ParseMode = tba.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var recipient int64
    if recipient, err = strconv.ParseInt(data[1], 10, 32); err != nil {
        logger(err.Error())
    }

    giver := int64(query.From.ID)

    logger(fmt.Sprintf("Trying to move liner '%s' from '%d' to '%d'", liner.ToString(), giver, recipient))

    if giver == recipient {
        msg := tba.NewMessage(chatID, "Сам у себя это как вообще?")

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    var giverCollection Collection
    if giverCollection, err = LoadCollection(giver); err != nil {
        msg := tba.NewMessage(chatID, "У тебя нет вкладышей, жмакай\n<code>/turbo@embercatbot</code>")
        msg.ParseMode = tba.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

    var recipientCollection Collection
    if recipientCollection, err = LoadCollection(recipient); err != nil {
        msg := tba.NewMessage(chatID, "Что-то сломалось")
        msg.ParseMode = tba.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

    if err = giverCollection.MoveTo(recipientCollection, liner); err != nil {
        logger(err.Error())
        msg := tba.NewMessage(chatID, fmt.Sprintf("Не получилось передать вкладыш, сорян.\n%s", err.Error()))
        msg.ParseMode = tba.ModeHTML
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    msg := tba.NewMessage(chatID, fmt.Sprintf("Вкладыш <b>%s</b> подарен!", liner.ToString()))
    msg.ParseMode = tba.ModeHTML
    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }

}
