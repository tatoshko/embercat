package handlerTurbo

import (
    redis2 "embercat/redis"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
    "strings"
)

func HandlerGive(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    chantID := update.Message.Chat.ID
    userID := update.Message.From.ID

    localCollectionKey := makeLocalKey(REDIS_KEY_TURBO_COLLECTION, userID)

    validCommand := regexp.MustCompile(`^/give\s\d{3}\s\@\w+$`)
    text := update.Message.Text

    if validCommand.MatchString(text) {
        parts := strings.Split(text, " ")

        liner := parts[1]
        recipient := parts[2]

        log.Printf("Trying to give %s to %s", liner, recipient)

        if stats, err := redis.ZRangeWithScores(localCollectionKey, 0, -1).Result(); err != nil {
            log.Printf("HandlerCollection ZRangeWithScores error %s", err.Error())
        } else {
            log.Printf("%v", stats)
        }
    } else {
        log.Printf("Incorrect command '%s'", text)
        msg := tgbotapi.NewMessage(chantID, "Неверный формат команды. Пример:\n<code>/give 001 @username</code>")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Send error %s", err.Error())
        }
    }
}
