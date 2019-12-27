// Package object implements the objects system for the language.
package object

import "fmt"

// T represents the type of the object.
type T string

// Object types.
const (
	BooleanT = "BOOLEAN"
	IntegerT = "INTEGER"
	NullT    = "NULL"
)

// Object represents an object or value type within the language.
type Object interface {
	Type() T
	Inspect() string
}

// Integer represents an object of that type.
type Integer struct {
	Value int64
}

// Inspect returns a representation of the object value.
func (i *Integer) Inspect() string { return fmt.Sprintf("%v", i.Value) }

// Type returns the type of the object.
func (i *Integer) Type() T { return IntegerT }

// Boolean represents an object of that type.
type Boolean struct {
	Value bool
}

// Inspect returns a representation of the object value.
func (i *Boolean) Inspect() string { return fmt.Sprintf("%v", i.Value) }

// Type returns the type of the object.
func (i *Boolean) Type() T { return BooleanT }

// Null represents an object of that type.
type Null struct{}

// Inspect returns a representation of the object value.
func (i *Null) Inspect() string { return "null" }

// Type returns the type of the object.
func (i *Null) Type() T { return NullT }
