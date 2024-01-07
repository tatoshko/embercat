package drawer

import (
    "embercat/bot/handlerQuote/service"
    "image"
    "image/color"
)

func MakeQuoted(quote *service.Quote, src *image.RGBA, position Position, color color.Color) (quotedPic *image.RGBA, err error) {
    var alpha *image.Alpha

    if alpha, err = MakeQuotePic(quote, src.Bounds(), color); err != nil {
        return
    }

    return Combine(src, alpha, position)
}

func AddQuoteBelow(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionBelow, color.White)
}

func AddQuoteAbove(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionAbove, color.White)
}

func AddQupteMix(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionMix, color.Black)
}
