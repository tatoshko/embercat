package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

var (
    Commands map[string]Handler
    Callbacks map[int]Handler
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
                        log.Printf("Command: '%s', data: '%s'\n", message.Command(), message.CommandArguments())
                        go handler(API, update)
                    }
                }
            } else if update.CallbackQuery != nil {
                //query := update.CallbackQuery
                //query.Data
            }
        }
    } else {
        log.Fatalf("NewAPIBot error %s\n", err.Error())
    }
}

func registerCommands() {
    Commands = make(map[string]Handler)

    Commands["thread"] = handleThread
    Commands["day"] = handleWednesday
    Commands["yn"] = handlerYesNo
}

func registerCallback(id int, f Handler) {
    Callbacks[id] = f
}

func unregisterCallback(id int) {
    if _, found := Callbacks[id]; found {
        delete(Callbacks, id)
    }
}