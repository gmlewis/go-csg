// Package irmf defines the Abstract Syntax Tree for IRMF.
// It translates an ast.Program to an irmf.Shader and
// can then output this shader as a valid .irmf file.
package irmf

import (
	"github.com/gmlewis/go-csg/ast"
)

// Shader represents an IRMF shader.
type Shader struct {
}

// String returns the strings representation of the IRMF Shader.
func (s *Shader) String() string {
	return ""
}

// New returns a new IRMF Shader from a CSG ast.Program.
func New(program *ast.Program) (*Shader, error) {
	return &Shader{}, nil
}
