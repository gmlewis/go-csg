package parser

import (
	"fmt"
	"testing"

	"github.com/gmlewis/go-csg/lexer"
)

func TestCSGPrimitives(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"circle($fn = 50, $fa = 12, $fs = 2, r = 0.75);", "circle($fn = 50, $fa = 12, $fs = 2, r = 0.75)"},
		{"cube(size = [10.4, 9.02857, 20.8], center = false);", "cube(size = [10.4, 9.02857, 20.8], center = false)"},
		{"cylinder($fn = 0, $fa = 12, $fs = 2, h = 8, r1 = 6, r2 = 6, center = true);", "cylinder($fn = 0, $fa = 12, $fs = 2, h = 8, r1 = 6, r2 = 6, center = true)"},
		{"group();", "group()"},
		{"polygon(points = [[0, 0], [20, 0], [0, 40]], paths = undef, convexity = 1);", "polygon(points = [[0, 0], [20, 0], [0, 40]], paths = undef, convexity = 1)"},
		{"polygon(points = [[0, 0], [40, 0], [0, 40]], faces = [[1, 2, 3], [2, 3, 4]], convexity = 1);", "polygon(points = [[0, 0], [40, 0], [0, 40]], faces = [[1, 2, 3], [2, 3, 4]], convexity = 1)"},
		{"sphere($fn = 0, $fa = 12, $fs = 2, r = 5);", "sphere($fn = 0, $fa = 12, $fs = 2, r = 5)"},
		{"square(size = [12, 9], center = false);", "square(size = [12, 9], center = false)"},
		{`text(text = "HeartyGFX", size = 3, spacing = 1, font = "Arial Black", direction = "ltr", language = "en", script = "Latn", halign = "left", valign = "center", $fn = 0, $fa = 12, $fs = 2);`, `text(text = "HeartyGFX", size = 3, spacing = 1, font = "Arial Black", direction = "ltr", language = "en", script = "Latn", halign = "left", valign = "center", $fn = 0, $fa = 12, $fs = 2)`},

		// Group primitives:
		{"color([0, 0, 0.545098, 1]) { sphere(); }", "color([0, 0, 0.545098, 1]) { sphere() }"},
		{"difference() { sphere(); }", "difference() { sphere() }"},
		{"group() { sphere(); }", "group() { sphere() }"},
		{"hull() { sphere(); }", "hull() { sphere() }"},
		{"intersection() { sphere(); }", "intersection() { sphere() }"},
		{"linear_extrude(height = 0.666667, center = false, convexity = 1, twist = 3, slices = 2, scale = [0.670925, 0.670925], $fn = 0, $fa = 12, $fs = 2) { sphere(); }", "linear_extrude(height = 0.666667, center = false, convexity = 1, twist = 3, slices = 2, scale = [0.670925, 0.670925], $fn = 0, $fa = 12, $fs = 2) { sphere() }"},
		{"minkowski(convexity = 0) { sphere(); }", "minkowski(convexity = 0) { sphere() }"},
		{"multmatrix([[0.0193379, 0.999813, 0, 0], [-0.999813, 0.0193379, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) { sphere(); }", "multmatrix([[0.0193379, 0.999813, 0, 0], [(-0.999813), 0.0193379, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]]) { sphere() }"},
		{"projection(cut = false, convexity = 0) { sphere(); }", "projection(cut = false, convexity = 0) { sphere() }"},
		{"rotate_extrude(convexity = 2, $fn = 100, $fa = 12, $fs = 2) { sphere(); }", "rotate_extrude(convexity = 2, $fn = 100, $fa = 12, $fs = 2) { sphere() }"},
		{"union() { sphere(); cube(); }", "union() { sphere(); cube() }"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			le := lexer.New(tt.input)
			p := New(le)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			got := program.String()
			if got != tt.want {
				t.Errorf("string = %v, want %v", got, tt.want)
			}
		})
	}
}
