package handlerQuote

import (
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
)

var (
    TTF   *truetype.Font
    face1 font.Face
)

func Init() {
    TTF, _ = truetype.Parse(gobold.TTF)
    face1 = truetype.NewFace(TTF, &truetype.Options{
        Size:    42.0,
        DPI:     72.0,
        Hinting: font.HintingNone,
    })
}
