package huggingface

type Config struct {
    Api   string `json:"api"`
    Model string `json:"model"`
    Token string `json:"token"`
    Proxy string `json:"proxy"`
}
