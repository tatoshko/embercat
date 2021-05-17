package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

var (
    API      *tba.BotAPI
    HANDLERS map[string]CommandHandler
)

type CommandHandler func(api *tba.BotAPI, update tba.Update)

func Start(token, hook string) {
    if API, err := tba.NewBotAPI(token); err == nil {
        if _, err := API.SetWebhook(tba.NewWebhook(hook + "/" + token)); err != nil {
            log.Printf("SetHoook error %s\n", err.Error())
        }

        registerHandlers()

        updates := API.ListenForWebhook("/" + API.Token)

        for update := range updates {
            if update.Message == nil {
                continue
            }

            message := update.Message

            if message.IsCommand() {
                if handler, found := HANDLERS[message.Command()]; found {
                    log.Printf("Command: '%s', data: '%s'\n", message.Command(), message.CommandArguments())
                    go handler(API, update)
                }
            }
        }
    } else {
        log.Fatalf("NewAPIBot error %s\n", err.Error())
    }
}

func registerHandlers() {
    HANDLERS = make(map[string]CommandHandler)

    HANDLERS["thread"] = handleThread
    HANDLERS["day"] = handleWednesday
}
