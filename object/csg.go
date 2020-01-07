package object

import (
	"bytes"
	"fmt"
	"strings"
)

// CubePrimitive represents an object of that type.
type CubePrimitive struct {
	Arguments []Object
}

// Inspect returns a representation of the object value.
func (c *CubePrimitive) Inspect() string {
	var out bytes.Buffer

	var args []string
	for _, p := range c.Arguments {
		args = append(args, p.Inspect())
	}

	out.WriteString("cube(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(");\n")

	return out.String()
}

// Type returns the type of the object.
func (c *CubePrimitive) Type() T { return CubePrimitiveT }

// GroupBlockPrimitive represents an object of that type.
type GroupBlockPrimitive struct {
	Body Object
}

// Inspect returns a representation of the object value.
func (g *GroupBlockPrimitive) Inspect() string {
	var out bytes.Buffer

	out.WriteString("group() {\n")
	out.WriteString(g.Body.Inspect())
	out.WriteString("\n}")

	return out.String()
}

// Type returns the type of the object.
func (g *GroupBlockPrimitive) Type() T { return GroupBlockPrimitiveT }

// NamedArgument represents an object of that type.
type NamedArgument struct {
	Name  string
	Value Object
}

// Inspect returns a representation of the object value.
func (n *NamedArgument) Inspect() string {
	return fmt.Sprintf("%v = %v", n.Name, n.Value.Inspect())
}

// Type returns the type of the object.
func (n *NamedArgument) Type() T { return NamedArgumentT }

// PolygonPrimitive represents an object of that type.
type PolygonPrimitive struct {
	Arguments []Object
}

// Inspect returns a representation of the object value.
func (c *PolygonPrimitive) Inspect() string {
	var out bytes.Buffer

	var args []string
	for _, p := range c.Arguments {
		args = append(args, p.Inspect())
	}

	out.WriteString("polygon(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(");\n")

	return out.String()
}

// Type returns the type of the object.
func (c *PolygonPrimitive) Type() T { return PolygonPrimitiveT }
