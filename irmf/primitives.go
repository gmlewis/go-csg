package irmf

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gmlewis/go-csg/ast"
)

var primitives = map[string]string{
	"circle": `float circle(in float radius, in vec3 xyz) {
	float r = length(xyz.xy);
	return r <= radius ? 1.0 : 0.0;
}
`,

	"cube": `float cube(in vec3 size, in bool center, in vec3 xyz) {
	xyz /= size;
	if (!center) { xyz -= vec3(0.5); }
	if (any(greaterThan(abs(xyz), vec3(0.5)))) { return 0.0; }
	return 1.0;
}
`,

	"cylinder": `float cylinder(in float h, in float r1, in float r2, in bool center, in vec3 xyz) {
	xyz.z /= h;
	float z = xyz.z;
	if (center) { z += 0.5; } else { xyz.z -= 0.5; }
	if (abs(xyz.z) > 0.5) { return 0.0; }
	float r = length(xyz.xy);
	float radius = mix(r1, r2, z);
	return r <= radius ? 1.0 : 0.0;
}
`,

	"polygon": `float polygon(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,

	"polyhedron": `float polyhedron(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,

	"rotAxis": `mat3 rotAxis(vec3 axis, float a) {
  // This is from: http://www.neilmendoza.com/glsl-rotation-about-an-arbitrary-axis/
  float s = sin(a);
  float c = cos(a);
  float oc = 1.0 - c;
  vec3 as = axis * s;
  mat3 p = mat3(axis.x * axis, axis.y * axis, axis.z * axis);
  mat3 q = mat3(c, - as.z, as.y, as.z, c, - as.x, - as.y, as.x, c);
  return p * oc + q;
}
`,

	"rotZ": `mat4 rotZ(float angle) {
  return mat4(rotAxis(vec3(0, 0, 1), angle));
}
`,

	"sphere": `float sphere(in float radius, in vec3 xyz) {
	float r = length(xyz);
	return r <= radius ? 1.0 : 0.0;
}
`,

	"square": `float square(in vec2 size, in bool center, in vec3 xyz) {
	xyz.xy /= size;
	if (!center) { xyz.xy -= vec2(0.5); }
	if (any(greaterThan(abs(xyz.xy), vec2(0.5)))) { return 0.0; }
	return 1.0;
}
`,
}

func (s *Shader) getArgs(exps []ast.Expression, names ...string) []string {
	result := make([]string, len(names))
	values := map[string]string{}
	var count int
	for _, exp := range exps {
		switch exp := exp.(type) {
		case *ast.NamedArgument:
			values[exp.Name.String()] = exp.Value.String()
		case *ast.StringLiteral, *ast.IntegerLiteral, *ast.FloatLiteral, *ast.BooleanLiteral, *ast.ArrayLiteral:
			result[count] = exp.String()
			count++
		default:
			log.Fatalf("getArgs: unhandled type %T (%+v)", exp, exp)
		}
	}

	for i, name := range names {
		if v, ok := values[name]; ok {
			result[i] = v
		}
	}
	return result
}

func (s *Shader) getMat4Args(exps []ast.Expression) []string {
	var result []string
	for _, exp := range exps {
		switch exp := exp.(type) {
		case *ast.ArrayLiteral:
			// Split up array into 4 expressions (each arrays).
			if len(exp.Elements) != 4 {
				log.Fatalf("getMat4Args: unhandled elements!=4 %T (%+v)", exp, exp)
			}
			for _, e := range exp.Elements {
				switch e := e.(type) {
				case *ast.ArrayLiteral:
					val := strings.Trim(e.String(), "[]")
					val = strings.ReplaceAll(val, "(", "")
					val = strings.ReplaceAll(val, ")", "")
					result = append(result, val)
				default:
					log.Fatalf("getMat4Args: unhandled element type %T (%+v)", e, e)
				}
			}
		default:
			log.Fatalf("getMat4Args: unhandled type %T (%+v)", exp, exp)
		}
	}
	return result
}

func parseVec4(s string) ([]float64, error) {
	if !strings.Contains(s, ",") {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q", s)
		}
		return []float64{v, v, v, v}, nil
	}

	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("error parsing %q", s)
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing x: %q: %v", s, err)
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing y: %q: %v", s, err)
	}

	z, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing z: %q: %v", s, err)
	}

	w, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing w: %q: %v", s, err)
	}

	return []float64{x, y, z, w}, nil
}

func parseVec3(s string) ([]float64, error) {
	if !strings.Contains(s, ",") {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q", s)
		}
		return []float64{v, v, v}, nil
	}

	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("error parsing %q", s)
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing x: %q: %v", s, err)
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing y: %q: %v", s, err)
	}

	z, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing z: %q: %v", s, err)
	}

	return []float64{x, y, z}, nil
}

