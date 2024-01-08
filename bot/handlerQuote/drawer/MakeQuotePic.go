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
    paddingY = 14
    paddingX = 28
)

func MakeQuotePic(quote *service.Quote, srcRect image.Rectangle, color color.Color) (alpha *image.Alpha, err error) {
    if quote.Len() <= 0 {
        return nil, EmptyStringErr
    }

    // Make rows
    ttf, _ := truetype.Parse(gobold.TTF)
    face := truetype.NewFace(ttf, &truetype.Options{Size: float64(fontSize)})
    drawer := font.Drawer{Src: image.NewUniform(color), Face: face}

    words := quote.Words()
    fixedR := fixed.R(0, 0, srcRect.Max.X-paddingX*2, srcRect.Max.Y-paddingY*2)

    currentRow := 0
    rows := []string{words[0]}
    for _, word := range words[1:] {
        newString := fmt.Sprintf("%s %s", rows[currentRow], word)

        _, advice := drawer.BoundString(newString)

        if advice <= fixedR.Max.X {
            rows[currentRow] = newString
        } else {
            rows = append(rows, word)
            currentRow++
        }
    }

    // Print to dst
    height := ((len(rows) + 1) * fontSize) + paddingY*2

    r := image.Rect(0, 0, srcRect.Max.X+paddingX, height+paddingY)
    alpha = image.NewAlpha(r)

    drawer.Dst = alpha

    for i, row := range rows {
        drawer.Dot = fixed.P(paddingX, fontSize*(i+1)+paddingY)
        drawer.DrawString(row)
    }

    ad := drawer.MeasureString(quote.UserName)

    drawer.Dot = fixed.P(srcRect.Max.X-paddingX-ad.Round(), height-paddingY)
    drawer.DrawString(quote.UserName)

    return
}
