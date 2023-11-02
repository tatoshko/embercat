package handlerTurbo

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
    "strings"
)

func HandlerGive(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chantID := update.Message.Chat.ID

    validCommand := regexp.MustCompile(`^/give\s\d{3}\s\@\w+$`)
    text := update.Message.Text

    if validCommand.MatchString(text) {
        parts := strings.Split(text, " ")

        log.Printf("%v", parts)
    } else {
        log.Printf("Incorrect command %s", text)
        msg := tgbotapi.NewMessage(chantID, "Неверный формат команды. Пример:\n<code>/give 001 @username</code>")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Send error %s", err.Error())
        }
    }
}
