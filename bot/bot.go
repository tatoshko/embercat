package bot

import (
    "embercat/bot/ai"
    "embercat/bot/callbacks"
    "embercat/bot/core"
    "embercat/bot/handlerCatacul"
    "embercat/bot/handlerDeepSeek"
    "embercat/bot/handlerDonate"
    "embercat/bot/handlerPic"
    "embercat/bot/handlerQuote"
    "embercat/bot/handlerThreat"
    "embercat/bot/handlerTurbo"
    "embercat/bot/handlerWednesday"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "math/rand"
    "strings"
)

var (
    Commands = make(map[string]core.Handler)
)

func Start(config Config) {
    if API, err := tba.NewBotAPI(config.Token); err == nil {
        if _, err := API.SetWebhook(tba.NewWebhook(config.Hook + "/" + config.Token)); err != nil {
            log.Printf("SetHoook error %s\n", err.Error())
        }

        API.Debug = false

        registerCommands()

        updates := API.ListenForWebhook("/" + API.Token)

        for update := range updates {
            if update.Message != nil {
                message := update.Message

                direct := int64(message.From.ID) == message.Chat.ID
                tagMe := strings.Index(message.CommandWithAt(), config.Name) != -1

                if message.IsCommand() && (tagMe || direct) {
                    if handler, found := Commands[message.Command()]; found {
                        log.Printf(
                            "MessageID: '%d', Command: '%s', Data: '%s', From: '%d'\n",
                            message.MessageID, message.Command(), message.CommandArguments(), message.From.ID,
                        )
                        go handler(API, update)
                    }
                } else {
                    if handlerWednesday.ItIsWednesdayMyDudes() && rand.Intn(19) == 0 {
                        go handlerWednesday.Check(API, update)
                    }

                    if rand.Intn(49) == 0 {
                        go handlerQuote.Rnd(API, update)
                    }

                    if update.Message.Photo != nil && rand.Intn(20) == 0 {
                        go handlerQuote.Pic(API, update)
                    }

                    if strings.HasPrefix(strings.ToLower(update.Message.Text), "уголек") {
                        go handlerDeepSeek.Prompt(API, update)
                    }

                    go ai.HandlerTextAnalizer(API, update)
                }
            } else if update.CallbackQuery != nil {
                data := update.CallbackQuery.Data

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
    // Other
    Commands["thread"] = handlerThreat.HandleThread

    // Catacul
    handlerCatacul.Init()
    Commands["day"] = handlerCatacul.Day
    Commands["hny"] = handlerCatacul.Hny

    // Turbo
    Commands["turbo"] = handlerTurbo.Roll
    Commands["collection"] = handlerTurbo.MyCollection
    Commands["show"] = handlerTurbo.Show
    Commands["want"] = handlerTurbo.Want
    Commands["doubles"] = handlerTurbo.Double
    callbacks.RegisterCallback("wantans", handlerTurbo.CallbackWant)

    // Wednesday
    Commands["wed"] = handlerWednesday.Check
    Commands["newfrog"] = handlerWednesday.Save

    // Donates
    Commands["donate"] = handlerDonate.Add
    Commands["donates"] = handlerDonate.Show

    // Pictures
    Commands["aserver"] = handlerPic.RndServer
    Commands["ashare"] = handlerPic.RndShare
    Commands["asave"] = handlerPic.Save

    // Quote
    Commands["add"] = handlerQuote.Add
    Commands["pic"] = handlerQuote.Pic
}
