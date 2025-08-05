package main

import "github.com/go-gl/gl/v2.1/gl"

type Color struct { r, g, b float32; }

const textureSize = 16.0

type RenderFunctionArgv struct {
    id          uint32
    x, y, z     float32
    visible     uint32
}

func MakeDisplayList(chunk *Chunk, xLocation, yLocation, zLocation uint32, distance float32) {    
    
    list := gl.GenLists(1)
    gl.NewList(list, gl.COMPILE);

    gl.Enable(gl.TEXTURE_2D);
    gl.BindTexture(gl.TEXTURE_2D, Atlas)

    gl.Enable(gl.BLEND);
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);

    gl.Begin(gl.QUADS)
    
    chunk.hasVisibleBlocks = false
    
    const (
        FIRST_PASS = 0
        SECOND_PASS = 1
        N_PASSES = 2
    )
    
    var skipDrawingOn = [2][N_BLOCK_IDS]bool{
        FIRST_PASS: {
            WATER: true,
            AIR: true,
        },
        SECOND_PASS: {
            AIR: true,
        }, 
    }
    
    var drawIfNextTo = [2][N_BLOCK_IDS]bool{
        FIRST_PASS: {
            WATER: true,
        },
        SECOND_PASS: {
            AIR: true,
            TORCH: true,
        },
    }
    
    for pass := 0; pass < N_PASSES; pass++ {
        for x := 0; x < SUBCHUNK_H; x++ {
            for y := 0; y < SUBCHUNK_H; y++ {
                for z := 0; z < SUBCHUNK_H; z++ {
                    
                    if !skipDrawingOn[pass][ chunk.block[x][y][z].id ] {
                        
                        base := IVec3{
                            x: uint32(x) + xLocation,
                            y: uint32(y) + yLocation,
                            z: uint32(z) + zLocation,
                        }
                        
                        var BlockNeighborhood = []IVec3{
                            0: IVec3{x:base.x+1, y:base.y  , z:base.z  },
                            1: IVec3{x:base.x-1, y:base.y  , z:base.z  },
                            2: IVec3{x:base.x  , y:base.y+1, z:base.z  },
                            3: IVec3{x:base.x  , y:base.y-1, z:base.z  },
                            4: IVec3{x:base.x  , y:base.y  , z:base.z+1},
                            5: IVec3{x:base.x  , y:base.y  , z:base.z-1},
                        }

                        visibleFaces := uint32(0)
                        
                        for i, v := range(BlockNeighborhood){
                            block := GetBlock(v)
                            if drawIfNextTo[pass][block.id] {
                                visibleFaces += 1 << uint32(i)
                            }
                        }

                        if visibleFaces != 0 {
                            chunk.hasVisibleBlocks = true
                                
                            drawCube(RenderFunctionArgv{
                                id: uint32(chunk.block[x][y][z].id),
                                x: float32(x) + float32(xLocation),
                                y: float32(y) + float32(yLocation),
                                z: float32(z) + float32(zLocation),
                                visible: visibleFaces,
                            })
                        }
                    }
                }
            }
        }
    }
    gl.End()
    gl.Disable(gl.TEXTURE_2D);
    gl.EndList()
    chunk.displaylist = list
}

func drawCube( arg RenderFunctionArgv) {
    id := arg.id
    worldX := arg.x
    worldY := arg.y
    worldZ := arg.z
    visibleFaces := arg.visible
    var v = [6][4][3]float32{
        {{1, 0, 1}, {1, 0, 0}, {1, 1, 0}, {1, 1, 1}}, // +X (bit 0)
        {{0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}}, // -X (bit 1)
        {{0, 1, 0}, {0, 1, 1}, {1, 1, 1}, {1, 1, 0}}, // +Y (bit 2)
        {{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}}, // -Y (bit 3)
        {{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 1, 1}}, // +Z (bit 4)
        {{1, 0, 0}, {0, 0, 0}, {0, 1, 0}, {1, 1, 0}}, // -Z (bit 5)
    }
    for f := uint32(0); f < 6; f++ {
        mask := uint32(1) << f
        if visibleFaces & mask == 0 {
            continue // face invisible
        }
        v0 := textVPos(id)
        v1 := textVPos(id+1)
        h0 := texHPos(  f, 0)
        h1 := texHPos(f+1, 0)
        gl.TexCoord2f(h1, v1);
        gl.Vertex3f(worldX+v[f][0][0], worldY+v[f][0][1], worldZ+v[f][0][2])
        gl.TexCoord2f(h0, v1);
        gl.Vertex3f(worldX+v[f][1][0], worldY+v[f][1][1], worldZ+v[f][1][2])
        gl.TexCoord2f(h0, v0);
        gl.Vertex3f(worldX+v[f][2][0], worldY+v[f][2][1], worldZ+v[f][2][2])
        gl.TexCoord2f(h1, v0);
        gl.Vertex3f(worldX+v[f][3][0], worldY+v[f][3][1], worldZ+v[f][3][2])
    }
}

func drawTransparentCube( argv RenderFunctionArgv ) {
	const (
		sx    = 1.1
		ex    = 1.0 - sx
		alpha = 0.5
	)
	var v = [6][4][3]float32{
		{{sx, ex, sx}, {sx, ex, ex}, {sx, sx, ex}, {sx, sx, sx}},
		{{ex, ex, ex}, {ex, ex, sx}, {ex, sx, sx}, {ex, sx, ex}},
		{{ex, ex, sx}, {sx, ex, sx}, {sx, sx, sx}, {ex, sx, sx}},
		{{sx, ex, ex}, {ex, ex, ex}, {ex, sx, ex}, {sx, sx, ex}},
		{{ex, sx, ex}, {ex, sx, sx}, {sx, sx, sx}, {sx, sx, ex}},
		{{ex, ex, ex}, {sx, ex, ex}, {sx, ex, sx}, {ex, ex, sx}},
	}
	gl.PushMatrix()
	gl.Begin(gl.QUADS)
	gl.Color4f(1.0, 1.0, 1.0, alpha)
	for f := 0; f < 6; f++ {
		for vtx := 0; vtx < 4; vtx++ {
			gl.Vertex3f(argv.x+v[f][vtx][0],argv.y+v[f][vtx][1],argv.z+v[f][vtx][2])
		}
	}
	gl.Color4f(1.0, 1.0, 1.0, 1.0)
	gl.End()
	gl.PopMatrix()
}    

func texHPos( face, light_level uint32) float32 {
	return float32(face)*(16.0/1024.0)
}

func textVPos( id uint32 ) float32 {
    return float32(id)*(16.0/1024.0)
}