func parseVec2(s string) ([]float64, error) {
	if !strings.Contains(s, ",") {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q", s)
		}
		return []float64{v, v}, nil
	}

	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("error parsing %q", s)
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing x: %q: %v", s, err)
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing y: %q: %v", s, err)
	}

	return []float64{x, y}, nil
}

func (s *Shader) processCubePrimitive(exps []ast.Expression) (string, *MBB) {
	s.Primitives["cube"] = true
	args := s.getArgs(exps, "size", "center")

	size := strings.Trim(args[0], "[]")
	if size == "" {
		size = "1"
	}

	vec3, err := parseVec3(size)
	if err != nil {
		log.Printf("error parsing cube size=%q, setting to 1", size)
		size = "1"
		vec3 = []float64{1, 1, 1}
	}

	var mbb *MBB

	center := args[1]
	if center == "true" {
		mbb = &MBB{XMin: -0.5 * vec3[0], YMin: -0.5 * vec3[1], ZMin: -0.5 * vec3[2], XMax: 0.5 * vec3[0], YMax: 0.5 * vec3[1], ZMax: 0.5 * vec3[2]}
	} else {
		center = "false"
		mbb = &MBB{XMax: vec3[0], YMax: vec3[1], ZMax: vec3[2]}
	}

	return fmt.Sprintf("cube(vec3(%v), %v, xyz)", size, center), mbb
}

func (s *Shader) processSpherePrimitive(exps []ast.Expression) (string, *MBB) {
	s.Primitives["sphere"] = true
	args := s.getArgs(exps, "r", "d")

	radius := args[0]
	diameter := args[1]
	if diameter != "" && radius == "" {
		if vec3, err := parseVec3(diameter); err == nil {
			radius = fmt.Sprintf("%v", 0.5*vec3[0])
		}
	}
	if radius == "" {
		radius = "1"
	}

	vec3, err := parseVec3(radius)
	if err != nil {
		log.Printf("error parsing sphere radius=%q, setting to 1", radius)
		radius = "1"
		vec3 = []float64{1, 1, 1}
	}

	mbb := &MBB{XMin: -vec3[0], YMin: -vec3[1], ZMin: -vec3[2], XMax: vec3[0], YMax: vec3[1], ZMax: vec3[2]}

	return fmt.Sprintf("sphere(float(%v), xyz)", radius), mbb
}

func (s *Shader) processCylinderPrimitive(exps []ast.Expression) (string, *MBB) {
	s.Primitives["cylinder"] = true
	args := s.getArgs(exps, "h", "r1", "r2", "center", "r", "d", "d1", "d2")

	h := args[0]
	r1 := args[1]
	r2 := args[2]
	center := args[3]
	r := args[4]
	d := args[5]
	d1 := args[6]
	d2 := args[7]

	if d2 != "" && r2 == "" {
		if vec3, err := parseVec3(d2); err == nil {
			r2 = fmt.Sprintf("%v", 0.5*vec3[0])
		}
	}
	if d1 != "" && r1 == "" {
		if vec3, err := parseVec3(d1); err == nil {
			r1 = fmt.Sprintf("%v", 0.5*vec3[0])
		}
	}
	if d != "" && r1 == "" && r2 == "" {
		if vec3, err := parseVec3(d); err == nil {
			r1 = fmt.Sprintf("%v", 0.5*vec3[0])
			r2 = r1
		}
	}
	if r != "" && r1 == "" && r2 == "" {
		if vec3, err := parseVec3(r); err == nil {
			r1 = fmt.Sprintf("%v", vec3[0])
			r2 = r1
		}
	}

	if h == "" {
		h = "1"
	}
	if r1 == "" {
		r1 = "1"
	}
	if r2 == "" {
		r2 = "1"
	}

	params := fmt.Sprintf("%v,%v,%v", h, r1, r2)
	vec3, err := parseVec3(params)
	if err != nil {
		log.Printf("error parsing cylinder params %q, setting to 1", params)
		vec3 = []float64{1, 1, 1}
	}

	radius := vec3[1]
	if vec3[2] > radius {
		radius = vec3[2]
	}

	mbb := &MBB{XMin: -radius, YMin: -radius, ZMin: -0.5 * vec3[0], XMax: radius, YMax: radius, ZMax: 0.5 * vec3[0]}
	if center != "true" {
		center = "false"
		mbb.ZMin = 0
		mbb.ZMax = vec3[0]
	}

	return fmt.Sprintf("cylinder(float(%v), float(%v), float(%v), %v, xyz)", h, r1, r2, center), mbb
}

