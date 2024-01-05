package drawer

import (
    "embercat/bot/handlerQuote/service"
    "image"
)

func MakeQuoted(quote *service.Quote, src *image.RGBA, position Position) (err error) {
    var quotedPic *image.RGBA

    if quotedPic, err = MakeQuotePic(quote, src.Bounds()); err != nil {
        return
    }

    return Combine(src, quotedPic, position)
}

func AddQuoteBottom(quote *service.Quote, src *image.RGBA) error {
    return MakeQuoted(quote, src, PositionBelow)
}

func AddQuoteTop(quote *service.Quote, src *image.RGBA) error {
    return MakeQuoted(quote, src, PositionAbove)
}

func MixQuote(quote *service.Quote, src *image.RGBA) error {
    return MakeQuoted(quote, src, PositionMix)
}
