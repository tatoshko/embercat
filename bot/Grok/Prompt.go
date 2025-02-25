package Grok

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "io"
    "log"
    "net/http"
    "net/url"
    "strings"
)

type Request struct {
    Inputs string `json:"inputs"`
}

type Response struct {
    Answer string `json:"answer"`
}

func Prompt(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var logger = getLogger("Prompt")

    text := strings.TrimPrefix(update.Message.Text, "уголек")
    text = strings.TrimPrefix(update.Message.Text, ",")
    text = strings.TrimPrefix(update.Message.Text, " ")

    reqBody := Request{Inputs: text}
    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        logger("Marshal error", err.Error())
        return
    }

    purl := url.URL{}
    url_proxy, _ := purl.Parse("socks5://127.0.0.1:10808")

    transport := http.Transport{}
    transport.Proxy = http.ProxyURL(url_proxy)
    transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    transport.ProxyConnectHeader.Set("Authorization", "Bearer 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")
    client := &http.Client{Transport: &transport}

    url := "https://api-inference.huggingface.co/models/gpt2"
    resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
    if err != nil {
        logger("Request error", err.Error())
        return
    }
    defer resp.Body.Close()

    var result Response
    if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
        logger("Decode error", err.Error())

        if resp.StatusCode == http.StatusOK {
            bodyBytes, err := io.ReadAll(resp.Body)
            if err != nil {
                log.Fatal(err)
            }
            bodyString := string(bodyBytes)
            logger("ACTUAL RESPONSE", bodyString)
        } else {
            logger("ACTUAL Code", fmt.Sprintf("%d", resp.StatusCode))
        }

        return
    }

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, result.Answer)
    msg.ReplyToMessageID = update.Message.MessageID

    if _, err = bot.Send(msg); err != nil {
        logger("Message send error", err.Error())
    }
}
