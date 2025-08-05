package main

func Terrain(newChunk ChunkPointerAndPosition) {
    GenChunk(newChunk.chunk, newChunk.p.x, newChunk.p.y, newChunk.p.z)
    newChunk.chunk.terrainComputed.Store(2)
}

func GenChunk(chunk *Chunk, xPos, yPos, zPos uint32) {
    genIsland(chunk, xPos, yPos, zPos)
}

const WORLD_SIZE = 64
const WORLD_DEPTH = 64
const DIST_BELOW_SPAWN = 30
func genIsland(chunk *Chunk, xPos, yPos, zPos uint32) {
	for i := 0; i < SUBCHUNK_H; i++ {
        for k := 0; k < SUBCHUNK_H; k++ {
            placesoil := true
            for j := int(SUBCHUNK_H-1); j >= 0; j-- {
                xW := xPos*SUBCHUNK_H + uint32(i)
                yW := yPos*SUBCHUNK_H + uint32(j)
                zW := zPos*SUBCHUNK_H + uint32(k)
                xC := Uint32Distance(SPAWN_POSITION.x, xW) < WORLD_SIZE
                yC := ((SPAWN_POSITION.y - yW) > DIST_BELOW_SPAWN) && ((SPAWN_POSITION.y - yW) < (WORLD_DEPTH+DIST_BELOW_SPAWN))
                zC := Uint32Distance(SPAWN_POSITION.z, zW) < WORLD_SIZE
                if xC && yC && zC {
                    if placesoil {
                        placesoil = false
                        chunk.block[i][j][k].id = GRASS
                    } else {
                        chunk.block[i][j][k].id = STONE
                    }
                } else {
                    chunk.block[i][j][k].id = AIR
                }
            }
        }
    }
}
