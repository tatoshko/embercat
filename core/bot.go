package core

import (
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
    }

    initStorage(config.DB)

    tbot.RegisterHandler("set", handleSet(config.DB))
    tbot.RegisterHandler("get", handleGet(config.DB))
    tbot.RegisterHandler("thread", handleThread)
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

func (bot *TBot) Watch() {
    for update := range bot.updates {
        text := update.Message.Text

        log.Printf("%q - %q", update.Message.Text, update.Message.ReplyToMessage.Text)

        if bot.commandMsg.MatchString(text) {
            match := reSubMatchMap(bot.commandMsg, text)

            if handler, found := bot.HANDLERS[match["command"]]; found {
                log.Printf("Command: '%s', data: '%s'", match["command"], match["data"])
                go handler(bot.Bot, update, strings.TrimSpace(match["data"]))
            }
        }
    }
}