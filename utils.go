// Reserved for things that probably have their place somewhere else

package main

import (
    "math"
    "github.com/go-gl/gl/v2.1/gl"
)

// Returns a pointer to a [16]float32 matrix usable with gl.LoadMatrixf
func LookAt(eye, center Vec3) {
	up := Vec3{0, 1, 0}

	// Forward vector
	z := normalize(Vec3{
		x: eye.x - center.x,
		y: eye.y - center.y,
		z: eye.z - center.z,
	})

	// Right vector
	x := normalize(cross(up, z))

	// True up vector
	y := cross(z, x)

	// Column-major layout for OpenGL
    var view [16]float32
    
	view = [16]float32{
		x.x, y.x, z.x, 0,
		x.y, y.y, z.y, 0,
		x.z, y.z, z.z, 0,
		-dot(x, eye), -dot(y, eye), -dot(z, eye), 1,
	}
    gl.LoadMatrixf(&view[0])
}

func ClampFloat32(f, low, high float32) float32 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func ClampFloat64(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func ScaleVec3(v Vec3, s float32) Vec3 {
    return Vec3{v.x*s,v.y*s,v.z*s}
}

func MaxUint8(a, b uint8) uint8 {
	if a > b {
		return a
	} else {
        return b
    }
}

func MinUint8(values ...uint8) uint8 {
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func Uint32Distance(a, b uint32) uint32 {
    if a > b {
        return a - b
    }
    return b - a
}

func cosf(f float32) float32 {
    return float32(math.Cos(float64(f)))
}
func sinf(f float32) float32 {
    return float32(math.Sin(float64(f)))
}

