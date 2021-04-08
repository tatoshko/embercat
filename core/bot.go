package core

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
)

type CommandHandler func(*tba.BotAPI, tba.Update)

type TBot struct {
    Bot *tba.BotAPI
    commandMsg *regexp.Regexp
    updates tba.UpdatesChannel

    HANDLERS map[string]CommandHandler
}

func StartBot(token, hook string) {
    bot := &TBot{
        Bot:        nil,
        commandMsg: regexp.MustCompile(`^/(?P<command>\w+)\s*(?P<data>.*)$`),
        HANDLERS:   make(map[string]CommandHandler),
    }

    bot.RegisterHandler("thread", handleThread)
    bot.RegisterHandler("day", handleWednesday)


    if Bot, err := tba.NewBotAPI(token); err == nil {
        bot.Bot = Bot

        if _, err := Bot.SetWebhook(tba.NewWebhook(hook + "/" + token)); err != nil {
            log.Printf("SetHoook error %s", err.Error())
        }

        bot.updates = Bot.ListenForWebhook("/" + Bot.Token)

        go bot.Watch()
    } else {
        log.Fatalf("NewAPIBot error %s", err.Error())
    }
}

func (bot *TBot) RegisterHandler(name string, f CommandHandler) {
    bot.HANDLERS[name] = f
}

func (bot *TBot) UnregisterHandler(name string) {
    delete(bot.HANDLERS, name)
}

func (bot *TBot) Watch() {
    for update := range bot.updates {
        if update.Message == nil {
            continue
        }

        text := update.Message.Text

        if bot.commandMsg.MatchString(text) {
            match := reSubMatchMap(bot.commandMsg, text)

            if handler, found := bot.HANDLERS[match["command"]]; found {
                log.Printf("Command: '%s', data: '%s'", match["command"], match["data"])
                go handler(bot.Bot, update)
            }
        }
    }
}