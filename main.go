//~ a recreation of rd-132211 written in go
//~ go build -o rd *.go

package main

import (
    "os"
    "log"
    "runtime"
    "math"
    "math/rand"
    "time"
    "sort"
    "sync/atomic"
    "strconv"
    "net/http"
    _ "net/http/pprof"

    "github.com/go-gl/gl/v2.1/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
)

const SUBCHUNK_H = 16

const SUBCHUNK_V = SUBCHUNK_H  // default size: 16

var RENDER_DISTANCE = uint32(8)

var ABYSS_CENTER_UP   IVec3 = IVec3{15000, 15000, 15000}
var ABYSS_CENTER_DOWN IVec3 = IVec3{15000,  1000, 15000}
var MAX_RADIUS      float32 = float32(1000)

var SPAWN_POSITION IVec3 = IVec3{15000, 15000, 15000}

var Player struct{
    pos Vec3
    vel Vec3
    acc Vec3
    yaw float32
    pitch float32
    dir Vec3
}

var (
    Atlas uint32 
	firstMouse bool  = true
	lastMouseX float64
	lastMouseY float64
	lCamYaw   float32
	lCamPitch float32
 
    StoneLayerSeed uint32
    
    World map[IVec3]*Chunk
)

type Chunk struct {
    block            [SUBCHUNK_H][SUBCHUNK_H][SUBCHUNK_H]Block
    displaylist      uint32
    hasVisibleBlocks bool

    terrainComputed  atomic.Int32
    openToRendering  atomic.Bool

}

func roundToMultiple(x float32, m uint32) uint32 {
    return (uint32(x) / m) * m
}

var strbuf []byte
var newline = []byte{'\n'}
func putInt(num int) {
    strbuf = strconv.AppendInt(strbuf[:0], int64(count), 10)
    os.Stdout.Write(strbuf)
    os.Stdout.Write(newline)
}

var count int
var start time.Time
//~ func printFPS() {
    //~ count++
    //~ if time.Since(start) >= time.Second {
        //~ putInt(count)
        //~ count = 0
        //~ start = time.Now()
    //~ }
//~ }

type ChunkPointerAndPosition struct {
    chunk *Chunk
    p IVec3
}

var WORLD_TICKS uint64


const SeaLevel = 450
const SeaMaxDepth = 300
const CloudLevel = 600



