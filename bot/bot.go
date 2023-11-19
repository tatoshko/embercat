package bot

import (
    "embercat/bot/ai"
    "embercat/bot/callbacks"
    "embercat/bot/handlerDonate"
    "embercat/bot/handlerThreat"
    "embercat/bot/handlerTurbo"
    "embercat/bot/types"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "strings"
)

var (
    Commands = make(map[string]types.Handler)
)

func Start(name, token, hook string) {
    if API, err := tba.NewBotAPI(token); err == nil {
        if _, err := API.SetWebhook(tba.NewWebhook(hook + "/" + token)); err != nil {
            log.Printf("SetHoook error %s\n", err.Error())
        }

        API.Debug = false

        registerCommands()

        updates := API.ListenForWebhook("/" + API.Token)

        for update := range updates {
            if update.Message != nil {
                message := update.Message

                direct := int64(message.From.ID) == message.Chat.ID
                tagMe := strings.Index(message.CommandWithAt(), name) != -1

                if message.IsCommand() && (tagMe || direct) {
                    if handler, found := Commands[message.Command()]; found {
                        log.Printf(
                            "MessageID: '%d', Command: '%s', Data: '%s', From: '%d'\n",
                            message.MessageID, message.Command(), message.CommandArguments(), message.From.ID,
                        )
                        go handler(API, update)
                    }
                } else {
                    go ai.HandlerTextAnalizer(API, update)
                }
            } else if update.CallbackQuery != nil {
                data := update.CallbackQuery.Data

                //log.Printf("CallBackQuery DATA: %v, user: %d, message: %q", data, update.CallbackQuery.From.ID, update.CallbackQuery.Message)

                var handlerID string
                if strings.HasPrefix(data, "/") {
                    parts := strings.SplitN(data, " ", 2)
                    handlerID = strings.TrimPrefix(parts[0], "/")
                } else {
                    handlerID = data
                }

                if handler, found := callbacks.GetHandler(handlerID); found {
                    handler(API, update)
                }
            }
        }
    } else {
        log.Fatalf("NewAPIBot error %s\n", err.Error())
    }
}

func registerCommands() {
    initHandleTime()

    Commands["thread"] = handlerThreat.HandleThread
    Commands["day"] = handlerDay
    Commands["wed"] = handleWednesday
    Commands["btn"] = handlerBtn
    Commands["time"] = handlerTime

    Commands["turbo"] = handlerTurbo.HandlerTurbo
    Commands["collection"] = handlerTurbo.HandlerCollection
    Commands["show"] = handlerTurbo.HandlerShow
    Commands["want"] = handlerTurbo.HandlerWant

    Commands["donate"] = handlerDonate.Add
    Commands["donates"] = handlerDonate.Show

    callbacks.RegisterCallback("wantans", handlerTurbo.CallbackWant)
}
