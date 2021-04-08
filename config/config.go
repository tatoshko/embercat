package config

type Config struct {
    Token string `json:"token"`
    Hook  string `json:"hook"`
    DB    string `json:"db"`
}
