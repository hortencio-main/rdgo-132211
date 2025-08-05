package main

import (
    "image"
    "image/draw"
    "image/png"
    //~ "os"
    "log"
    "github.com/go-gl/gl/v2.1/gl"
    "bytes"
    
    _ "embed"
)

//go:embed atlas.png
var atlasData []byte

//~ func LoadTexture(filename string) uint32 {
    //~ file, err := os.Open(filename)
    
func LoadTexture() uint32 {
    //~ file := atlas
    
    img, err := png.Decode(bytes.NewReader(atlasData))
    if err != nil {
        log.Println("INFO: Failed to decode embedded texture: atlas.png")
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

    log.Println("INFO: atlas Texture loaded successfully")
    
    return texture
}
