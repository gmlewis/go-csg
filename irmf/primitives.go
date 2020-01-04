package irmf

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gmlewis/go-csg/ast"
)

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
		return []float64{v, v, v}, nil
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
