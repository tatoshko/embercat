package bot

import (
    "embercat/bot/callbacks"
    "embercat/bot/core"
    "embercat/bot/handlerCatacul"
    "embercat/bot/handlerDonate"
    "embercat/bot/handlerPic"
    "embercat/bot/handlerQuote"
    "embercat/bot/handlerReview"
    "embercat/bot/handlerThreat"
    "embercat/bot/handlerTurbo"
    "embercat/bot/handlerWednesday"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
    "math/rand"
    "net"
    "net/http"
    "strings"
    "time"
)

var (
    Commands = make(map[string]core.Handler)
)

func Start(config Config) {
    log.Printf("Stating new bot...\nName: %s, Token: %s, Hook: %s\n", config.Name, config.Token, config.Hook)

    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout: 10 * time.Second,
        // Force using specific IP for api.telegram.org
    }

    // "Пока" не надо, поменяли сервер, все работает, вроде
    // Подменяем разрешение имени
    //transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
    //    if addr == "api.telegram.org:443" {
    //        // Используем один из известных IP Telegram
    //        // Можно попробовать несколько: 149.154.167.220, 149.154.167.221, 149.154.167.51
    //        addr = "149.154.167.220:443"
    //    }
    //    // Для всех остальных адресов используем стандартный Dialer
    //    return (&net.Dialer{
    //        Timeout:   30 * time.Second,
    //        KeepAlive: 30 * time.Second,
    //    }).DialContext(ctx, network, addr)
    //}

    hc := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }

    var err error
    var API *tba.BotAPI
    if API, err = tba.NewBotAPIWithClient(config.Token, tba.APIEndpoint, hc); err != nil {
        log.Fatalf("NewAPIBot error: %s\n", err.Error())
        return
    }

    log.Printf("Bot API successfully initialized\n")

    log.Printf("Delete old webhook\n")
    dwh := tba.DeleteWebhookConfig{DropPendingUpdates: true}
    if _, err := API.Request(dwh); err != nil {
        log.Printf("DeleteWebhook error %s\n", err.Error())
    }
    log.Printf("Successfully deleted old webhook\n")

    whs := config.Hook + "/" + config.Token
    wh, _ := tba.NewWebhook(whs)

    log.Printf("Trying set new webhook: %s\n", whs)

    if _, err := API.Request(wh); err != nil {
        log.Printf("SetHoook error %s\n", err.Error())
    }

    log.Printf("Successfully setted new webhook\n")

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

                if update.Message.Chat.ID == 300318981 {
                    if update.Message.Photo != nil {
                        go handlerQuote.Pic(API, update)
                    }
                }

                if update.Message.Photo != nil && rand.Intn(20) == 0 {
                    go handlerQuote.Pic(API, update)
                }
            }
        } else if update.CallbackQuery != nil {
            data := update.CallbackQuery.Data

            log.Printf("CallBackQuery %s", data)

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

    // Frog Review
    Commands["reviewstart"] = handlerReview.Start
    Commands["reviewnext"] = handlerReview.Next
    Commands["reviewstop"] = handlerReview.Stop
    callbacks.RegisterCallback(handlerReview.CBFRRemove, handlerReview.CallbackRemove)
    callbacks.RegisterCallback(handlerReview.CBFRStay, handlerReview.CallbackStay)
}
