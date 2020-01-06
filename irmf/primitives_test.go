package irmf

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/go-csg/lexer"
	"github.com/gmlewis/go-csg/parser"
)

func TestParseVec2(t *testing.T) {
	tests := []struct {
		s    string
		want []float64
	}{
		{"1", []float64{1, 1}},
		{" 1, 2 ", []float64{1, 2}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got, err := parseVec2(tt.s)
			if err != nil {
				t.Fatalf("parseVec2 = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVec2 = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseVec3(t *testing.T) {
	tests := []struct {
		s    string
		want []float64
	}{
		{"1", []float64{1, 1, 1}},
		{" 1, 2, 3 ", []float64{1, 2, 3}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got, err := parseVec3(tt.s)
			if err != nil {
				t.Fatalf("parseVec3 = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVec3 = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseVec4(t *testing.T) {
	tests := []struct {
		s    string
		want []float64
	}{
		{"1", []float64{1, 1, 1, 1}},
		{" 1, 2, 3, 4 ", []float64{1, 2, 3, 4}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got, err := parseVec4(tt.s)
			if err != nil {
				t.Fatalf("parseVec4 = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVec4 = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestProcessCubePrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src:  "cube();",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(1), false, xyz)")},
			mbb:  &MBB{XMax: 1, YMax: 1, ZMax: 1},
		},
		{
			src:  "cube(2);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(2), false, xyz)")},
			mbb:  &MBB{XMax: 2, YMax: 2, ZMax: 2},
		},
		{
			src:  "cube(center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(1), true, xyz)")},
			mbb:  &MBB{XMin: -0.5, YMin: -0.5, ZMin: -0.5, XMax: 0.5, YMax: 0.5, ZMax: 0.5},
		},
		{
			src:  "cube(size=5);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(5), false, xyz)")},
			mbb:  &MBB{XMax: 5, YMax: 5, ZMax: 5},
		},
		{
			src:  "cube(size= [ 5 , 4 , 3 ]);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(5, 4, 3), false, xyz)")},
			mbb:  &MBB{XMax: 5, YMax: 4, ZMax: 3},
		},
		{
			src:  "cube(center = false, size = [ 5 , 4 , 3 ]);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(5, 4, 3), false, xyz)")},
			mbb:  &MBB{XMax: 5, YMax: 4, ZMax: 3},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %+v, want %+v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessSpherePrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src:  "sphere();",
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(1), xyz)")},
			mbb:  &MBB{XMin: -1, YMin: -1, ZMin: -1, XMax: 1, YMax: 1, ZMax: 1},
		},
		{
			src:  "sphere(2);",
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(2), xyz)")},
			mbb:  &MBB{XMin: -2, YMin: -2, ZMin: -2, XMax: 2, YMax: 2, ZMax: 2},
		},
		{
			src:  "sphere(r = 2);",
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(2), xyz)")},
			mbb:  &MBB{XMin: -2, YMin: -2, ZMin: -2, XMax: 2, YMax: 2, ZMax: 2},
		},
		{
			src:  "sphere(r=5);",
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(5), xyz)")},
			mbb:  &MBB{XMin: -5, YMin: -5, ZMin: -5, XMax: 5, YMax: 5, ZMax: 5},
		},
		{
			src:  "sphere(r=3.14);",
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(3.14), xyz)")},
			mbb:  &MBB{XMin: -3.14, YMin: -3.14, ZMin: -3.14, XMax: 3.14, YMax: 3.14, ZMax: 3.14},
		},
		{
			src:  "sphere(d = 2);",
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(1), xyz)")},
			mbb:  &MBB{XMin: -1, YMin: -1, ZMin: -1, XMax: 1, YMax: 1, ZMax: 1},
		},
		{
			src:  "sphere(d = 20, r=1);", // r overrides d.
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(1), xyz)")},
			mbb:  &MBB{XMin: -1, YMin: -1, ZMin: -1, XMax: 1, YMax: 1, ZMax: 1},
		},
		{
			src:  "sphere($fn = 0, $fa = 12, $fs = 2, r = 1);",
			want: []string{fmt.Sprintf(mainBodyFmt, "sphere(float(1), xyz)")},
			mbb:  &MBB{XMin: -1, YMin: -1, ZMin: -1, XMax: 1, YMax: 1, ZMax: 1},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %+v, want %+v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessCylinderPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src:  "cylinder();",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(1), float(1), float(1), false, xyz)")},
			mbb:  &MBB{XMin: -1, XMax: 1, YMin: -1, YMax: 1, ZMin: 0, ZMax: 1},
		},

		// Equivalent:
		{
			src:  "cylinder(h=15, r1=9.5, r2=19.5, center=false);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(9.5), float(19.5), false, xyz)")},
			mbb:  &MBB{XMin: -19.5, XMax: 19.5, YMin: -19.5, YMax: 19.5, ZMin: 0, ZMax: 15},
		},
		{
			src:  "cylinder(  15,    9.5,    19.5, false);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(9.5), float(19.5), false, xyz)")},
			mbb:  &MBB{XMin: -19.5, XMax: 19.5, YMin: -19.5, YMax: 19.5, ZMin: 0, ZMax: 15},
		},
		{
			src:  "cylinder(  15,    9.5,    19.5);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(9.5), float(19.5), false, xyz)")},
			mbb:  &MBB{XMin: -19.5, XMax: 19.5, YMin: -19.5, YMax: 19.5, ZMin: 0, ZMax: 15},
		},
		{
			src:  "cylinder(  15,    9.5, d2=39  );",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(9.5), float(19.5), false, xyz)")},
			mbb:  &MBB{XMin: -19.5, XMax: 19.5, YMin: -19.5, YMax: 19.5, ZMin: 0, ZMax: 15},
		},
		{
			src:  "cylinder(  15, d1=19,  d2=39  );",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(9.5), float(19.5), false, xyz)")},
			mbb:  &MBB{XMin: -19.5, XMax: 19.5, YMin: -19.5, YMax: 19.5, ZMin: 0, ZMax: 15},
		},
		{
			src:  "cylinder(  15, d1=19,  r2=19.5);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(9.5), float(19.5), false, xyz)")},
			mbb:  &MBB{XMin: -19.5, XMax: 19.5, YMin: -19.5, YMax: 19.5, ZMin: 0, ZMax: 15},
		},

		// Equivalent:
		{
			src:  "cylinder(h=15, r1=10, r2=0, center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(10), float(0), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -7.5, ZMax: 7.5},
		},
		{
			src:  "cylinder(  15,    10,    0,        true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(10), float(0), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -7.5, ZMax: 7.5},
		},
		{
			src:  "cylinder(h=15, d1=20, d2=0, center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(15), float(10), float(0), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -7.5, ZMax: 7.5},
		},

		// Equivalent:
		{
			src:  "cylinder(h=20, r=10, center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(20), float(10), float(10), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -10, ZMax: 10},
		},
		{
			src:  "cylinder(  20,   10, 10,true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(20), float(10), float(10), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -10, ZMax: 10},
		},
		{
			src:  "cylinder(  20, d=20, center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(20), float(10), float(10), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -10, ZMax: 10},
		},
		{
			src:  "cylinder(  20,r1=10, d2=20, center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(20), float(10), float(10), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -10, ZMax: 10},
		},
		{
			src:  "cylinder(  20,r1=10, d2=20, center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cylinder(float(20), float(10), float(10), true, xyz)")},
			mbb:  &MBB{XMin: -10, XMax: 10, YMin: -10, YMax: 10, ZMin: -10, ZMax: 10},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %+v, want %+v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessSquarePrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src:  "square();",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(1), false, xyz)")},
			mbb:  &MBB{XMax: 1, YMax: 1},
		},
		{
			src:  "square([20,10],true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(20, 10), true, xyz)")},
			mbb:  &MBB{XMin: -10, YMin: -5, XMax: 10, YMax: 5},
		},
		{
			src:  "square(size = 10);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
		},
		{
			src:  "square(10);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
		},
		{
			src:  "square([10,10]);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10, 10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
		},
		{
			src:  "square(10,false);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
		},
		{
			src:  "square([10,10],false);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10, 10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
		},
		{
			src:  "square([10,10],center=false);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10, 10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
		},
		{
			src:  "square(size = [10, 10], center = false);",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10, 10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
		},
		{
			src:  "square(center = false,size = [10, 10] );",
			want: []string{fmt.Sprintf(mainBodyFmt, "square(vec2(10, 10), false, xyz)")},
			mbb:  &MBB{XMax: 10, YMax: 10},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %+v, want %+v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessMultmatrixBlockPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src: "multmatrix([[1, 0, 0, -19], [0, 1, 0, -0.5], [0, 0, 1, 0], [0, 0, 0, 1]]) {sphere();}",
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, 19), vec4(0, 1, 0, 0.5), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return sphere(float(1), xyz);
}
`,
				fmt.Sprintf(mainBodyFmt, "multimatrixBlock0(xyz)"),
			},
			mbb: &MBB{XMin: -20, XMax: -18, YMin: -1.5, YMax: 0.5, ZMin: -1, ZMax: 1},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %#v, want %#v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessUnionBlockPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src: `union() {
	sphere($fn = 100, $fa = 12, $fs = 2, r = 1);
	multmatrix([[1, 0, 0, 2], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) {
		sphere($fn = 100, $fa = 12, $fs = 2, r = 2);
	}
}
`,
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, -2), vec4(0, 1, 0, 0), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return sphere(float(2), xyz);
}
`,
				"float union1(in vec3 xyz) {\n\treturn clamp(sphere(float(1), xyz) + multimatrixBlock0(xyz), 0.0, 1.0);\n}\n",
				fmt.Sprintf(mainBodyFmt, "union1(xyz)"),
			},
			mbb: &MBB{XMin: -1, XMax: 4, YMin: -2, YMax: 2, ZMin: -2, ZMax: 2},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %#v, want %#v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessDifferenceBlockPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src: `difference() {
	sphere($fn = 100, $fa = 12, $fs = 2, r = 1);
	multmatrix([[1, 0, 0, 2], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) {
		sphere($fn = 100, $fa = 12, $fs = 2, r = 2);
	}
}
`,
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, -2), vec4(0, 1, 0, 0), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return sphere(float(2), xyz);
}
`,
				"float difference1(in vec3 xyz) {\n\treturn clamp(sphere(float(1), xyz) - multimatrixBlock0(xyz), 0.0, 1.0);\n}\n",
				fmt.Sprintf(mainBodyFmt, "difference1(xyz)"),
			},
			mbb: &MBB{XMin: -1, XMax: 4, YMin: -2, YMax: 2, ZMin: -2, ZMax: 2},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %#v, want %#v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessIntersectionBlockPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src: `intersection() {
	sphere($fn = 100, $fa = 12, $fs = 2, r = 1);
	multmatrix([[1, 0, 0, 2], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) {
		sphere($fn = 100, $fa = 12, $fs = 2, r = 2);
	}
}
`,
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, -2), vec4(0, 1, 0, 0), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return sphere(float(2), xyz);
}
`,
				"float intersection1(in vec3 xyz) {\n\treturn clamp(sphere(float(1), xyz) * multimatrixBlock0(xyz), 0.0, 1.0);\n}\n",
				fmt.Sprintf(mainBodyFmt, "intersection1(xyz)"),
			},
			mbb: &MBB{XMin: -1, XMax: 4, YMin: -2, YMax: 2, ZMin: -2, ZMax: 2},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %#v, want %#v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessLinearExtrudeBlockPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src: `linear_extrude(height = 10, center = false, convexity = 1, scale = [1, 1], $fn = 0, $fa = 12, $fs = 2) {
	multmatrix([[1, 0, 0, 1], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) {
		circle($fn = 100, $fa = 12, $fs = 2, r = 1);
	}
}`,
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, -1), vec4(0, 1, 0, 0), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return circle(float(1), xyz);
}
`,
				`float linearExtrudeBlock1(in vec3 xyz) {
	xyz.z /= float(10);
	float z = xyz.z;
	if (false) { z += 0.5; } else { xyz.z -= 0.5; }
	if (abs(xyz.z) > 0.5) { return 0.0; }
	vec2 s = mix(vec2(1),vec2(1,1),z);
	xyz.xy /= s;
	return multimatrixBlock0(xyz);
}
`,
				fmt.Sprintf(mainBodyFmt, "linearExtrudeBlock1(xyz)"),
			},
			mbb: &MBB{XMin: 0, XMax: 2, YMin: -1, YMax: 1, ZMin: 0, ZMax: 10},
		},
		{
			src: `linear_extrude(height = 10, center = true, convexity = 1, scale = [1, 1], $fn = 0, $fa = 12, $fs = 2) {
	multmatrix([[1, 0, 0, 1], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) {
		circle($fn = 100, $fa = 12, $fs = 2, r = 1);
	}
}`,
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, -1), vec4(0, 1, 0, 0), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return circle(float(1), xyz);
}
`,
				`float linearExtrudeBlock1(in vec3 xyz) {
	xyz.z /= float(10);
	float z = xyz.z;
	if (true) { z += 0.5; } else { xyz.z -= 0.5; }
	if (abs(xyz.z) > 0.5) { return 0.0; }
	vec2 s = mix(vec2(1),vec2(1,1),z);
	xyz.xy /= s;
	return multimatrixBlock0(xyz);
}
`,
				fmt.Sprintf(mainBodyFmt, "linearExtrudeBlock1(xyz)"),
			},
			mbb: &MBB{XMin: 0, XMax: 2, YMin: -1, YMax: 1, ZMin: -5, ZMax: 5},
		},
		{
			src: `linear_extrude(height = 10, center = false, convexity = 1, twist = 90, slices = 7, scale = [1, 1], $fn = 0, $fa = 12, $fs = 2) {
	multmatrix([[1, 0, 0, 1], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) {
		circle($fn = 100, $fa = 12, $fs = 2, r = 1);
	}
}`,
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, -1), vec4(0, 1, 0, 0), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return circle(float(1), xyz);
}
`,
				`float linearExtrudeBlock1(in vec3 xyz) {
	xyz.z /= float(10);
	float z = xyz.z;
	if (false) { z += 0.5; } else { xyz.z -= 0.5; }
	if (abs(xyz.z) > 0.5) { return 0.0; }
	float angle = mix(0.0, float(90)*3.1415926535897932384626433832795/180.0, z);
	vec2 s = mix(vec2(1),vec2(1,1),z);
	xyz.xy /= s;
	xyz = (vec4(xyz, 1) * rotZ(angle)).xyz;
	return multimatrixBlock0(xyz);
}
`,
				fmt.Sprintf(mainBodyFmt, "linearExtrudeBlock1(xyz)"),
			},
			mbb: &MBB{XMin: -1, XMax: 2, YMin: -2, YMax: 1, ZMin: 0, ZMax: 10},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %#v, want %#v", shader.Functions, tt.want)
			}

			// X and Y values are approximate.
			const tol = 0.5
			dxmin := math.Abs(tt.mbb.XMin - shader.MBB.XMin)
			dxmax := math.Abs(tt.mbb.XMax - shader.MBB.XMax)
			dymin := math.Abs(tt.mbb.YMin - shader.MBB.YMin)
			dymax := math.Abs(tt.mbb.YMax - shader.MBB.YMax)
			dzmin := math.Abs(tt.mbb.ZMin - shader.MBB.ZMin)
			dzmax := math.Abs(tt.mbb.ZMax - shader.MBB.ZMax)
			if dxmin > tol || dxmax > tol || dymin > tol || dymax > tol || dzmin > tol || dzmax > tol {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}

