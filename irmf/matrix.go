package irmf

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
