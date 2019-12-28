// increase the visual detail
$fn = 100;

// the main body :
// a cylinder
rocket_d = 30; 				// 3 cm wide
rocket_r = rocket_d / 2;
rocket_h = 100; 			// 10 cm tall
cylinder(d = rocket_d, h = rocket_h);

// the head :
// a cone
head_d = 40;  				// 4 cm wide
head_r = head_d / 2;
head_h = 40;  				// 4 cm tall
// prepare a triangle
tri_base = head_r;
tri_height = head_h;
tri_points = [[0,			 0],
			  [tri_base,	 0],
			  [0,	tri_height]];
// rotation around X-axis and then 360° around Z-axis
// put it on top of rocket's body
translate([0,0,rocket_h])
rotate_extrude(angle = 360)
	polygon(tri_points);

// the wings :
// 3x triangles
wing_w = 2;					// 2 mm thick
many = 3;					// 3x wings
wing_l = 40;				// length
wing_h = 40;				// height
wing_points = [[0,0],[wing_l,0],[0,wing_h]];

module wing() {
	// let it a bit inside the main body
	in_by = 1;				// 1 mm
	// set it up on the rocket's perimeter
	translate([rocket_r - in_by,0,0])
	// set it upright by rotating around X-axis
	rotate([90,0,0])
	// set some width and center it
	linear_extrude(height = wing_w,center = true)
	// make a triangle
		polygon(wing_points);
}

for (i = [0: many - 1])
	rotate([0, 0, 360 / many * i])
	wing();
