// Rather kludgy module for determining bounding box from intersecting projections
module BoundingBox()
{
	intersection()
	{
		translate([0,0,0])
		linear_extrude(height = 1000, center = true, convexity = 10, twist = 0)
		projection(cut=false) intersection()
		{
			rotate([0,90,0])
			linear_extrude(height = 1000, center = true, convexity = 10, twist = 0)
			projection(cut=false)
			rotate([0,-90,0])
			children(0);

			rotate([90,0,0])
			linear_extrude(height = 1000, center = true, convexity = 10, twist = 0)
			projection(cut=false)
			rotate([-90,0,0])
			children(0);
		}
		rotate([90,0,0])
		linear_extrude(height = 1000, center = true, convexity = 10, twist = 0)
		projection(cut=false)
		rotate([-90,0,0])
		intersection()
		{
			rotate([0,90,0])
			linear_extrude(height = 1000, center = true, convexity = 10, twist = 0)
			projection(cut=false)
			rotate([0,-90,0])
			children(0);

			rotate([0,0,0])
			linear_extrude(height = 1000, center = true, convexity = 10, twist = 0)
			projection(cut=false)
			rotate([0,0,0])
			children(0);
		}
	}
}

// Test module on ellipsoid
translate([0,0,40]) scale([1,2,3]) sphere(r=5);
BoundingBox() scale([1,2,3]) sphere(r=5);
