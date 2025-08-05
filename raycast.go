package main


import (
    "math"
    "github.com/go-gl/gl/v2.1/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
)

func ApplyOverNeighborhood( v IVec3, fn func(v IVec3) ) {
    var Neighborhood = []IVec3{
        0: IVec3{x:v.x+1, y:v.y  , z:v.z  },
        1: IVec3{x:v.x-1, y:v.y  , z:v.z  },
        2: IVec3{x:v.x  , y:v.y+1, z:v.z  },
        3: IVec3{x:v.x  , y:v.y-1, z:v.z  },
        4: IVec3{x:v.x  , y:v.y  , z:v.z+1},
        5: IVec3{x:v.x  , y:v.y  , z:v.z-1},
    }
    for _, element := range Neighborhood {
        fn(element)
    }
}

// state machines for the mouse 
// TODO: make their states into enums rather than 1/0
// TODO2: do thing above
var stateMouseButtonLeft int
var stateMouseButtonRight int

var stateBlockBreaking int

func Raycast(window *glfw.Window) {
    
    camDirection := Vec3{
        x: cosf(lCamYaw) * cosf(lCamPitch),
        y: sinf(lCamPitch),
        z: sinf(lCamYaw) * cosf(lCamPitch),
    }
    
	directions := []Vec3{
		0:{-0.5, 0.0, 0.5}, 1:{ 0.0, 0.0, 0.5}, 2:{ 0.5, 0.0, 0.5}, //players head collision
		3:{-0.5, 0.0, 0.0},                     4:{ 0.5, 0.0, 0.0},
        5:{-0.5, 0.0,-0.5}, 6:{ 0.0, 0.0,-0.5}, 7:{ 0.5, 0.0,-0.5},
        
		8:{-0.5,-0.5, 0.5}, 9:{ 0.0,-0.5, 0.5},10:{ 0.5,-0.5, 0.5}, // playres leg collision
       11:{-0.5,-0.5, 0.0},                    12:{ 0.5,-0.5, 0.0},
       13:{-0.5,-0.5,-0.5},14:{ 0.0,-0.5,-0.5},15:{ 0.5,-0.5,-0.5},

	   16:{ 0.0, 0.5, 0.0}, // vertical collision above
	   17:{ 0.0,-0.5, 0.0}, // vertical collision down
       
	   18:{ 0.0, Player.vel.y, 0.0},
       
       19:ScaleVec3(camDirection,1.0),
	}

	collisions := make([]bool, len(directions))
    
    var highlightpos Vec3
    var previousBlock Vec3
    var underwater bool
    
	for i, dir := range directions {
		length := float32(math.Sqrt(float64(dir.x*dir.x + dir.y*dir.y + dir.z*dir.z)))
        
		dx := dir.x / length
		dy := dir.y / length
		dz := dir.z / length

		x := int(math.Floor(float64(Player.pos.x)))
		y := int(math.Floor(float64(Player.pos.y)))
		z := int(math.Floor(float64(Player.pos.z)))

		stepX := 1
		if dx <= 0 { stepX = -1 }
		stepY := 1
		if dy <= 0 { stepY = -1 }
		stepZ := 1
		if dz <= 0 { stepZ = -1 }

		tMaxX := float32(math.Inf(1))
		if dx != 0 {
			if stepX > 0 {
				tMaxX = (float32(x+1) - Player.pos.x) / float32(math.Abs(float64(dx)))
			} else {
				tMaxX = (Player.pos.x - float32(x)) / float32(math.Abs(float64(dx)))
			}
		}
		tMaxY := float32(math.Inf(1))
		if dy != 0 {
			if stepY > 0 {
				tMaxY = (float32(y+1) - Player.pos.y) / float32(math.Abs(float64(dy)))
			} else {
				tMaxY = (Player.pos.y - float32(y)) / float32(math.Abs(float64(dy)))
			}
		}
		tMaxZ := float32(math.Inf(1))
		if dz != 0 {
			if stepZ > 0 {
				tMaxZ = (float32(z+1) - Player.pos.z) / float32(math.Abs(float64(dz)))
			} else {
				tMaxZ = (Player.pos.z - float32(z)) / float32(math.Abs(float64(dz)))
			}
		}

		tDeltaX := float32(math.Inf(1))
		if dx != 0 {
			tDeltaX = 1.0 / float32(math.Abs(float64(dx)))
		}
		tDeltaY := float32(math.Inf(1))
		if dy != 0 {
			tDeltaY = 1.0 / float32(math.Abs(float64(dy)))
		}
		tDeltaZ := float32(math.Inf(1))
		if dz != 0 {
			tDeltaZ = 1.0 / float32(math.Abs(float64(dz)))
		}

		distance := float32(0)
        var maxDistance float32
        if i == 19 {
            maxDistance = 5.0
        } else if i == 17 {
            maxDistance = 1.5
        } else {
            maxDistance = 0.25
        }
        

		for distance < maxDistance {
        
			if tMaxX < tMaxY {
				if tMaxX < tMaxZ {
					x += stepX
					distance = tMaxX
					tMaxX += tDeltaX
				} else {
					z += stepZ
					distance = tMaxZ
					tMaxZ += tDeltaZ
				}
			} else {
				if tMaxY < tMaxZ {
					y += stepY
					distance = tMaxY
					tMaxY += tDeltaY
				} else {
					z += stepZ
					distance = tMaxZ
					tMaxZ += tDeltaZ
				}
			}

            // if our ray hits an block
            
            var transparent = [N_BLOCK_IDS]bool{
                AIR: true,
                WATER: true,
            }
            
			if b := GetBlock(IVec3{x:uint32(math.Round(float64(x))),y:uint32(math.Round(float64(y))),z:uint32(math.Round(float64(z)))}); !transparent[b.id] {
				collisions[i] = true

                if i == 19 {
                    highlightpos = Vec3{x:float32(math.Round(float64(x))),y:float32(math.Round(float64(y))),z:float32(math.Round(float64(z)))}
                }
                break

			} else {
                if i == 17 {
                    if b.id == WATER {
                        underwater = true
                    }
                }
                if i == 19 {
                    previousBlock = Vec3{x:float32(math.Round(float64(x))),y:float32(math.Round(float64(y))),z:float32(math.Round(float64(z)))}
                }
            }
            
		}
	}


    if collisions[19] {
        if (window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press) && (stateMouseButtonLeft == 0) {

            blockpos := Vec3ToIVec3(previousBlock)
            _, ok := World[IVec3{blockpos.x/SUBCHUNK_H, blockpos.y/SUBCHUNK_H, blockpos.z/SUBCHUNK_H}]
            if ok {
                
                GetBlock(blockpos).id = COBBLE
                ApplyOverNeighborhood( blockpos, func(pos IVec3){
                    chunkp, ok := GetChunk(pos)
                    if !ok { return }
                    gl.DeleteLists( chunkp.displaylist , 1)
                    chunkp.displaylist = 0
                }) 
            }
            stateMouseButtonLeft = 1
        }

        if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press) && (stateMouseButtonRight == 0) {

            //~ blockpos := Vec3ToIVec3(highlightpos)
            blockpos := Vec3ToIVec3(highlightpos)
            _, ok := World[IVec3{blockpos.x/SUBCHUNK_H, blockpos.y/SUBCHUNK_H, blockpos.z/SUBCHUNK_H}]
            if ok {
                GetBlock(blockpos).id = AIR
                
                ApplyOverNeighborhood( blockpos, func(pos IVec3){
                    chunkp, ok := GetChunk(pos)
                    if !ok { return }
                    gl.DeleteLists( chunkp.displaylist , 1)
                    chunkp.displaylist = 0
                }) 
            }
            stateMouseButtonRight = 1
        }
        
        drawTransparentCube(RenderFunctionArgv{
            x:highlightpos.x,
            y:highlightpos.y,
            z:highlightpos.z,
        })
    }
    if (window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Release) {
        stateMouseButtonLeft = 0
    }
    if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Release) {
        stateMouseButtonRight = 0
    }

	//~ if collisions[2] || collisions[4] || collisions[7] || collisions[10] || collisions[12] || collisions[15] {
	if collisions[4] || collisions[12]  {
		if Player.vel.x > 0.0 {
			Player.vel.x = 0
		}
	}
	//~ if collisions[0] || collisions[3] || collisions[5] || collisions[8] || collisions[11] || collisions[13] {
	if collisions[3] || collisions[11] {
		if Player.vel.x < 0.0 {
			Player.vel.x = 0
		}
	}
	//~ if collisions[0] || collisions[1] || collisions[2] || collisions[8]  || collisions[9]  || collisions[10] {
	if  collisions[1] || collisions[9]  {
		if Player.vel.z > 0.0 {
			Player.vel.z = 0
		}
	}
	//~ if collisions[5] || collisions[6] || collisions[7] || collisions[13] || collisions[14] || collisions[15] {
	if  collisions[6] || collisions[14] {
		if Player.vel.z < 0.0 {
			Player.vel.z = 0
		}
	}
	if collisions[16] {
		if Player.vel.y > 0.0 {
			Player.vel.y = 0
		}
	}
	if !collisions[17] {
        if underwater {
            Player.vel.y -= 0.075 / 4.0
        } else {
            Player.vel.y -= 0.075
        }
	} else {
        Player.acc.y = 0
        if window.GetKey(glfw.KeySpace) == glfw.Press {
            Player.acc.y = 0.10
        }
    }
	if collisions[18] {
		if Player.vel.y < 0.0 {
			Player.vel.y = 0
		}
	}
}

