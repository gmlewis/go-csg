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
			mbb:  &MBB{xmax: 1, ymax: 1, zmax: 1},
		},
		{
			src:  "cube(2);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(2), false, xyz)")},
			mbb:  &MBB{xmax: 2, ymax: 2, zmax: 2},
		},
		{
			src:  "cube(center=true);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(1), true, xyz)")},
			mbb:  &MBB{xmin: -0.5, ymin: -0.5, zmin: -0.5, xmax: 0.5, ymax: 0.5, zmax: 0.5},
		},
		{
			src:  "cube(size=5);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(5), false, xyz)")},
			mbb:  &MBB{xmax: 5, ymax: 5, zmax: 5},
		},
		{
			src:  "cube(size= [ 5 , 4 , 3 ]);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(5, 4, 3), false, xyz)")},
			mbb:  &MBB{xmax: 5, ymax: 4, zmax: 3},
		},
		{
			src:  "cube(center = false, size = [ 5 , 4 , 3 ]);",
			want: []string{fmt.Sprintf(mainBodyFmt, "cube(vec3(5, 4, 3), false, xyz)")},
			mbb:  &MBB{xmax: 5, ymax: 4, zmax: 3},
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
