package irmf

var primitives = map[string]string{
	"circle": `float circle(in vec3 xyz) {
		// TODO
		return 1.0;
	}
	`,
	"cube": `float cube(in vec3 size, in bool center, in vec3 xyz) {
	xyz /= size;
	if (!center) { xyz -= vec3(0.5); }
	if (any(greaterThan(abs(xyz), vec3(0.5)))) { return 0.0; }
	return 1.0;
}
`,
	"cylinder": `float cylinder(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,
	"polygon": `float polygon(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,
	"polyhedron": `float polyhedron(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,
	"sphere": `float sphere(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,
	"square": `float square(in vec3 xyz) {
	// TODO
	return 1.0;
}
`,
}
