package drawer

import (
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "log"
)

var (
    point         fixed.Point26_6
    srcImg        *image.RGBA
    imageBounds   fixed.Rectangle26_6
    DPI                         = 72.0
    white                       = color.RGBA{R: 255, G: 255, B: 255, A: 255}
    TTF, _                      = truetype.Parse(gobold.TTF)
    StartFontSize fixed.Int26_6 = 42.0
)

func AddLabel(img *image.RGBA, label string) {
    srcImg = img
    ib := img.Bounds()
    imageBounds = fixed.R(ib.Min.X, ib.Min.Y, ib.Max.X, ib.Max.Y)
    drawString(label, StartFontSize)
}

func drawString(label string, size fixed.Int26_6) {
    if size == 0 {
        return
    }

    point = fixed.Point26_6{X: 16, Y: imageBounds.Max.Y - size - 64}
    drawer := &font.Drawer{Dst: srcImg, Src: image.NewUniform(white), Face: getFace(size), Dot: point}

    sb, _ := drawer.BoundString(label)

    if sb.In(imageBounds) {
        log.Printf("OK")
        drawer.DrawString(label)
    } else {
        log.Printf("FAIL SIZE: %d | STRING_BOUND: %v | IGMBOUND: %v", size, sb, imageBounds)
        drawString(label, size-1)
    }
}

func getFace(size fixed.Int26_6) font.Face {
    return truetype.NewFace(TTF, &truetype.Options{Size: float64(size), DPI: DPI, Hinting: font.HintingNone})
}
