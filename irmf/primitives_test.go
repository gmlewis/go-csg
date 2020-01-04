package irmf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/go-csg/lexer"
	"github.com/gmlewis/go-csg/parser"
)

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

func TestProcessCubePrimitive(t *testing.T) {
	tests := []struct {
		src  string
		want []string
		mbb  *MBB
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

			shader := New(program)
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
		src  string
		want []string
		mbb  *MBB
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

			shader := New(program)
			if !reflect.DeepEqual(shader.Functions, tt.want) {
				t.Errorf("functions = %+v, want %+v", shader.Functions, tt.want)
			}
			if !reflect.DeepEqual(shader.MBB, tt.mbb) {
				t.Errorf("mbb = %+v, want %+v", shader.MBB, tt.mbb)
			}
		})
	}
}