func init() {
    
    if os.Getenv("APP_ENV") == "development" {
        log.Println("Enabling pprof for profiling")
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
    
    StoneLayerSeed = rand.Uint32()
    World = make(map[IVec3]*Chunk)
    runtime.LockOSThread()
    Player.pos = Vec3{
        float32(SPAWN_POSITION.x),
        float32(SPAWN_POSITION.y),
        float32(SPAWN_POSITION.z),
    }
}

func main() {
    
    if err := glfw.Init(); err != nil {
        log.Fatalln("failed to initialize glfw:", err)
    }
    defer glfw.Terminate()

    glfw.WindowHint(glfw.ContextVersionMajor, 2)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)

    window, err := glfw.CreateWindow(800, 600, "3D Cube with Vertex3f", nil, nil)
    if err != nil {
        log.Fatalln("failed to create window:", err)
    }
    window.MakeContextCurrent()

    if err := gl.Init(); err != nil {
        log.Fatalln("failed to initialize OpenGL:", err)
    }
    
    Atlas = LoadTexture("atlas.png")

    window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
    window.SetCursorPosCallback(mouseCursorCallback)
    
    var fogColor [4]float32 = [4]float32{ 0: 0.529, 0.808, 0.922, 1.0 }
    
    
    
    for !window.ShouldClose() {

        
        width, height := window.GetFramebufferSize()
        
        gl.Viewport(0, 0, int32(width), int32(height))
        gl.Enable(gl.DEPTH_TEST)
        gl.MatrixMode(gl.PROJECTION)
        gl.LoadIdentity()

        const fov = 45.0
        const PI = 3.14159
        
        aspect := 800.0 / 600.0;
        near   := 0.1
        far    := 1000.0
        top    := near * math.Tan(fov * PI / 360.0)
        bottom := -top
        right  := top * aspect
        left   := -right
        gl.Frustum(left, right, bottom, top, near, far)
        
        gl.MatrixMode(gl.MODELVIEW)
        
        WORLD_TICKS++
        
        //~ printFPS()
        controls(window)

        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.ClearColor(fogColor[0], fogColor[1], fogColor[2], fogColor[3])
        
        gl.LoadIdentity()

        Player.dir = Vec3{
            x: cosf(lCamYaw) * cosf(lCamPitch),
            y: sinf(lCamPitch),
            z: sinf(lCamYaw) * cosf(lCamPitch),
        }
        
        LookAt( Player.pos, AddVec3(Player.pos,Player.dir) )

        type RenderChunk struct {
            Chunk     *Chunk
            I, J, K      uint32
            Distance2 float32 // squared distance for performance
        }
        var renderQueue[]RenderChunk
        
        var QueueTerrainGeneration []ChunkPointerAndPosition
        
        camPos := Player.pos
        camDir := Player.dir

        xPos := uint32(camPos.x) / SUBCHUNK_H
        yPos := uint32(camPos.y) / SUBCHUNK_H
        zPos := uint32(camPos.z) / SUBCHUNK_H

        for i := (xPos - RENDER_DISTANCE); i < (xPos + RENDER_DISTANCE); i++ {
            for j := (yPos - RENDER_DISTANCE); j < (yPos + RENDER_DISTANCE); j++ {
                for k := (zPos - RENDER_DISTANCE); k < (zPos + RENDER_DISTANCE); k++ {

                    neighbors := []IVec3{
                        {x: i  , y: j  , z: k  },
                        {x: i+1, y: j  , z: k  },
                        {x: i-1, y: j  , z: k  },
                        {x: i  , y: j+1, z: k  },
                        {x: i  , y: j-1, z: k  },
                        {x: i  , y: j  , z: k+1},
                        {x: i  , y: j  , z: k-1},
                    }
                    
                    {
                        _, ok := World[IVec3{i,j,k}]
                        if !ok {
                            newChunk := &Chunk{} // allocate a new chunk
                            World[IVec3{i,j,k}] = newChunk
                            goto cannotRender
                        }
                    }

                    //~ for _, dir := range neighbors {
                    if chunk, _ := World[IVec3{i,j,k}]; chunk.terrainComputed.Load() == 0 {
                    //~ {
                        
                        QueueTerrainGeneration = append( QueueTerrainGeneration, ChunkPointerAndPosition{ chunk:chunk, p:IVec3{i,j,k} } )
                        chunk.terrainComputed.Store(1)
                        goto cannotRender
                    }
                    
                    for _, dir := range neighbors {
                        chunk, ok := World[dir]
                        if !ok {
                            goto cannotRender
                        }
                        if chunk.terrainComputed.Load() != 2 {
                            goto cannotRender
                        }
                    }

                    {
                        chunk, initialized := World[IVec3{i,j,k}]
                        
                        chunk.openToRendering.Store(true)
                        
                        if initialized {
                            // Compute distance from camera
                            centerX := float32(i*SUBCHUNK_H + SUBCHUNK_H/2)
                            centerY := float32(j*SUBCHUNK_H + SUBCHUNK_H/2)
                            centerZ := float32(k*SUBCHUNK_H + SUBCHUNK_H/2)
                            toChunk := Vec3{
                                x: centerX - camPos.x,
                                y: centerY - camPos.y,
                                z: centerZ - camPos.z,
                            }

                            // Check if in front of camera (dot product > 0)
                            dot := toChunk.x*camDir.x + toChunk.y*camDir.y + toChunk.z*camDir.z
                            distSq := toChunk.x*toChunk.x + toChunk.y*toChunk.y + toChunk.z*toChunk.z
                            if (dot > 0.0) || (distSq < (2*SUBCHUNK_H*SUBCHUNK_H) ) {

                                renderQueue = append(renderQueue, RenderChunk{
                                    Chunk:     chunk,
                                    I:         i,
                                    J:         j,
                                    K:         k,
                                    Distance2: distSq,
                                })
                            }
                        }
                    }

                    cannotRender:
                }
            }
        }
        sort.SliceStable(renderQueue, func(a, b int) bool {
            return renderQueue[a].Distance2 > renderQueue[b].Distance2
        })
        
        listsgeneratedthisframe := 0
        for i := len(renderQueue) - 1; i >= 0; i-- {
            if listsgeneratedthisframe < 5 {
            entry := renderQueue[i]
            if !entry.Chunk.openToRendering.Load() { continue }
            if entry.Chunk.displaylist == 0 {
                listsgeneratedthisframe++
                    MakeDisplayList(entry.Chunk, entry.I*SUBCHUNK_H, entry.J*SUBCHUNK_H, entry.K*SUBCHUNK_H, entry.Distance2 )
                }
            }
        }

        for _, entry := range renderQueue {
            if !entry.Chunk.openToRendering.Load() { continue }
            if entry.Chunk.displaylist != 0 {
                if entry.Chunk.hasVisibleBlocks {

                    gl.Enable(gl.FOG);
                    if ((SPAWN_POSITION.y - uint32(entry.J*SUBCHUNK_H)) < 41) {
                        gl.Fogfv(gl.FOG_COLOR, &fogColor[0])
                    } else {
                        black := [4]float32{ 0: 0.0, 0.0, 0.0, 1.0 }
                        gl.Fogfv(gl.FOG_COLOR, &black[0]) 
                        
                    }
                    gl.Fogi(gl.FOG_MODE, gl.LINEAR)
                    gl.Fogf(gl.FOG_START, 0.0 )
                    
                    //~ log.Println( (SPAWN_POSITION.y - uint32(entry.J*SUBCHUNK_H)) , (SPAWN_POSITION.y - uint32(entry.J*SUBCHUNK_H)) < WORLD_SIZE )
                    
                    if ((SPAWN_POSITION.y - uint32(entry.J*SUBCHUNK_H)) < 41) {
                        gl.Fogf(gl.FOG_END, float32((RENDER_DISTANCE-1)*SUBCHUNK_H))
                    } else {
                        gl.Fogf(gl.FOG_END, float32(8))
                    }


                    gl.CallList(entry.Chunk.displaylist)
                    
                    gl.Disable(gl.FOG);
                }
            }
        }

        for _, entry := range QueueTerrainGeneration {
            go Terrain(entry)
        }

        _, wasgenerated := World[IVec3{xPos,yPos,zPos}]
        if wasgenerated {
            Raycast(window)
        }
        
        window.SwapBuffers()
        glfw.PollEvents()
        
        Player.pos = Vec3{ x: Player.pos.x + Player.vel.x, y: Player.pos.y + Player.vel.y, z: Player.pos.z + Player.vel.z }
        Player.vel = Vec3{ x: 0.88*Player.vel.x + Player.acc.x, y: 0.88*Player.vel.y + Player.acc.y, z: 0.88*Player.vel.z + Player.acc.z }
        Player.acc = Vec3{ x: Player.acc.x*0.98, y: Player.acc.y*0.98, z: Player.acc.z*0.98 }
    }
}


