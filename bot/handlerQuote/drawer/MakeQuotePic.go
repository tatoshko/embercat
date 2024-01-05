package drawer

import (
    "embercat/bot/handlerQuote/service"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "strings"
)

var (
    defaultFontSize = 42
)

func MakeQuotePic(quote *service.Quote, srcBounds image.Rectangle) (img *image.RGBA, err error) {
    rows := makeRows(quote.Words())

    rowsCount := len(rows)
    height := (rowsCount + 1) * defaultFontSize

    r := image.Rect(0, 0, srcBounds.Max.X, height)

    img = image.NewRGBA(r)

    ttf, _ := truetype.Parse(gobold.TTF)
    face := truetype.NewFace(ttf, &truetype.Options{Size: float64(defaultFontSize)})
    drawer := font.Drawer{Dst: img, Src: image.NewUniform(color.White), Face: face}

    for i, row := range rows {
        drawer.Dot = fixed.P(0, defaultFontSize*(i+1))
        drawer.DrawString(row)
    }

    drawer.Dot = fixed.P(0, height)
    drawer.DrawString(quote.UserName)

    return
}

func makeRows(words []string) (rows []string) {
    inRowCount := 5

    for len(words) > 0 {
        l := len(words)

        if l < inRowCount {
            rows = append(rows, strings.Join(words[:l], " "))
            words = words[l:]
        } else {
            rows = append(rows, strings.Join(words[:inRowCount], " "))
            words = words[inRowCount:]
        }
    }

    return
}
