package drawer

import (
    "embercat/bot/handlerQuote/service"
    "errors"
    "fmt"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
)

var (
    EmptyStringErr = errors.New("empty string")
)

const (
    fontSize = 28
)

func MakeQuotePic(quote *service.Quote, srcRect image.Rectangle, color color.Color) (alpha *image.Alpha, err error) {
    if quote.Len() <= 0 {
        return nil, EmptyStringErr
    }

    // Make rows
    r := image.Rect(0, 0, srcRect.Max.X, fontSize)
    a := image.NewAlpha(r)

    ttf, _ := truetype.Parse(gobold.TTF)
    face := truetype.NewFace(ttf, &truetype.Options{Size: float64(fontSize)})
    drawer := font.Drawer{Dst: a, Src: image.NewUniform(color), Face: face}

    words := quote.Words()
    fixedR := fixed.R(srcRect.Min.X, srcRect.Min.Y, srcRect.Max.X, srcRect.Max.Y)

    currentRow := 0
    rows := []string{words[0]}
    for _, word := range words[1:] {
        newString := fmt.Sprintf("%s %s", rows[currentRow], word)

        bounds, _ := drawer.BoundString(newString)

        if bounds.In(fixedR) {
            rows[currentRow] = newString
        } else {
            rows = append(rows, word)
            currentRow++
        }
    }

    // Print to dst
    height := (len(rows) + 1) * fontSize

    r = image.Rect(0, 0, srcRect.Max.X, height)
    alpha = image.NewAlpha(r)

    drawer.Dst = alpha

    for i, row := range rows {
        drawer.Dot = fixed.P(0, fontSize*(i+1))
        drawer.DrawString(row)
    }

    drawer.Dot = fixed.P(0, height)
    drawer.DrawString(quote.UserName)

    return
}
