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
	MBB        *MBB
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
	s := &Shader{
		Program:    program,
		Primitives: map[string]bool{},
	}

	calls, mbb := s.getCalls(program.Statements)
	if len(calls) > 0 {
		mainFunc := fmt.Sprintf(`void mainModel4(out vec4 materials, in vec3 xyz) {
	materials[0] = %v;
}
`, strings.Join(calls, " + "))
		s.Functions = append(s.Functions, mainFunc)
		s.MBB = mbb
	}

	return s
}

// MBB represents a minimum bounding box.
type MBB struct {
	xmin, xmax float64
	ymin, ymax float64
	zmin, zmax float64
}

func (mbb *MBB) update(other *MBB) {
	if other.xmin < mbb.xmin {
		mbb.xmin = other.xmin
	}
	if other.ymin < mbb.ymin {
		mbb.ymin = other.ymin
	}
	if other.zmin < mbb.zmin {
		mbb.zmin = other.zmin
	}
	if other.xmax > mbb.xmax {
		mbb.xmax = other.xmax
	}
	if other.ymax > mbb.ymax {
		mbb.ymax = other.ymax
	}
	if other.zmax > mbb.zmax {
		mbb.zmax = other.zmax
	}
}

func (s *Shader) getCalls(stmts []ast.Statement) ([]string, *MBB) {
	var mbb *MBB
	var calls []string
	for _, stmt := range stmts {
		if call, callMBB := s.processStatement(stmt); call != "" {
			calls = append(calls, call)
			if mbb == nil {
				mbb = callMBB
			} else {
				mbb.update(callMBB)
			}
		}
	}
	return calls, mbb
}

func (s *Shader) processStatement(stmt ast.Statement) (string, *MBB) {
	switch node := stmt.(type) {
	case *ast.ExpressionStatement:
		return s.processExpression(node.Expression)
	default:
		log.Fatalf("unhandled statement type %T (%+v)", node, node)
	}
	return "", nil
}

func (s *Shader) processExpression(exp ast.Expression) (string, *MBB) {
	switch node := exp.(type) {
	case *ast.CallExpression:
		log.Printf("WARNING: node currently not supported. Skipping: %v", node.String())
	case *ast.CirclePrimitive:
		s.Primitives["circle"] = true
		// TODO: make a new function to call this primitive.
		return "circle(TODO)", nil
	case *ast.ColorBlockPrimitive: // Currently, color itself is a NOOP.
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("colorBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.CubePrimitive:
		return s.processCubePrimitive(node.Arguments)
	case *ast.CylinderPrimitive:
		s.Primitives["cylinder"] = true
		// TODO: make a new function to call this primitive.
		return "cylinder(TODO)", nil
	case *ast.DifferenceBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("differenceBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.GroupBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("groupBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(xyz)", fName), mbb
			}
		}
	case *ast.HullBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("hullBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.IntersectionBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("intersectionBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.LinearExtrudeBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a linear extrude.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("linearExtrudeBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.MinkowskiBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a minkowski.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("minkowskiBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.MultmatrixBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after a matrix multiply.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("multimatrixBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.PolygonPrimitive:
		s.Primitives["polygon"] = true
		// TODO: make a new function to call this primitive.
		return "polygon(TODO)", nil
	case *ast.PolyhedronPrimitive:
		s.Primitives["polyhedron"] = true
		// TODO: make a new function to call this primitive.
		return "polyhedron(TODO)", nil
	case *ast.ProjectionBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a projection.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("projectionBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.RotateExtrudeBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a rotate extrude.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("rotateExtrudeBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	case *ast.SpherePrimitive:
		s.Primitives["sphere"] = true
		// TODO: make a new function to call this primitive.
		return "sphere(TODO)", nil
	case *ast.SquarePrimitive:
		s.Primitives["square"] = true
		// TODO: make a new function to call this primitive.
		return "square(TODO)", nil
	case *ast.UnionBlockPrimitive:
		if node.Body != nil {
			// TODO: make a new function to call these statements after wrapping in a union.
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("unionBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(TODO) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(TODO)", fName), mbb
			}
		}
	default:
		log.Fatalf("unhandled expression type %T (%+v)", node, node)
	}
	return "", nil
}