func TestProcessRotateExtrudeBlockPrimitive(t *testing.T) {
	tests := []struct {
		src    string
		center bool
		want   []string
		mbb    *MBB
	}{
		{
			src: `rotate_extrude(convexity = 2, $fn = 100, $fa = 12, $fs = 2) {
	multmatrix([[1, 0, 0, 1], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) {
		circle($fn = 100, $fa = 12, $fs = 2, r = 1);
	}
}`,
			want: []string{
				`float multimatrixBlock0(in vec3 xyz) {
	mat4 xfm = mat4(vec4(1, 0, 0, -1), vec4(0, 1, 0, 0), vec4(0, 0, 1, 0), vec4(0, 0, 0, 1));
	xyz = (vec4(xyz, 1.0) * xfm).xyz;
	return circle(float(1), xyz);
}
`,
				`float rotateExtrudeBlock1(in vec3 xyz) {
	float angle = atan(xyz.y, xyz.x);
	if (angle<0.) { angle+=(2.*3.1415926535897932384626433832795); }
	if (angle>float(360)*3.1415926535897932384626433832795/180.0) { return 0.0; }
	vec3 slice=(vec4(xyz,1)*rotZ(-angle)).xyz;
	xyz = slice.xzy;
	return multimatrixBlock0(xyz);
}
`,
				fmt.Sprintf(mainBodyFmt, "rotateExtrudeBlock1(xyz)"),
			},
			mbb: &MBB{XMin: -2, XMax: 2, YMin: -2, YMax: 2, ZMin: -1, ZMax: 1},
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

			shader := New(program, tt.center)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %#v, want %#v", shader.Functions, tt.want)
			}

			// X and Y values are approximate.
			const tol = 0.5
			dxmin := math.Abs(tt.mbb.XMin - shader.MBB.XMin)
			dxmax := math.Abs(tt.mbb.XMax - shader.MBB.XMax)
			dymin := math.Abs(tt.mbb.YMin - shader.MBB.YMin)
			dymax := math.Abs(tt.mbb.YMax - shader.MBB.YMax)
			dzmin := math.Abs(tt.mbb.ZMin - shader.MBB.ZMin)
			dzmax := math.Abs(tt.mbb.ZMax - shader.MBB.ZMax)
			if dxmin > tol || dxmax > tol || dymin > tol || dymax > tol || dzmin > tol || dzmax > tol {
				t.Errorf("mbb = %#v, want %#v", shader.MBB, tt.mbb)
			}
		})
	}
}
