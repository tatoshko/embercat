package core

import (
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
    "strings"
)

type CommandHandler func(*tba.BotAPI, tba.Update, string)

type TBot struct {
    Bot *tba.BotAPI
    commandMsg *regexp.Regexp
    updates tba.UpdatesChannel

    HANDLERS map[string]CommandHandler
    REPLAYS map[string]string
}

var err error
var Bot *tba.BotAPI

func InitBot(config Config) (tbot *TBot) {
    token := config.Token
    hook := config.Hook

    tbot = &TBot{
        Bot:        nil,
        commandMsg: regexp.MustCompile(`^\/(?P<command>\w+)\s*(?P<data>.*)$`),
        HANDLERS:   make(map[string]CommandHandler),
        REPLAYS:    make(map[string]string),
    }

    initStorage(config.DB)

    tbot.RegisterHandler("set", handleSet(config.DB))
    tbot.RegisterHandler("get", handleGet(config.DB))
    tbot.RegisterHandler("thread", handleThread)
    tbot.RegisterHandler("rebus", handleRebus(tbot))
    tbot.RegisterHandler("wednesday", handleWednesday)
    //tbot.RegisterHandler("twoch", handle2ch)


    if Bot, err = tba.NewBotAPI(token); err != nil {
        panic(err)
    }

    tbot.Bot = Bot

    Bot.SetWebhook(tba.NewWebhook(hook + "/" + token))
    tbot.updates = Bot.ListenForWebhook("/" + Bot.Token)

    return
}

func (bot *TBot) RegisterHandler(name string, f CommandHandler) {
    bot.HANDLERS[name] = f
}

func (bot *TBot) UnregisterHandler(name string) {
    delete(bot.HANDLERS, name)
}

func (bot *TBot) RegisterReplay(id, answer string) {
    bot.REPLAYS[id] = answer
}

func (bot *TBot) UnregisterReplay(id string) {
    delete(bot.REPLAYS, id)
}

func (bot *TBot) Watch() {
    for update := range bot.updates {
        text := update.Message.Text

        if bot.commandMsg.MatchString(text) {
            match := reSubMatchMap(bot.commandMsg, text)

            if handler, found := bot.HANDLERS[match["command"]]; found {
                log.Printf("Command: '%s', data: '%s'", match["command"], match["data"])
                go handler(bot.Bot, update, strings.TrimSpace(match["data"]))
            }
        }

        if update.Message.ReplyToMessage == nil {
            continue
        }

        id := fmt.Sprintf("id_%d", update.Message.ReplyToMessage.MessageID)
        if value, found := bot.REPLAYS[id]; found {
            if text == value {
                msg := tba.NewMessage(update.Message.Chat.ID, "Верно!")
                msg.ReplyToMessageID = update.Message.MessageID
                bot.Bot.Send(msg)
                bot.UnregisterReplay(id)
            }
        }
    }
}