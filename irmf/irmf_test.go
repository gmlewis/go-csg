package irmf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/go-csg/lexer"
	"github.com/gmlewis/go-csg/parser"
)

func TestNew(t *testing.T) {
	tests := []struct {
		program string
		want    *Shader
	}{
		{
			program: "cube()",
			want:    &Shader{},
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

			got, err := New(program)
			if err != nil {
				t.Fatalf("New: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New =\n%+v\nwant\n%+v", got, tt.want)
			}
		})
	}
}
