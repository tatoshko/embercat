package bot

import (
    "context"
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
    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout: 10 * time.Second,
        // Force using specific IP for api.telegram.org
    }

    // Подменяем разрешение имени
    transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
        if addr == "api.telegram.org:443" {
            // Используем один из известных IP Telegram
            // Можно попробовать несколько: 149.154.167.220, 149.154.167.221, 149.154.167.51
            addr = "149.154.167.220:443"
        }

        // Для всех остальных адресов используем стандартный Dialer
        return (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext(ctx, network, addr)
    }

    hc := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }

    if API, err := tba.NewBotAPIWithClient(config.Token, tba.APIEndpoint, hc); err == nil {

        wh, _ := tba.NewWebhook(config.Hook + "/" + config.Token)
        if _, err := API.Request(wh); err != nil {
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

    // Frog Review
    Commands["reviewstart"] = handlerReview.Start
    Commands["reviewnext"] = handlerReview.Next
    Commands["reviewstop"] = handlerReview.Stop
    callbacks.RegisterCallback(handlerReview.CBFRRemove, handlerReview.CallbackRemove)
    callbacks.RegisterCallback(handlerReview.CBFRStay, handlerReview.CallbackStay)
}
