package handlerTurbo

import (
    redisServ "embercat/redis"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "gopkg.in/redis.v3"
    "log"
    "regexp"
    "strings"
)

func HandlerGive(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    redisInst := redisServ.GetClient()
    if redisInst == nil {
        return
    }
    defer redisInst.Close()

    chantID := update.Message.Chat.ID
    userID := update.Message.From.ID

    localCollectionKey := makeLocalKey(REDIS_KEY_TURBO_COLLECTION, userID)

    validCommand := regexp.MustCompile(`^/give\@embercatbot\s\d{3}\s\@\w+$`)
    text := update.Message.Text

    if !validCommand.MatchString(text) {
        log.Printf("Incorrect command '%s'", text)

        msg := tgbotapi.NewMessage(chantID, "Неверный формат команды. Пример:\n<code>/give 001 @username</code>")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Send error %s", err.Error())
        }

        return
    }

    parts := strings.Split(text, " ")

    liner := parts[1]
    recipient := parts[2]

    log.Printf("Trying to give %s to %s", liner, recipient)

    var stats []redis.Z
    if stats, err = redisInst.ZRangeWithScores(localCollectionKey, 0, -1).Result(); err != nil {
        log.Printf("HandlerCollection ZRangeWithScores error %s", err.Error())

        msg := tgbotapi.NewMessage(chantID, "У тебя нет вкладышей, жмакай\n<code>/turbo@embercatbot</code>")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Send error %s", err.Error())
        }

        return
    }

    score := GetScore(stats, liner)

    if score <= 0 {
        log.Printf("You don't have %s", liner)

        msg := tgbotapi.NewMessage(chantID, fmt.Sprintf("У тебя нет вклладыша <b>%s</b>", liner))
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Send error %s", err.Error())
        }

        return
    }

    var member tgbotapi.ChatMember
    if member, err = GetChatMember(bot, recipient); err != nil {
        log.Printf("Get chatMember error %s", err.Error())

        msg := tgbotapi.NewMessage(chantID, fmt.Sprintf("Не могу найти пользователя <b>%s</b>", recipient))
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Send error %s", err.Error())
        }

        return
    }

    log.Printf("Chat member %v", member.ID)

    // пользователь играет
    // отнять у текущего
    // передать новому

}
