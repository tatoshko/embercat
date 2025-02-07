package handlerDeepSeek

import (
    "context"
    ds "embercat/deepseek"
    "github.com/go-deepseek/deepseek"
    "github.com/go-deepseek/deepseek/request"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Prompt(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var logger = getLogger("Prompt")
    var client = ds.GetClient()

    chatReq := &request.ChatCompletionsRequest{
        Model:  deepseek.DEEPSEEK_CHAT_MODEL,
        Stream: false,
        Messages: []*request.Message{
            {
                Role:    "user",
                Content: update.Message.Text, // set your input message
            },
        },
    }

    if chatResp, err := client.CallChatCompletionsChat(context.Background(), chatReq); err != nil {
        logger(err.Error())
    } else {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, chatResp.Choices[0].Message.Content)
        msg.ReplyToMessageID = update.Message.MessageID

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }
}
