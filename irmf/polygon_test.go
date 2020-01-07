package irmf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/go-csg/evaluator"
	"github.com/gmlewis/go-csg/lexer"
	"github.com/gmlewis/go-csg/parser"
)

func TestProcessSimplyPolygonPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src: "polygon(points = [[0, 0], [100, 0], [130, 50], [30, 50]], paths = undef, convexity = 1);",
			want: []string{
				`float simplePolygon0(in vec3 xyz) {
	if (any(lessThan(xyz.xy, vec2(0,0))) || any(greaterThan(xyz.xy, vec2(130,50)))) { return 0.0; }
	if (xyz.y >= float(0) && xyz.y <= float(50)) { return testTwoLineSegments(vec2(0,0),vec2(30,50),vec2(100,0),vec2(130,50),xyz.xy); }
	return 1.0;
}
`,
				fmt.Sprintf(mainBodyFmt, "simplePolygon0(xyz)"),
			},
			mbb: &MBB{XMin: 0, XMax: 130, YMin: 0, YMax: 50},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			le := lexer.New(tt.src)
			p := parser.New(le)
			program := p.ParseProgram()
			if errs := p.Errors(); len(errs) != 0 {
				t.Fatalf("ParseProgram: %v", strings.Join(errs, "\n"))
			}

			obj := evaluator.Eval(program, nil)

			shader := New(obj, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %#v, want %#v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}
