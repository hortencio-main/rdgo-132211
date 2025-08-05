package main

import "math"

type IVec3 struct {
    x uint32
    y uint32
    z uint32
}

type Vec3 struct {
    x float32
    y float32
    z float32
}

func normalize(v Vec3) Vec3 {
	len := float32(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
	return Vec3{v.x / len, v.y / len, v.z / len}
}

func cross(a, b Vec3) Vec3 {
	return Vec3{
		a.y*b.z - a.z*b.y,
		a.z*b.x - a.x*b.z,
		a.x*b.y - a.y*b.x,
	}
}

func dot(a, b Vec3) float32 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}


func AddIVec3( a, b IVec3) IVec3 {
    return IVec3{
        x: a.x+b.x, 
        y: a.y+b.y, 
        z: a.z+b.z,
    }
}

func AddVec3( a, b Vec3) Vec3 {
    return Vec3{
        x: a.x+b.x, 
        y: a.y+b.y, 
        z: a.z+b.z,
    }
}

func Vec3ToIVec3( v Vec3 ) IVec3 {
    return IVec3{uint32(v.x),uint32(v.y),uint32(v.z)}
}
