package irmf

import (
	"fmt"
	"log"
	"strings"
)

func matrixMult(mbb *MBB, vec0, vec1, vec2, vec3 []float64) *MBB {
	mult := func(x, y, z float64) (float64, float64, float64) {
		a := x*vec0[0] + y*vec0[1] + z*vec0[2] + vec0[3]
		b := x*vec1[0] + y*vec1[1] + z*vec1[2] + vec1[3]
		c := x*vec2[0] + y*vec2[1] + z*vec2[2] + vec2[3]
		return a, b, c
	}
	xmin, ymin, zmin := mult(mbb.XMin, mbb.YMin, mbb.ZMin)
	xmax, ymax, zmax := mult(mbb.XMax, mbb.YMax, mbb.ZMax)
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}
	if zmin > zmax {
		zmin, zmax = zmax, zmin
	}
	return &MBB{
		XMin: xmin,
		YMin: ymin,
		ZMin: zmin,
		XMax: xmax,
		YMax: ymax,
		ZMax: zmax,
	}
}

func matrixInverse(vec0, vec1, vec2, vec3 []float64) (inv0, inv1, inv2, inv3 []float64) {
	f := func(v0, v1, v2 []float64) []float64 {
		return []float64{
			v0[1]*v1[2]*v2[3] + v1[1]*v2[2]*v0[3] + v2[1]*v0[2]*v1[3] - v2[1]*v1[2]*v0[3] - v1[1]*v0[2]*v2[3] - v0[1]*v2[2]*v1[3],
			v0[0]*v1[2]*v2[3] + v1[0]*v2[2]*v0[3] + v2[0]*v0[2]*v1[3] - v2[0]*v1[2]*v0[3] - v1[0]*v0[2]*v2[3] - v0[0]*v2[2]*v1[3],
			v0[0]*v1[1]*v2[3] + v1[0]*v2[1]*v0[3] + v2[0]*v0[1]*v1[3] - v2[0]*v1[1]*v0[3] - v1[0]*v0[1]*v2[3] - v0[0]*v2[1]*v1[3],
			v0[0]*v1[1]*v2[2] + v1[0]*v2[1]*v0[2] + v2[0]*v0[1]*v1[2] - v2[0]*v1[1]*v0[2] - v1[0]*v0[1]*v2[2] - v0[0]*v2[1]*v1[2],
		}
	}

	// Matrix of Minors:
	inv0 = f(vec1, vec2, vec3)
	inv1 = f(vec0, vec2, vec3)
	inv2 = f(vec0, vec1, vec3)
	inv3 = f(vec0, vec1, vec2)

	// Matrix of Cofactors:
	inv0 = []float64{inv0[0], -inv0[1], inv0[2], -inv0[3]}
	inv1 = []float64{-inv1[0], inv1[1], -inv1[2], inv1[3]}
	inv2 = []float64{inv2[0], -inv2[1], inv2[2], -inv2[3]}
	inv3 = []float64{-inv3[0], inv3[1], -inv3[2], inv3[3]}

	det := vec0[0]*inv0[0] + vec1[0]*inv1[0] + vec2[0]*inv2[0] + vec3[0]*inv3[0]
	if det == 0 {
		log.Fatalf("Got singular 4x4 matrix with determinant 0:\n%+v\n%+v\n%+v\n%+v", vec0, vec1, vec2, vec3)
	}

	// Matrix of Adjoints:
	inv0[1], inv0[2], inv0[3], inv1[0], inv2[0], inv3[0] = inv1[0], inv2[0], inv3[0], inv0[1], inv0[2], inv0[3]
	inv1[2], inv1[3], inv2[1], inv3[1] = inv2[1], inv3[1], inv1[2], inv1[3]
	inv2[3], inv3[2] = inv3[2], inv2[3]

	// Multiply by one over the determinant:
	det = 1.0 / det
	inv0 = []float64{det * inv0[0], det * inv0[1], det * inv0[2], det * inv0[3]}
	inv1 = []float64{det * inv1[0], det * inv1[1], det * inv1[2], det * inv1[3]}
	inv2 = []float64{det * inv2[0], det * inv2[1], det * inv2[2], det * inv2[3]}
	inv3 = []float64{det * inv3[0], det * inv3[1], det * inv3[2], det * inv3[3]}

	return inv0, inv1, inv2, inv3
}

func vs(vec []float64) string {
	var result []string
	for _, v := range vec {
		s := fmt.Sprintf("%v", v)
		if s == "-0" {
			s = "0"
		}
		result = append(result, s)
	}
	return strings.Join(result, ", ")
}
