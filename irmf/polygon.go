package irmf

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gmlewis/go-csg/ast"
	"github.com/gmlewis/go-csg/evaluator"
	"github.com/gmlewis/go-csg/object"
)

func (s *Shader) processPolygonPrimitiveObject(objs []object.Object) ([]string, *MBB) {
	args := s.getArgObjects(objs, "points", "paths")

	points, ok := args[0].(*object.Array)
	if !ok || points == nil {
		log.Fatalf("missing polygon points")
	}

	if args[1] == nil {
		return s.processSimplePolygonPrimitiveObject(points)
	}

	// TODO: Support "paths".

	return nil, nil
}

func (s *Shader) processPolygonPrimitive(exps []ast.Expression) (string, *MBB) {
	var points, paths *object.Array

	for _, exp := range exps {
		switch exp := exp.(type) {
		case *ast.NamedArgument:
			switch exp.Name.String() {
			case "points":
				if val, ok := exp.Value.(*ast.ArrayLiteral); ok {
					obj := evaluator.Eval(val, nil)
					if v, ok := obj.(*object.Array); ok {
						points = v
					} else {
						log.Fatalf("polygon unexpected points type %T (%+v)", obj, obj)
					}
				} else {
					log.Fatalf("polygon unexpected points type %T (%+v)", exp.Value, exp.Value)
				}
			case "paths":
				if val, ok := exp.Value.(*ast.ArrayLiteral); ok {
					obj := evaluator.Eval(val, nil)
					if v, ok := obj.(*object.Array); ok {
						paths = v
					} else {
						log.Fatalf("polygon unexpected paths type %T (%+v)", obj, obj)
					}
				}
			}
		default:
			log.Fatalf("exp=%T (%+v)", exp, exp)
		}
	}

	if points == nil {
		log.Fatalf("missing polygon points")
	}

	if paths == nil {
		return s.processSimplePolygonPrimitive(points)
	}

	// TODO: Support "paths".

	return "", nil
}

type ptT struct {
	x, y float64
}

func getPT(array *object.Array) ptT {
	if len(array.Elements) != 2 {
		log.Fatalf("polygon expected 2 elements per point: %+v", array)
	}
	result := ptT{}

	switch x := array.Elements[0].(type) {
	case *object.Float:
		result.x = x.Value
	case *object.Integer:
		result.x = float64(x.Value)
	default:
		log.Fatalf("polygon unexpected point x type %T %+v", x, x)
	}

	switch y := array.Elements[1].(type) {
	case *object.Float:
		result.y = y.Value
	case *object.Integer:
		result.y = float64(y.Value)
	default:
		log.Fatalf("polygon unexpected point y type %T %+v", y, y)
	}

	return result
}

