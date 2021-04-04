package core

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
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

func InitBot(token, hook string) (tbot *TBot) {
    tbot = &TBot{
        Bot:        nil,
        commandMsg: regexp.MustCompile(`^\/(?P<command>\w+)\s*(?P<data>.*)$`),
        HANDLERS:   make(map[string]CommandHandler),
    }

    tbot.RegisterHandler("set", handleSet)
    tbot.RegisterHandler("get", handleGet)


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

        if bot.commandMsg.MatchString(text) {
            match := reSubMatchMap(bot.commandMsg, text)

            log.Printf("Command: '%s', data: '%s'", match["command"], match["data"])

            if handler, found := bot.HANDLERS[match["command"]]; found {
                go handler(bot.Bot, update, match["data"])
            }
        }
    }
}