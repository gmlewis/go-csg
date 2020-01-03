// Package irmf defines the Abstract Syntax Tree for IRMF.
// It translates an ast.Program to an irmf.Shader and
// can then output this shader as a valid .irmf file.
package irmf

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gmlewis/go-csg/ast"
)

// Shader represents an IRMF shader.
type Shader struct {
	Program    *ast.Program
	Functions  []string
	Primitives map[string]bool
}

// String returns the strings representation of the IRMF Shader.
func (s *Shader) String() string {
	var result []string

	// First, output all used primitives, sorted by name (for consistency).
	var primNames []string
	for key := range s.Primitives {
		primNames = append(primNames, key)
	}
	sort.Strings(primNames)
	for _, primName := range primNames {
		result = append(result, primitives[primName])
	}

	result = append(result, s.Functions...)

	return strings.Join(result, "\n")
}

// New returns a new IRMF Shader from a CSG ast.Program.
func New(program *ast.Program) *Shader {
	shader := &Shader{
		Program:    program,
		Primitives: map[string]bool{},
	}

	var calls []string
	for _, stmt := range program.Statements {
		if call := shader.processStatement(stmt); call != "" {
			calls = append(calls, call)
		}
	}

	if len(calls) > 0 {
		mainFunc := fmt.Sprintf(`void mainModel4(out vec4 materials, in vec3 xyz) {
	materials[0] = %v;
}
`, strings.Join(calls, " + "))
		shader.Functions = append(shader.Functions, mainFunc)
	}

	return shader
}

func (s *Shader) processStatement(stmt ast.Statement) string {
	switch node := stmt.(type) {
	case *ast.ExpressionStatement:
		return s.processExpression(node.Expression)
	default:
		log.Fatalf("unhandled statement type %T (%+v)", node, node)
	}
	return ""
}

func (s *Shader) processExpression(exp ast.Expression) string {
	switch node := exp.(type) {
	case *ast.CallExpression:
		log.Printf("WARNING: node currently not supported. Skipping: %v", node.String())
	case *ast.CirclePrimitive:
		s.Primitives["circle"] = true
		// TODO: make a new function to call this primitive.
		return "circle(TODO)"
	case *ast.ColorBlockPrimitive: // Currently, color itself is a NOOP.
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("colorBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.CubePrimitive:
		s.Primitives["cube"] = true
		var size, center string
		for _, exp := range node.Arguments {
			arg := exp.String()
			switch {
			case strings.HasPrefix(arg, "size = ["):
				size = arg[8 : len(arg)-1]
			case strings.HasPrefix(arg, "center = "):
				center = arg[9:]
			default:
				log.Printf("arg=%v", arg)
			}
		}
		return fmt.Sprintf("cube(vec3(%v), %v, xyz)", size, center)
	case *ast.CylinderPrimitive:
		s.Primitives["cylinder"] = true
		// TODO: make a new function to call this primitive.
		return "cylinder(TODO)"
	case *ast.DifferenceBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("differenceBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.GroupBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("groupBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(xyz)", fName)
			}
		}
	case *ast.HullBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("hullBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.IntersectionBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("intersectionBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.LinearExtrudeBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a linear extrude.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("linearExtrudeBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.MinkowskiBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a minkowski.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("minkowskiBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.MultmatrixBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after a matrix multiply.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("multimatrixBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.PolygonPrimitive:
		s.Primitives["polygon"] = true
		// TODO: make a new function to call this primitive.
		return "polygon(TODO)"
	case *ast.PolyhedronPrimitive:
		s.Primitives["polyhedron"] = true
		// TODO: make a new function to call this primitive.
		return "polyhedron(TODO)"
	case *ast.ProjectionBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a projection.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("projectionBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.RotateExtrudeBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a rotate extrude.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("rotateExtrudeBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	case *ast.SpherePrimitive:
		s.Primitives["sphere"] = true
		// TODO: make a new function to call this primitive.
		return "sphere(TODO)"
	case *ast.SquarePrimitive:
		s.Primitives["square"] = true
		// TODO: make a new function to call this primitive.
		return "square(TODO)"
	case *ast.UnionBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a union.
			var calls []string
			for _, stmt := range node.Body.Statements {
				if call := s.processStatement(stmt); call != "" {
					calls = append(calls, call)
				}
			}
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("unionBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName)
			}
		}
	default:
		log.Fatalf("unhandled expression type %T (%+v)", node, node)
	}
	return ""
}
