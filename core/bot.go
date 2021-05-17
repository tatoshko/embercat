package core

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
)

type CommandHandler func(api *tba.BotAPI, data string, chatID int64, update tba.Update)

type TBot struct {
    API        *tba.BotAPI
    commandMsg *regexp.Regexp
    updates    tba.UpdatesChannel

    HANDLERS map[string]CommandHandler
}

func StartBot(token, hook string) {
    bot := &TBot{
        API:        nil,
        commandMsg: regexp.MustCompile(`^/(?P<command>\w+)\s*(?P<data>.*)$`),
        HANDLERS:   make(map[string]CommandHandler),
    }

    bot.RegisterHandler("thread", handleThread)
    bot.RegisterHandler("day", handleWednesday)

    if API, err := tba.NewBotAPI(token); err == nil {
        bot.API = API

        if _, err := API.SetWebhook(tba.NewWebhook(hook + "/" + token)); err != nil {
            log.Printf("SetHoook error %s", err.Error())
        }

        bot.updates = API.ListenForWebhook("/" + API.Token)
        bot.Watch()
    } else {
        log.Fatalf("NewAPIBot error %s", err.Error())
    }
}

func (bot *TBot) RegisterHandler(command string, f CommandHandler) {
    bot.HANDLERS[command] = f
}

func (bot *TBot) Watch() {
    for update := range bot.updates {
        if update.Message == nil {
            continue
        }

        text := update.Message.Text

        if bot.commandMsg.MatchString(text) {
            match := ParseCommand(bot.commandMsg, text)

            if handler, found := bot.HANDLERS[match["command"]]; found {
                log.Printf("Command: '%s', data: '%s'", match["command"], match["data"])
                go handler(bot.API, match["data"], update.Message.Chat.ID, update)
            }
        }
    }
}
