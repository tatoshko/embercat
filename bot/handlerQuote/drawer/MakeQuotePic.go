package drawer

import (
    "embercat/bot/handlerQuote/service"
    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "log"
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

    drawer := font.Drawer{Dst: img, Src: image.NewUniform(color.RGBA{B: 0xFF, A: 0xFF}), Face: basicfont.Face7x13}

    for i, row := range rows {
        log.Printf("Trying '%s' at %dx%d", row, 0, fixed.Int26_6(defaultFontSize*i*72))

        drawer.Dot = fixed.Point26_6{X: 0, Y: fixed.Int26_6(defaultFontSize * i * 72)}
        drawer.DrawString(row)
    }

    drawer.Dot = fixed.Point26_6{X: 0, Y: fixed.Int26_6(height)}
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
