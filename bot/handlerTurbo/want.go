package handlerTurbo

import (
    "embercat/bot/core"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
    "strconv"
    "strings"
)

func Want(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    logger := getLogger("WANT")

    chatID := update.Message.Chat.ID
    args := update.Message.CommandArguments()

    validateArgs := regexp.MustCompile(`^\d{3}$`)

    if !validateArgs.MatchString(args) {
        msg := tgbotapi.NewMessage(chatID, "Неверный номер вкладыша, должно быть три цифры, например: 001")

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    var liner Liner
    if liner, err = NewLinerFromString(args); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не могу понять, что ты хочешь.\n<code>%s</code>", err.Error()))
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s хочет получить в дар вкладыш <b>%s</b>", core.GetUserName(update.Message.From, true), liner.ToString()))
    msg.ParseMode = tgbotapi.ModeHTML

    button := tgbotapi.NewInlineKeyboardButtonData("Подарить", fmt.Sprintf("/wantans %s %d", liner.ToString(), update.Message.From.ID))
    row := tgbotapi.NewInlineKeyboardRow(button)
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row)

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}

func CallbackWant(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    logger := getLogger("WANTCALLBACK")

    query := update.CallbackQuery
    callback := tgbotapi.NewCallback(query.ID, query.Data)
    chatID := query.Message.Chat.ID

    if _, err := bot.AnswerCallbackQuery(callback); err != nil {
        logger(err.Error())
        return
    }

    data := strings.Split(strings.TrimLeft(query.Data, "/wantans "), " ")

    var liner Liner
    if liner, err = NewLinerFromString(data[0]); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Неверный номер вкладыша <b>%s</b>", data[0]))
        msg.ParseMode = tgbotapi.ModeHTML

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

    log.Printf("Trying to move liner '%s' from '%d' to '%d'", liner.ToString(), giver, recipient)

    if giver == recipient {
        msg := tgbotapi.NewMessage(chatID, "Сам у себя это как вообще?")

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    var giverCollection Collection
    if giverCollection, err = LoadCollection(giver); err != nil {
        msg := tgbotapi.NewMessage(chatID, "У тебя нет вкладышей, жмакай\n<code>/turbo@embercatbot</code>")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

    var recipientCollection Collection
    if recipientCollection, err = LoadCollection(recipient); err != nil {
        msg := tgbotapi.NewMessage(chatID, "Что-то сломалось")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

    if err = giverCollection.MoveTo(recipientCollection, liner); err != nil {
        logger(err.Error())
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не получилось передать вкладыш, сорян.\n%s", err.Error()))
        msg.ParseMode = tgbotapi.ModeHTML
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Вкладыш <b>%s</b> подарен!", liner.ToString()))
    msg.ParseMode = tgbotapi.ModeHTML
    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }

}