func (s *Shader) processSimplePolygonPrimitiveObject(points *object.Array) ([]string, *MBB) {
	var xvals, yvals []float64
	var pts []ptT
	for _, el := range points.Elements {
		switch el := el.(type) {
		case *object.Array:
			pt := getPT(el)
			pts = append(pts, pt)
			xvals = append(xvals, pt.x)
			yvals = append(yvals, pt.y)
		default:
			log.Fatalf("polygon unexpected element type %T (%+v)", el, el)
		}
	}

	if len(xvals) < 3 || len(yvals) < 3 {
		log.Fatalf("polygon expected to have a least 3 points")
	}
	sort.Float64s(xvals)
	sort.Float64s(yvals)

	xmin := xvals[0]
	xmax := xvals[len(xvals)-1]
	ymin := yvals[0]
	ymax := yvals[len(yvals)-1]
	mbb := &MBB{XMin: xmin, YMin: ymin, XMax: xmax, YMax: ymax}

	fNum := len(s.Functions)
	fName := fmt.Sprintf("simplePolygon%v", fNum)

	s.Primitives["testTwoLineSegments"] = true

	lines := processSimplePolygonSegments(pts, yvals)

	newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	if (any(lessThan(xyz.xy, vec2(%v,%v))) || any(greaterThan(xyz.xy, vec2(%v,%v)))) { return 0.0; }
	%v
	return 1.0;
}
`, fName, xmin, ymin, xmax, ymax, strings.Join(lines, "\n\t"))
	s.Functions = append(s.Functions, newFunc)

	return []string{fmt.Sprintf("%v(xyz)", fName)}, mbb
}

func (s *Shader) processSimplePolygonPrimitive(points *object.Array) (string, *MBB) {
	var xvals, yvals []float64
	var pts []ptT
	for _, el := range points.Elements {
		switch el := el.(type) {
		case *object.Array:
			pt := getPT(el)
			pts = append(pts, pt)
			xvals = append(xvals, pt.x)
			yvals = append(yvals, pt.y)
		default:
			log.Fatalf("polygon unexpected element type %T (%+v)", el, el)
		}
	}

	if len(xvals) < 3 || len(yvals) < 3 {
		log.Fatalf("polygon expected to have a least 3 points")
	}
	sort.Float64s(xvals)
	sort.Float64s(yvals)

	xmin := xvals[0]
	xmax := xvals[len(xvals)-1]
	ymin := yvals[0]
	ymax := yvals[len(yvals)-1]
	mbb := &MBB{XMin: xmin, YMin: ymin, XMax: xmax, YMax: ymax}

	fNum := len(s.Functions)
	fName := fmt.Sprintf("simplePolygon%v", fNum)

	s.Primitives["testTwoLineSegments"] = true

	lines := processSimplePolygonSegments(pts, yvals)

	newFunc := fmt.Sprintf(`float %v(in vec3 xyz) {
	if (any(lessThan(xyz.xy, vec2(%v,%v))) || any(greaterThan(xyz.xy, vec2(%v,%v)))) { return 0.0; }
	%v
	return 1.0;
}
`, fName, xmin, ymin, xmax, ymax, strings.Join(lines, "\n\t"))
	s.Functions = append(s.Functions, newFunc)

	return fmt.Sprintf("%v(xyz)", fName), mbb
}

func processSimplePolygonSegments(pts []ptT, yvals []float64) []string {
	var result []string

	for i := 0; i < len(yvals)-1; i++ {
		if line := processSegment(yvals[i], yvals[i+1], pts); line != "" {
			result = append(result, line)
		}
	}

	return result
}

type segT struct {
	ly ptT
	uy ptT
}

func processSegment(ly, uy float64, pts []ptT) string {
	var leftSeg *segT
	var rightSeg *segT

	for i := 0; i < len(pts); i++ {
		lypt := pts[i]
		uypt := pts[(i+1)%len(pts)]
		if lypt.y == uypt.y {
			continue // horizontal line
		}
		if lypt.y > uypt.y {
			lypt, uypt = uypt, lypt
		}

		if ly >= uypt.y || uy <= lypt.y {
			continue
		}

		if leftSeg == nil {
			leftSeg = &segT{ly: lypt, uy: uypt}
		} else if rightSeg == nil {
			rightSeg = &segT{ly: lypt, uy: uypt}
			if lypt.x < leftSeg.ly.x || uypt.x < leftSeg.uy.x {
				// TODO: Check if lines cross?!?
				leftSeg, rightSeg = rightSeg, leftSeg
			}
		} else {
			log.Fatalf("concave polygon not yet supported: pts=%+v", pts)
		}
	}

	if leftSeg == nil || rightSeg == nil {
		return ""
	}

	return fmt.Sprintf("if (xyz.y >= float(%v) && xyz.y <= float(%v)) { return testTwoLineSegments(vec2(%v,%v),vec2(%v,%v),vec2(%v,%v),vec2(%v,%v),xyz.xy); }",
		ly, uy,
		leftSeg.ly.x, leftSeg.ly.y, leftSeg.uy.x, leftSeg.uy.y,
		rightSeg.ly.x, rightSeg.ly.y, rightSeg.uy.x, rightSeg.uy.y,
	)
}
