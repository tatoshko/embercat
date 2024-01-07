package drawer

import (
    "embercat/bot/handlerQuote/service"
    "fmt"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "strings"
)

const (
    inRowCharsCount = 30
)

func MakeQuotePic(quote *service.Quote, srcBounds image.Rectangle, color color.Color) (alpha *image.Alpha, err error) {
    fontSize := int((float64(srcBounds.Bounds().Max.X) * 0.6) / inRowCharsCount)

    rows := makeRows(quote.Words())

    rowsCount := len(rows)
    height := (rowsCount + 1) * fontSize

    r := image.Rect(0, 0, srcBounds.Max.X, height)
    alpha = image.NewAlpha(r)

    ttf, _ := truetype.Parse(gobold.TTF)
    face := truetype.NewFace(ttf, &truetype.Options{Size: float64(fontSize)})
    drawer := font.Drawer{Dst: alpha, Src: image.NewUniform(color), Face: face}

    for i, row := range rows {
        drawer.Dot = fixed.P(0, fontSize*(i+1))
        drawer.DrawString(row)
    }

    drawer.Dot = fixed.P(0, height)
    drawer.DrawString(quote.UserName)

    return
}

func makeRows(words []string) (rows []string) {
    rows = []string{""}

    currentRowIdx := 0
    for _, word := range words {
        lenOfNewStr := len([]rune(rows[currentRowIdx])) + len([]rune(word))

        if lenOfNewStr > inRowCharsCount {
            strings.TrimSpace(rows[currentRowIdx])
            rows = append(rows, "")
            currentRowIdx++
        }

        rows[currentRowIdx] += fmt.Sprintf("%s ", word)
    }

    return
}
