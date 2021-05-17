package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

var (
    Commands = make(map[string]Handler)
    Callbacks = make(map[string]Handler)
)

type Handler func(api *tba.BotAPI, update tba.Update)

func Start(token, hook string) {
    if API, err := tba.NewBotAPI(token); err == nil {
        if _, err := API.SetWebhook(tba.NewWebhook(hook + "/" + token)); err != nil {
            log.Printf("SetHoook error %s\n", err.Error())
        }

        registerCommands()

        updates := API.ListenForWebhook("/" + API.Token)

        for update := range updates {
            if update.Message != nil {
                message := update.Message

                if message.IsCommand() {
                    if handler, found := Commands[message.Command()]; found {
                        log.Printf("MessageID: '%d', Command: '%s', data: '%s'\n", message.MessageID, message.Command(), message.CommandArguments())
                        go handler(API, update)
                    }
                }
            } else if update.CallbackQuery != nil {
                data := update.CallbackQuery.Data

                if handler, found := Callbacks[data]; found {
                    handler(API, update)
                }
            }
        }
    } else {
        log.Fatalf("NewAPIBot error %s\n", err.Error())
    }
}

func registerCommands() {
    Commands["thread"] = handleThread
    Commands["day"] = handleWednesday
    Commands["btn"] = handlerBtn
}

func registerCallback(id string, f Handler) {
    Callbacks[id] = f
}

func unregisterCallback(id string) {
    if _, found := Callbacks[id]; found {
        delete(Callbacks, id)
    }
}