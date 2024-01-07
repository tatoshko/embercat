package drawer

import (
    "embercat/bot/handlerQuote/service"
    "image"
)

func MakeQuoted(quote *service.Quote, src *image.RGBA, position Position) (quotedPic *image.RGBA, err error) {
    var alpha *image.Alpha

    if alpha, err = MakeQuotePic(quote, src.Bounds()); err != nil {
        return
    }

    return Combine(src, alpha, position)
}

func AddQuoteBelow(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionBelow)
}

func AddQuoteAbove(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionAbove)
}

func AddQupteMix(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionMix)
}
