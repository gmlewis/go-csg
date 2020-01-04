package irmf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gmlewis/go-csg/lexer"
	"github.com/gmlewis/go-csg/parser"
)

func TestShader_String(t *testing.T) {
	tests := []struct {
		program string
		center  bool
		want    string
	}{
		{
			program: `group() {
				cube(size = [1, 1, 1], center = false);
			}
			`,
			want: primitives["cube"] + `
float groupBlock0(in vec3 xyz) {
	return cube(vec3(1, 1, 1), false, xyz);
}

void mainModel4(out vec4 materials, in vec3 xyz) {
	materials[0] = groupBlock0(xyz);
}
`,
		},
		{
			program: `group() {
				cube(size = [1, 1, 1], center = false);
			}
			`,
			center: true,
			want: primitives["cube"] + `
float groupBlock0(in vec3 xyz) {
	return cube(vec3(1, 1, 1), false, xyz);
}

void mainModel4(out vec4 materials, in vec3 xyz) {
	xyz += vec3(0.5, 0.5, 0.5);
	materials[0] = groupBlock0(xyz);
}
`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			le := lexer.New(tt.program)
			p := parser.New(le)
			program := p.ParseProgram()
			if errs := p.Errors(); len(errs) != 0 {
				t.Fatalf("ParseProgram: %v", strings.Join(errs, "\n"))
			}

			shader := New(program, tt.center)

			if got := shader.String(); got != tt.want {
				t.Errorf("shader.String =\n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
}