func (s *Shader) processSquarePrimitive(exps []ast.Expression) (string, *MBB) {
	s.Primitives["square"] = true
	args := s.getArgs(exps, "size", "center")

	size := strings.Trim(args[0], "[]")
	if size == "" {
		size = "1"
	}

	vec2, err := parseVec2(size)
	if err != nil {
		log.Printf("error parsing square size=%q, setting to 1", size)
		size = "1"
		vec2 = []float64{1, 1}
	}

	var mbb *MBB

	center := args[1]
	if center == "true" {
		mbb = &MBB{XMin: -0.5 * vec2[0], YMin: -0.5 * vec2[1], XMax: 0.5 * vec2[0], YMax: 0.5 * vec2[1]}
	} else {
		center = "false"
		mbb = &MBB{XMax: vec2[0], YMax: vec2[1]}
	}

	return fmt.Sprintf("square(vec2(%v), %v, xyz)", size, center), mbb
}

func (s *Shader) processCirclePrimitive(exps []ast.Expression) (string, *MBB) {
	s.Primitives["circle"] = true
	args := s.getArgs(exps, "r", "d")

	radius := args[0]
	diameter := args[1]
	if diameter != "" && radius == "" {
		if vec3, err := parseVec3(diameter); err == nil {
			radius = fmt.Sprintf("%v", 0.5*vec3[0])
		}
	}
	if radius == "" {
		radius = "1"
	}

	vec3, err := parseVec3(radius)
	if err != nil {
		log.Printf("error parsing circle radius=%q, setting to 1", radius)
		radius = "1"
		vec3 = []float64{1, 1, 1}
	}

	mbb := &MBB{XMin: -vec3[0], YMin: -vec3[1], XMax: vec3[0], YMax: vec3[1]}

	return fmt.Sprintf("circle(float(%v), xyz)", radius), mbb
}

func (s *Shader) processMultmatrixBlockPrimitive(args []ast.Expression, exps []ast.Statement) (string, *MBB) {
	calls, mbb := s.getCalls(exps)
	if len(calls) == 0 || mbb == nil {
		return "", nil
	}

	fNum := len(s.Functions)
	fName := fmt.Sprintf("multimatrixBlock%v", fNum)
	vec4s := s.getMat4Args(args)
	newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	mat4 xfm = inverse(mat4(vec4(%v), vec4(%v), vec4(%v), vec4(%v)));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return %v;
}
`, fName, vec4s[0], vec4s[1], vec4s[2], vec4s[3], strings.Join(calls, " + "))
	s.Functions = append(s.Functions, newFunc)

	vec0, err := parseVec4(vec4s[0])
	if err != nil {
		log.Fatalf("vec0: %v", err)
	}
	vec1, err := parseVec4(vec4s[1])
	if err != nil {
		log.Fatalf("vec1: %v", err)
	}
	vec2, err := parseVec4(vec4s[2])
	if err != nil {
		log.Fatalf("vec2: %v", err)
	}
	vec3, err := parseVec4(vec4s[3])
	if err != nil {
		log.Fatalf("vec3: %v", err)
	}

	newMBB := matrixMult(mbb, vec0, vec1, vec2, vec3)

	return fmt.Sprintf("%v(xyz)", fName), newMBB
}

func (s *Shader) processUnionBlockPrimitive(exps []ast.Statement) (string, *MBB) {
	calls, mbb := s.getCalls(exps)
	if len(calls) == 0 || mbb == nil {
		return "", nil
	}

	fNum := len(s.Functions)
	fName := fmt.Sprintf("union%v", fNum)
	newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	return clamp(%v, 0.0, 1.0);
}
`, fName, strings.Join(calls, " + "))
	s.Functions = append(s.Functions, newFunc)

	return fmt.Sprintf("%v(xyz)", fName), mbb
}

func (s *Shader) processDifferenceBlockPrimitive(exps []ast.Statement) (string, *MBB) {
	calls, mbb := s.getCalls(exps)
	if len(calls) == 0 || mbb == nil {
		return "", nil
	}

	fNum := len(s.Functions)
	fName := fmt.Sprintf("difference%v", fNum)
	newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	return clamp(%v, 0.0, 1.0);
}
`, fName, strings.Join(calls, " - "))
	s.Functions = append(s.Functions, newFunc)

	return fmt.Sprintf("%v(xyz)", fName), mbb
}

func (s *Shader) processIntersectionBlockPrimitive(exps []ast.Statement) (string, *MBB) {
	calls, mbb := s.getCalls(exps)
	if len(calls) == 0 || mbb == nil {
		return "", nil
	}

	fNum := len(s.Functions)
	fName := fmt.Sprintf("intersection%v", fNum)
	newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	return clamp(%v, 0.0, 1.0);
}
`, fName, strings.Join(calls, " * "))
	s.Functions = append(s.Functions, newFunc)

	return fmt.Sprintf("%v(xyz)", fName), mbb
}

