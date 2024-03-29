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

func Combine(dst *image.RGBA, src *image.Alpha, position Position) (*image.RGBA, error) {
    if dst == nil {
        return nil, NilDstError
    }

    if src == nil {
        return nil, NilSrcError
    }

    var combined *image.RGBA

    zeroPoint := image.Pt(0, 0)

    switch position {
    case PositionAbove:
        combined = image.NewRGBA(image.Rect(0, 0, dst.Bounds().Max.X, dst.Bounds().Max.Y+src.Bounds().Max.Y))
        draw.Draw(combined, src.Bounds(), src, zeroPoint, draw.Src)
        draw.Draw(combined, dst.Bounds().Add(image.Pt(0, src.Bounds().Max.Y)), dst, zeroPoint, draw.Src)
    case PositionBelow:
        combined = image.NewRGBA(image.Rect(0, 0, dst.Bounds().Max.X, dst.Bounds().Max.Y+src.Bounds().Max.Y))
        draw.Draw(combined, dst.Bounds(), dst, zeroPoint, draw.Src)
        draw.Draw(combined, src.Bounds().Add(image.Pt(0, dst.Bounds().Max.Y)), src, zeroPoint, draw.Src)
    case PositionMix:
        combined = dst
        draw.Draw(combined, src.Bounds(), src, zeroPoint, draw.Src)
    }

    return combined, nil
}
