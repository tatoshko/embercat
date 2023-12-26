package handlerQuote

import (
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font/gofont/gobold"
)

var (
    TTF *truetype.Font
)

func Init() {
    TTF, _ = truetype.Parse(gobold.TTF)
}