var StateKeyTab = 0
var StateEnter = 0
var StateKeyF = 0
func controls(window *glfw.Window) {
    speed := float32(0.03)
    if window.GetKey(glfw.KeyEscape) == glfw.Press {
        window.SetShouldClose(true)
    }
    
    if window.GetKey(glfw.KeyW) == glfw.Press {
        Player.vel.x += cosf(lCamYaw) * speed; Player.vel.z += sinf(lCamYaw) * speed;
    }
    if window.GetKey(glfw.KeyS) == glfw.Press {
        Player.vel.x += -cosf(lCamYaw) * speed; Player.vel.z += -sinf(lCamYaw) * speed;
    }
    if window.GetKey(glfw.KeyA) == glfw.Press {
        Player.vel.x += sinf(lCamYaw) * speed; Player.vel.z += -cosf(lCamYaw) * speed;
    }
    if window.GetKey(glfw.KeyD) == glfw.Press {
        Player.vel.x += -sinf(lCamYaw) * speed; Player.vel.z += cosf(lCamYaw) * speed;
    }
    
    if (window.GetKey(glfw.KeyF) == glfw.Press) && ( StateKeyF == 0) {
        RENDER_DISTANCE++
        RENDER_DISTANCE = 1+(RENDER_DISTANCE%16)
        StateKeyF = 1
    } else if window.GetKey(glfw.KeyF) == glfw.Release {
        StateKeyF = 0
    }
    
    if (window.GetKey(glfw.KeyLeftShift) == glfw.Press)  {
        Player.pos.y -= speed
    } 
    
}

func mouseCursorCallback(window *glfw.Window, xpos, ypos float64) {
	if firstMouse {
		lastMouseX = xpos
		lastMouseY = ypos
		firstMouse = false
	}

	sensitivity := float32(0.0025)
	xoffset := float32(xpos - lastMouseX) * sensitivity
	yoffset := float32(lastMouseY - ypos) * sensitivity // Reversed Y

	lastMouseX = xpos
	lastMouseY = ypos

	lCamYaw += xoffset
	lCamPitch += yoffset

	if lCamPitch > 1.57 {
		lCamPitch = 1.57
	}
	if lCamPitch < -1.57 {
		lCamPitch = -1.57
	}
}

func GetChunk(v IVec3) (*Chunk, bool) {
    chunk, ok := World[IVec3{ x: v.x/SUBCHUNK_H, y: v.y/SUBCHUNK_H, z: v.z/SUBCHUNK_H }]
    return chunk, ok
}

func GetBlock(v IVec3) *Block{
    
    chunk, _ := World[IVec3{
        x: v.x/SUBCHUNK_H,
        y: v.y/SUBCHUNK_H,
        z: v.z/SUBCHUNK_H,
    }]
    
    return &chunk.block[v.x%SUBCHUNK_H][v.y%SUBCHUNK_H][v.z%SUBCHUNK_H]
}
