package main

import (
	"math"
)

var hash = [256]uint8{
	208, 34, 231, 213, 32, 248, 233, 56, 161, 78, 24, 140, 71, 48, 140, 254, 245, 255, 247, 247, 40,
	185, 248, 251, 245, 28, 124, 204, 204, 76, 36, 1, 107, 28, 234, 163, 202, 224, 245, 128, 167, 204,
	9, 92, 217, 54, 239, 174, 173, 102, 193, 189, 190, 121, 100, 108, 167, 44, 43, 77, 180, 204, 8, 81,
	70, 223, 11, 38, 24, 254, 210, 210, 177, 32, 81, 195, 243, 125, 8, 169, 112, 32, 97, 53, 195, 13,
	203, 9, 47, 104, 125, 117, 114, 124, 165, 203, 181, 235, 193, 206, 70, 180, 174, 0, 167, 181, 41,
	164, 30, 116, 127, 198, 245, 146, 87, 224, 149, 206, 57, 4, 192, 210, 65, 210, 129, 240, 178, 105,
	228, 108, 245, 148, 140, 40, 35, 195, 38, 58, 65, 207, 215, 253, 65, 85, 208, 76, 62, 3, 237, 55, 89,
	232, 50, 217, 64, 244, 157, 199, 121, 252, 90, 17, 212, 203, 149, 152, 140, 187, 234, 177, 73, 174,
	193, 100, 192, 143, 97, 53, 145, 135, 19, 103, 13, 90, 135, 151, 199, 91, 239, 247, 33, 39, 145,
	101, 120, 99, 3, 186, 86, 99, 41, 237, 203, 111, 79, 220, 135, 158, 42, 30, 154, 120, 67, 87, 167,
	135, 176, 183, 191, 253, 115, 184, 21, 233, 58, 129, 233, 142, 39, 128, 211, 118, 137, 139, 255,
	114, 20, 218, 113, 154, 27, 127, 246, 250, 1, 8, 198, 250, 209, 92, 222, 173, 21, 88, 102, 219,
}

func noise2(x, y, seed uint32) uint32 {
	yindex := (y + seed) % 256
	if yindex < 0 {
		yindex += 256
	}
	xindex := uint32(hash[yindex]+uint8(x)) % 256
	if xindex < 0 {
		xindex += 256
	}
	return uint32(hash[xindex])
}

func linInter(x, y, s float64) float64 {
	return x + s*(y-x)
}

func smoothInter(x, y, s float64) float64 {
	return linInter(x, y, s*s*(3-2*s))
}

func noise2d(x, y float64, seed uint32) float64 {
	xInt := uint32(math.Floor(x))
	yInt := uint32(math.Floor(y))
	xFrac := x - float64(xInt)
	yFrac := y - float64(yInt)

	s := float64(noise2(xInt, yInt, seed))
	t := float64(noise2(xInt+1, yInt, seed))
	u := float64(noise2(xInt, yInt+1, seed))
	v := float64(noise2(xInt+1, yInt+1, seed))

	low := smoothInter(s, t, xFrac)
	high := smoothInter(u, v, xFrac)

	return smoothInter(low, high, yFrac)
}

func Perlin2D(x, y, freq float64, depth uint32, seed uint32) float64 {
	xa := x * freq
	ya := y * freq
	amp := 1.0
	fin := 0.0
	div := 0.0

	for i := uint32(0); i < depth; i++ {
		div += 256 * amp
		fin += noise2d(xa, ya, seed) * amp
		amp /= 2
		xa *= 2
		ya *= 2
	}

	return fin / div
}


// slow 3D Perlin noise function
func Perlin3D(x, y, z, freq float64, depth uint32, seed uint32) float64 {
    return ( Perlin2D(x,y,freq,depth,seed) + Perlin2D(y,z,freq,depth,seed) + Perlin2D(z,x,freq,depth,seed) )/3.0
}


// interpolates a value inside a cube
// c000, c001, ..., c111 are the values at the 8 cube corners.
// tx, ty, tz are x, y, z for within the cube. normalized [0..1]
func TrilinearInterpolate(
	c000, c001, c010, c011,
	c100, c101, c110, c111 float64,
	tx, ty, tz float64,
) float64 {
	// Interpolate along x for each of the 4 lower and upper edges
	c00 := c000*(1-tx) + c100*tx
	c01 := c001*(1-tx) + c101*tx
	c10 := c010*(1-tx) + c110*tx
	c11 := c011*(1-tx) + c111*tx

	// Interpolate along y
	c0 := c00*(1-ty) + c10*ty
	c1 := c01*(1-ty) + c11*ty

	// Interpolate along z
	return c0*(1-tz) + c1*tz
}

// fast thread-safe random number generator 🚀
// this file might not be the best place to put it
func LGC(seed uint32) uint32 {
    seed = seed*1664525 + 1013904223
    return seed
}

