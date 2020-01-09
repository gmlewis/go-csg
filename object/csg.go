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

// CylinderPrimitive represents an object of that type.
type CylinderPrimitive struct {
	Arguments []Object
}

// Inspect returns a representation of the object value.
func (c *CylinderPrimitive) Inspect() string {
	var out bytes.Buffer

	var args []string
	for _, p := range c.Arguments {
		args = append(args, p.Inspect())
	}

	out.WriteString("cylinder(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(");\n")

	return out.String()
}

// Type returns the type of the object.
func (c *CylinderPrimitive) Type() T { return CylinderPrimitiveT }

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

// MultmatrixBlockPrimitive represents an object of that type.
type MultmatrixBlockPrimitive struct {
	Arguments []Object
	Body      Object
}

// Inspect returns a representation of the object value.
func (m *MultmatrixBlockPrimitive) Inspect() string {
	var out bytes.Buffer

	var args []string
	for _, p := range m.Arguments {
		args = append(args, p.Inspect())
	}

	out.WriteString("multmatrix(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.Inspect())
	out.WriteString("\n}")

	return out.String()
}

// Type returns the type of the object.
func (m *MultmatrixBlockPrimitive) Type() T { return MultmatrixBlockPrimitiveT }

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

// SpherePrimitive represents an object of that type.
type SpherePrimitive struct {
	Arguments []Object
}

// Inspect returns a representation of the object value.
func (c *SpherePrimitive) Inspect() string {
	var out bytes.Buffer

	var args []string
	for _, p := range c.Arguments {
		args = append(args, p.Inspect())
	}

	out.WriteString("sphere(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(");\n")

	return out.String()
}

// Type returns the type of the object.
func (c *SpherePrimitive) Type() T { return SpherePrimitiveT }

// SquarePrimitive represents an object of that type.
type SquarePrimitive struct {
	Arguments []Object
}

// Inspect returns a representation of the object value.
func (c *SquarePrimitive) Inspect() string {
	var out bytes.Buffer

	var args []string
	for _, p := range c.Arguments {
		args = append(args, p.Inspect())
	}

	out.WriteString("square(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(");\n")

	return out.String()
}

// Type returns the type of the object.
func (c *SquarePrimitive) Type() T { return SquarePrimitiveT }
