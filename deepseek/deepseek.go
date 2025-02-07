package deepseek

import (
    "fmt"
    "github.com/go-deepseek/deepseek"
)

var (
    client deepseek.Client
)

func Init(config Config) {
    var err error
    if client, err = deepseek.NewClient(config.Token); err != nil {
        fmt.Printf("DeepSeek init error: %s", err.Error())
    }
}

func GetClient() deepseek.Client {
    if client == nil {
        panic("DeepSeek: Needed to be inited")
    }

    return client
}
