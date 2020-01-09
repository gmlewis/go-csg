// Package irmf defines the Abstract Syntax Tree for IRMF.
// It translates an ast.Program to an irmf.Shader and
// can then output this shader as a valid .irmf file.
package irmf

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strings"

	"github.com/gmlewis/go-csg/ast"
	"github.com/gmlewis/go-csg/object"
)

const (
	mainBodyFmt = `void mainModel4(out vec4 materials, in vec3 xyz) {
	materials[0] = %v;
}
`
	mainBodyCenterFmt = `void mainModel4(out vec4 materials, in vec3 xyz) {
	xyz += vec3(%v, %v, %v);
	materials[0] = %v;
}
`
)

// Shader represents an IRMF shader.
type Shader struct {
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

// New returns a new IRMF Shader from a CSG Object.
func New(obj object.Object, center bool) *Shader {
	s := &Shader{
		Primitives: map[string]bool{},
	}

	calls, mbb := s.processObject(obj)
	// calls, mbb := s.getCalls(program.Statements)
	if len(calls) > 0 {
		mainFunc := fmt.Sprintf(mainBodyFmt, strings.Join(calls, " + "))
		if center {
			cx := 0.5 * (mbb.XMax + mbb.XMin)
			cy := 0.5 * (mbb.YMax + mbb.YMin)
			cz := 0.5 * (mbb.ZMax + mbb.ZMin)
			dx := 0.5 * math.Round(0.5+mbb.XMax-mbb.XMin) // Round up
			dy := 0.5 * math.Round(0.5+mbb.YMax-mbb.YMin)
			dz := 0.5 * math.Round(0.5+mbb.ZMax-mbb.ZMin)
			mbb = &MBB{XMin: -dx, YMin: -dy, ZMin: -dz, XMax: dx, YMax: dy, ZMax: dz}

			// Modify mainFunc to center the design.
			mainFunc = fmt.Sprintf(mainBodyCenterFmt, cx, cy, cz, strings.Join(calls, " + "))
		}

		s.Functions = append(s.Functions, mainFunc)
		s.MBB = mbb
	}

	return s
}

// MBB represents a minimum bounding box.
type MBB struct {
	XMin, XMax float64
	YMin, YMax float64
	ZMin, ZMax float64
}

func (mbb *MBB) update(other *MBB) {
	if other == nil {
		return
	}
	if other.XMin < mbb.XMin {
		mbb.XMin = other.XMin
	}
	if other.YMin < mbb.YMin {
		mbb.YMin = other.YMin
	}
	if other.ZMin < mbb.ZMin {
		mbb.ZMin = other.ZMin
	}
	if other.XMax > mbb.XMax {
		mbb.XMax = other.XMax
	}
	if other.YMax > mbb.YMax {
		mbb.YMax = other.YMax
	}
	if other.ZMax > mbb.ZMax {
		mbb.ZMax = other.ZMax
	}
}

func (s *Shader) processObject(obj object.Object) ([]string, *MBB) {
	switch obj := obj.(type) {
	case *object.CubePrimitive:
		return s.processCubePrimitiveObject(obj.Arguments)

	case *object.CylinderPrimitive:
		return s.processCylinderPrimitiveObject(obj.Arguments)

	case *object.GroupBlockPrimitive:
		if obj.Body != nil {
			return s.processGroupBlockPrimitiveObject(obj.Body)
		}

	case *object.MultmatrixBlockPrimitive:
		if obj.Body != nil {
			return s.processMultmatrixBlockPrimitiveObject(obj.Arguments, obj.Body)
		}

	case *object.PolygonPrimitive:
		return s.processPolygonPrimitiveObject(obj.Arguments)

	case *object.SpherePrimitive:
		return s.processSpherePrimitiveObject(obj.Arguments)

	case *object.SquarePrimitive:
		return s.processSquarePrimitiveObject(obj.Arguments)

	default:
		log.Fatalf("unhandled object type %T (%+v)", obj, obj)
	}

	return nil, nil
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
		return s.processCirclePrimitive(node.Arguments)
	case *ast.ColorBlockPrimitive: // Currently, color itself is a NOOP.
		if node.Body != nil {
			calls, mbb := s.getCalls(node.Body.Statements)
			if len(calls) > 0 {
				fNum := len(s.Functions)
				fName := fmt.Sprintf("colorBlock%v", fNum)
				newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	return %v;
}
`, fName, strings.Join(calls, " + "))
				s.Functions = append(s.Functions, newFunc)
				return fmt.Sprintf("%v(xyz)", fName), mbb
			}
		}
	case *ast.CubePrimitive:
		return s.processCubePrimitive(node.Arguments)
	case *ast.CylinderPrimitive:
		return s.processCylinderPrimitive(node.Arguments)
	case *ast.DifferenceBlockPrimitive:
		if node.Body != nil {
			return s.processDifferenceBlockPrimitive(node.Body.Statements)
		}
	case *ast.GroupBlockPrimitive:
		if node.Body != nil {
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
			return s.processIntersectionBlockPrimitive(node.Body.Statements)
		}
	case *ast.LinearExtrudeBlockPrimitive:
		if node.Body != nil {
			return s.processLinearExtrudeBlockPrimitive(node.Arguments, node.Body.Statements)
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
			return s.processMultmatrixBlockPrimitive(node.Arguments, node.Body.Statements)
		}
	case *ast.PolygonPrimitive:
		return s.processPolygonPrimitive(node.Arguments)
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
			return s.processRotateExtrudeBlockPrimitive(node.Arguments, node.Body.Statements)
		}
	case *ast.SpherePrimitive:
		return s.processSpherePrimitive(node.Arguments)
	case *ast.SquarePrimitive:
		return s.processSquarePrimitive(node.Arguments)
	case *ast.UnionBlockPrimitive:
		if node.Body != nil {
			return s.processUnionBlockPrimitive(node.Body.Statements)
		}
	default:
		log.Fatalf("unhandled expression type %T (%+v)", node, node)
	}
	return "", nil
}
