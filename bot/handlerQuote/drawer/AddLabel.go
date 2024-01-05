package drawer

import (
    "embercat/bot/handlerQuote/service"
    "image"
)

func MakeQuoted(quote *service.Quote, src *image.RGBA, position Position) (quotedPic *image.RGBA, err error) {
    if quotedPic, err = MakeQuotePic(quote, src.Bounds()); err != nil {
        return
    }

    err = Combine(src, quotedPic, position)

    return
}

func AddQuoteBottom(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionBelow)
}

func AddQuoteTop(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionAbove)
}

func MixQuote(quote *service.Quote, src *image.RGBA) (*image.RGBA, error) {
    return MakeQuoted(quote, src, PositionMix)
}
