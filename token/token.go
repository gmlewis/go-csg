// Package token tokenizes the input text.
package token

// T represents a token type.
type T string

// Token represents a token.
type Token struct {
	Type    T
	Literal string
}

// These constants represent the various types of possible tokens.
const (
	// Special types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	FLOAT  = "FLOAT"  // 1343456.0 or 5.0
	STRING = "STRING" // "abcdef"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!" // also show only modifier
	ASTERISK = "*" // also disable modifier
	SLASH    = "/"
	// DOLLAR   = "$"
	QUESTION = "?"

	LT = "<"
	GT = ">"

	EQ    = "=="
	NOTEQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Comments
	LINECOMMENT = "//"

	// Modifiers
	POUND   = "#" // highlight / debug
	PERCENT = "%" // transparent / background

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// https://www.openscad.org/cheatsheet/index.html

	// // Syntax
	// CONDITIONAL = "CONDITIONAL"
	// MODULE      = "MODULE"
	// INCLUDE     = "INCLUDE"
	// USE         = "USE"

	// Constants
	UNDEF = "UNDEF"
	// PI    = "PI"

	// 2D
	CIRCLE  = "CIRCLE"
	SQUARE  = "SQUARE"
	POLYGON = "POLYGON"
	TEXT    = "TEXT"
	// IMPORT     = "IMPORT"
	PROJECTION = "PROJECTION"

	// 3D
	SPHERE   = "SPHERE"
	CUBE     = "CUBE"
	CYLINDER = "CYLINDER"
	// POLYHEDRON     = "POLYHEDRON"
	LINEAR_EXTRUDE = "LINEAR_EXTRUDE"
	ROTATE_EXTRUDE = "ROTATE_EXTRUDE"
	// SURFACE        = "SURFACE"

	// Transformations
	// TRANSLATE  = "TRANSLATE"
	// ROTATE     = "ROTATE"
	// SCALE      = "SCALE"
	// RESIZE     = "RESIZE"
	// MIRROR     = "MIRROR"
	MULTMATRIX = "MULTMATRIX"
	COLOR      = "COLOR"
	OFFSET     = "OFFSET"
	HULL       = "HULL"
	MINKOWSKI  = "MINKOWSKI"

	// Boolean Operations
	UNION        = "UNION"
	DIFFERENCE   = "DIFFERENCE"
	INTERSECTION = "INTERSECTION"
	GROUP        = "GROUP" // Undocumented

	// Flow Control
	// FOR              = "FOR"
	// INTERSECTION_FOR = "INTERSECTION_FOR"

	// Type test functions
	// IS_UNDEF  = "IS_UNDEF"
	// IS_BOOL   = "IS_BOOL"
	// IS_NUM    = "IS_NUM"
	// IS_STRING = "IS_STRING"
	// IS_LIST   = "IS_LIST"

	// Other
	// ECHO     = "ECHO"
	// RENDER   = "RENDER"
	// CHILDREN = "CHILDREN"
	// ASSERT   = "ASSERT"

	// Functions
	// CONCAT        = "CONCAT"
	// LOOKUP        = "LOOKUP"
	// STR           = "STR"
	// CHR           = "CHR"
	// ORD           = "ORD"
	// SEARCH        = "SEARCH"
	// VERSION       = "VERSION"
	// VERSION_NUM   = "VERSION_NUM"
	// PARENT_MODULE = "PARENT_MODULE"

	// Mathematical
	// ABS   = "ABS"
	// SIGN  = "SIGN"
	// SIN   = "SIN"
	// COS   = "COS"
	// TAN   = "TAN"
	// ACOS  = "ACOS"
	// ASIN  = "ASIN"
	// ATAN  = "ATAN"
	// ATAN2 = "ATAN2"
	// FLOOR = "FLOOR"
	// ROUND = "ROUND"
	// CEIL  = "CEIL"
	// LN    = "LN"
	// LEN = "LEN"
	// LET = "LET"
	// LOG   = "LOG"
	// POW   = "POW"
	// SQRT  = "SQRT"
	// EXP   = "EXP"
	// RANDS = "RANDS"
	// MIN   = "MIN"
	// MAX   = "MAX"
	// NORM  = "NORM"
	// CROSS = "CROSS"
)

var keywords = map[string]T{
	"function": FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,

	// https://www.openscad.org/cheatsheet/index.html

	// Syntax
	// "cond":    CONDITIONAL,
	// "module":  MODULE,
	// "include": INCLUDE,
	// "use":     USE,

	// Constants
	"undef": UNDEF,
	// "PI":    PI,

	// 2D
	"circle":  CIRCLE,
	"square":  SQUARE,
	"polygon": POLYGON,
	// "text":    TEXT,
	// "import":     IMPORT,
	"projection": PROJECTION,

	// 3D
	"sphere":   SPHERE,
	"cube":     CUBE,
	"cylinder": CYLINDER,
	// "polyhedron":     POLYHEDRON,
	"linear_extrude": LINEAR_EXTRUDE,
	"rotate_extrude": ROTATE_EXTRUDE,
	// "surface":        SURFACE,

	// Transformations
	// "translate":  TRANSLATE,
	// "rotate":     ROTATE,
	// "scale":      SCALE,
	// "resize":     RESIZE,
	// "mirror":     MIRROR,
	"multmatrix": MULTMATRIX,
	"color":      COLOR,
	"offset":     OFFSET,
	"hull":       HULL,
	"minkowski":  MINKOWSKI,

	// Boolean Operations
	"union":        UNION,
	"difference":   DIFFERENCE,
	"intersection": INTERSECTION,
	"group":        GROUP,

	// Flow Control
	// "for":              FOR,
	// "intersection_for": INTERSECTION_FOR,

	// Type test functions
	// "is_undef":  IS_UNDEF,
	// "is_bool":   IS_BOOL,
	// "is_num":    IS_NUM,
	// "is_string": IS_STRING,
	// "is_list":   IS_LIST,

	// Other
	// "echo":     ECHO,
	// "render":   RENDER,
	// "children": CHILDREN,
	// "assert":   ASSERT,

	// Functions
	// "concat":        CONCAT,
	// "lookup":        LOOKUP,
	// "str":           STR,
	// "chr":           CHR,
	// "ord":           ORD,
	// "search":        SEARCH,
	// "version":       VERSION,
	// "version_num":   VERSION_NUM,
	// "parent_module": PARENT_MODULE,

	// Mathematical
	// "abs":   ABS,
	// "sign":  SIGN,
	// "sin":   SIN,
	// "cos":   COS,
	// "tan":   TAN,
	// "acos":  ACOS,
	// "asin":  ASIN,
	// "atan":  ATAN,
	// "atan2": ATAN2,
	// "floor": FLOOR,
	// "round": ROUND,
	// "ceil":  CEIL,
	// "ln":    LN,
	// "len":   LEN,
	// "let":   LET,
	// "log":   LOG,
	// "pow":   POW,
	// "sqrt":  SQRT,
	// "exp":   EXP,
	// "rands": RANDS,
	// "min":   MIN,
	// "max":   MAX,
	// "norm":  NORM,
	// "cross": CROSS,
}

// LookupIdent looks up an identifier and returns the type of token.
func LookupIdent(ident string) T {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
