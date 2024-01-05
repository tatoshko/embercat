package drawer

import (
    "errors"
    "golang.org/x/image/draw"
    "image"
)

var (
    NilDstError = errors.New("dst param is nil")
    NilSrcError = errors.New("src param is nil")
)

type Position int

const (
    PositionAbove Position = iota
    PositionBelow
    PositionMix
)

func Combine(dst *image.RGBA, src *image.RGBA, position Position) error {
    if dst == nil {
        return NilDstError
    }

    if src == nil {
        return NilSrcError
    }

    switch position {
    case PositionAbove:
        combined := image.NewRGBA(image.Rect(0, 0, dst.Bounds().Max.X, dst.Bounds().Max.Y+src.Bounds().Max.Y))
        draw.Draw(combined, src.Bounds(), src, image.Point{X: 0, Y: 0}, draw.Src)
        draw.Draw(combined, dst.Bounds(), dst, image.Point{X: 0, Y: src.Bounds().Max.Y}, draw.Src)
        dst = combined
    case PositionBelow:
        combined := image.NewRGBA(image.Rect(0, 0, dst.Bounds().Max.X, dst.Bounds().Max.Y+src.Bounds().Max.Y))
        draw.Draw(combined, dst.Bounds(), dst, image.Point{X: 0, Y: 0}, draw.Src)
        draw.Draw(combined, src.Bounds(), src, image.Point{X: 0, Y: dst.Bounds().Max.Y}, draw.Src)
        dst = combined
    case PositionMix:
        draw.Draw(dst, dst.Bounds(), src, dst.Bounds().Min, draw.Src)
    }

    return nil
}
