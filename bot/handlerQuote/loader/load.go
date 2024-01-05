package loader

import (
    "golang.org/x/image/draw"
    "image"
    "net/http"
)

func LoadPicByURL(fileURL string) (m *image.RGBA, err error) {
    var resp *http.Response
    if resp, err = http.Get(fileURL); err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var pic image.Image
    pic, _, err = image.Decode(resp.Body)

    b := pic.Bounds()
    m = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    draw.Draw(m, m.Bounds(), pic, b.Min, draw.Src)

    return
}
