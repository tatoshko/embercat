package drawer

import (
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
)

var (
    scrText       string
    point         fixed.Point26_6
    srcImg        *image.RGBA
    fakeImg       *image.RGBA
    imageBounds   fixed.Rectangle26_6
    white                 = color.RGBA{R: 255, G: 255, B: 255, A: 255}
    TTF, _                = truetype.Parse(gobold.TTF)
    StartFontSize         = fixed.Int26_6(42)
    DPI           float64 = 72
    padding               = 16
)

func AddLabel(img *image.RGBA, text string) {
    srcImg = img
    ib := img.Bounds()
    imageBounds = fixed.R(ib.Min.X, ib.Min.Y, ib.Max.X, ib.Max.Y)

    scrText = text

    fakeImg = image.NewRGBA(image.Rect(0, 0, ib.Max.X-padding*2, ib.Max.Y))
    draw(computeSize(StartFontSize))
}

func computeSize(size fixed.Int26_6) fixed.Int26_6 {
    if size <= 0 {
        return 0
    }

    point = fixed.Point26_6{X: 0, Y: imageBounds.Max.Y - size*fixed.Int26_6(DPI/4)}
    drawer := &font.Drawer{Dst: fakeImg, Src: image.NewUniform(white), Face: getFace(size), Dot: point}

    sb, _ := drawer.BoundString(scrText)

    if !sb.In(imageBounds) {
        return size
    } else {
        return computeSize(size - 1)
    }
}

func draw(size fixed.Int26_6) {
    if size == 0 {
        return
    }

    point = fixed.Point26_6{X: fixed.Int26_6(padding) * fixed.Int26_6(DPI/4), Y: imageBounds.Max.Y - size*fixed.Int26_6(DPI/4)}
    drawer := &font.Drawer{Dst: srcImg, Src: image.NewUniform(white), Face: getFace(size), Dot: point}
    drawer.DrawString(scrText)
}

func getFace(size fixed.Int26_6) font.Face {
    return truetype.NewFace(TTF, &truetype.Options{Size: float64(size), DPI: DPI})
}