func (s *Shader) processLinearExtrudeBlockPrimitive(args []ast.Expression, exps []ast.Statement) (string, *MBB) {
	calls, mbb := s.getCalls(exps)
	if len(calls) == 0 || mbb == nil {
		return "", nil
	}

	argVals := s.getArgs(args, "height", "center", "twist", "scale")

	height, err := strconv.ParseFloat(argVals[0], 64)
	if err != nil {
		log.Fatalf("unable to parse height %q: %v", argVals[0], err)
	}

	if argVals[1] == "true" {
		mbb.ZMin = -0.5 * height
		mbb.ZMax = 0.5 * height
	} else {
		argVals[1] = "false"
		mbb.ZMin = 0.0
		mbb.ZMax = height
	}

	argVals[2] = strings.Trim(argVals[2], "()")
	argVals[3] = strings.Trim(argVals[3], "[]")
	scaleVec, err := parseVec2(argVals[3])
	if err != nil {
		log.Fatalf("unable to parse linear_extrude scale %q: %v", argVals[3], err)
	}

	fNum := len(s.Functions)
	fName := fmt.Sprintf("linearExtrudeBlock%v", fNum)
	var newFunc string
	if argVals[2] == "" { // No twist.
		newFunc = fmt.Sprintf(`float %v(in vec3 xyz) {
	xyz.z /= float(%v);
	float z = xyz.z;
	if (%v) { z += 0.5; } else { xyz.z -= 0.5; }
	if (abs(xyz.z) > 0.5) { return 0.0; }
	vec2 s = mix(vec2(1),vec2(%v,%v),z);
	xyz.xy /= s;
	return %v;
}
`, fName, argVals[0], argVals[1], scaleVec[0], scaleVec[1], strings.Join(calls, " + "))
		// Modify MBB based on scale.
		if scaleVec[0] > 1.0 {
			cx := 0.5 * (mbb.XMax + mbb.XMin)
			dx := 0.5 * (mbb.XMax - mbb.XMin)
			mbb.XMin = cx - scaleVec[0]*dx
			mbb.XMax = cx + scaleVec[0]*dx
		}
		if scaleVec[1] > 1.0 {
			cy := 0.5 * (mbb.YMax + mbb.YMin)
			dy := 0.5 * (mbb.YMax - mbb.YMin)
			mbb.YMin = cy - scaleVec[1]*dy
			mbb.YMax = cy + scaleVec[1]*dy
		}
	} else {
		// With twist.
		s.Primitives["rotAxis"] = true
		s.Primitives["rotZ"] = true

		newFunc = fmt.Sprintf(`float %v(in vec3 xyz) {
	xyz.z /= float(%v);
	float z = xyz.z;
	if (%v) { z += 0.5; } else { xyz.z -= 0.5; }
	if (abs(xyz.z) > 0.5) { return 0.0; }
	float angle = mix(0.0, float(%v)*3.1415926535897932384626433832795/180.0, z);
	vec2 s = mix(vec2(1),vec2(%v,%v),z);
	xyz.xy /= s;
	xyz = (vec4(xyz, 1) * rotZ(angle)).xyz;
	return %v;
}
`, fName, argVals[0], argVals[1], argVals[2], scaleVec[0], scaleVec[1], strings.Join(calls, " + "))

		// Modify MBB based on twist and scale.
		twist, err := strconv.ParseFloat(argVals[2], 64)
		if err != nil {
			log.Fatalf("unable to parse twist %q: %v", argVals[2], err)
		}

		cx := 0.5 * (mbb.XMax + mbb.XMin)
		dx := 0.5 * (mbb.XMax - mbb.XMin)
		cy := 0.5 * (mbb.YMax + mbb.YMin)
		dy := 0.5 * (mbb.YMax - mbb.YMin)

		const vertSteps = 33
		for t := 0.0; t <= 1.0; t += (1.0 / vertSteps) {
			sx := lerp(1.0, scaleVec[0], t)
			xmin := cx - sx*dx
			xmax := cx + sx*dx
			sy := lerp(1.0, scaleVec[1], t)
			ymin := cy - sy*dy
			ymax := cy + sy*dy
			rot := lerp(0.0, -twist*math.Pi/180.0, t)

			s := math.Sin(rot)
			c := math.Cos(rot)

			x1 := xmin*c - ymin*s
			y1 := xmin*s + ymin*c
			x2 := xmax*c - ymax*s
			y2 := xmax*s + ymax*c
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			m := &MBB{XMin: x1, YMin: y1, XMax: x2, YMax: y2}
			mbb.update(m)
		}
	}
	s.Functions = append(s.Functions, newFunc)

	return fmt.Sprintf("%v(xyz)", fName), mbb
}

func lerp(a, b, t float64) float64 {
	return (1.0-t)*a + t*b
}
