package main

import (
    "image"
    "image/draw"
    "image/png"
    "os"
    "log"
    "github.com/go-gl/gl/v2.1/gl"
)

// go allocates memory for these strings even when they are not written
// todo change for a non-escaping way of printing
func LoadTexture(filename string) uint32 {
    file, err := os.Open(filename)
    if err != nil {
        log.Println("INFO: Failed to open texture file:", filename)
        return 0
    }
    defer file.Close()

    img, err := png.Decode(file)
    if err != nil {
        log.Println("INFO: Failed to decode texture:", filename)
        return 0
    }

    rgba := image.NewRGBA(img.Bounds())
    draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

    width := int32(rgba.Bounds().Dx())
    height := int32(rgba.Bounds().Dy())

    var texture uint32
    gl.GenTextures(1, &texture)
    gl.BindTexture(gl.TEXTURE_2D, texture)

    gl.TexImage2D(
        gl.TEXTURE_2D,
        0,
        gl.RGBA,
        width,
        height,
        0,
        gl.RGBA,
        gl.UNSIGNED_BYTE,
        gl.Ptr(rgba.Pix),
    )

    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

    log.Println("INFO: Texture loaded successfully", filename, texture)
    
    return texture
}
