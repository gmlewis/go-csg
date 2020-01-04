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
		case *ast.StringLiteral:
			result[count] = exp.String()
			count++
		case *ast.IntegerLiteral:
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
		mbb = &MBB{xmin: -0.5 * vec3[0], ymin: -0.5 * vec3[1], zmin: -0.5 * vec3[2], xmax: 0.5 * vec3[0], ymax: 0.5 * vec3[1], zmax: 0.5 * vec3[2]}
	} else {
		center = "false"
		mbb = &MBB{xmax: vec3[0], ymax: vec3[1], zmax: vec3[2]}
	}

	return fmt.Sprintf("cube(vec3(%v), %v, xyz)", size, center), mbb
}

var primitives = map[string]string{
	"circle": `float circle(in vec3 xyz) {
		// TODO
		return 1.0;
	}
	`,
	"cube": `float cube(in vec3 size, in bool center, in vec3 xyz) {
	xyz /= size;
	if (!center) { xyz -= vec3(0.5); }
	if (any(greaterThan(abs(xyz), vec3(0.5)))) { return 0.0; }
	return 1.0;
}
`,
	"cylinder": `float cylinder(in vec3 xyz) {
	// TODO
	return 1.0;
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
	"sphere": `float sphere(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,
	"square": `float square(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,
}
