package irmf

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gmlewis/go-csg/ast"
	"github.com/gmlewis/go-csg/token"
)

func TestGetArgs(t *testing.T) {
	size := &ast.NamedArgument{
		Name: &ast.StringLiteral{
			Token: token.Token{Literal: "size"},
		},
		Value: &ast.StringLiteral{
			Token: token.Token{Literal: "[1, 1, 1]"},
		},
	}
	center := &ast.NamedArgument{
		Name: &ast.StringLiteral{
			Token: token.Token{Literal: "center"},
		},
		Value: &ast.StringLiteral{
			Token: token.Token{Literal: "true"},
		},
	}

	tests := []struct {
		exps  []ast.Expression
		names []string
		want  []string
	}{
		{
			exps:  nil,
			names: []string{"size", "center"},
			want:  []string{"", ""},
		},
		{
			exps:  []ast.Expression{size},
			names: []string{"size", "center"},
			want:  []string{"[1, 1, 1]", ""},
		},
		{
			exps:  []ast.Expression{center},
			names: []string{"size", "center"},
			want:  []string{"", "true"},
		},
		{
			exps:  []ast.Expression{size, center},
			names: []string{"size", "center"},
			want:  []string{"[1, 1, 1]", "true"},
		},
		{
			exps:  []ast.Expression{center, size},
			names: []string{"size", "center"},
			want:  []string{"[1, 1, 1]", "true"},
		},
		{
			exps: []ast.Expression{
				&ast.StringLiteral{Token: token.Token{Literal: "[4, 5, 6]"}},
				&ast.StringLiteral{Token: token.Token{Literal: "true"}},
			},
			names: []string{"size", "center"},
			want:  []string{"[4, 5, 6]", "true"},
		},
		{
			exps: []ast.Expression{
				&ast.StringLiteral{Token: token.Token{Literal: "[4, 5, 6]"}},
				center,
			},
			names: []string{"size", "center"},
			want:  []string{"[4, 5, 6]", "true"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			s := &Shader{}
			got := s.getArgs(tt.exps, tt.names...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArgs = %+v, want %+v", got, tt.want)
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
