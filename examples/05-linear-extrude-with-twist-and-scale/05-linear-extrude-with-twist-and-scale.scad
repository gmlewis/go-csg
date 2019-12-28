//Linear Extrude with Twist and Scale as interpolated functions
// This module does not need to be modified,
// - unless default parameters want to be changed
// - or additional parameters want to be forwarded
module linear_extrude_ftfs(height=1,isteps=20,slices=0){
  //union of piecewise generated extrudes
  union(){
   for(i=[0:1:isteps-1]){
    translate([0,0,i*height/isteps])
     linear_extrude(
      height=height/isteps,
      twist=leftfs_ftw((i+1)/isteps)-leftfs_ftw(i/isteps),
      scale=leftfs_fsc((i+1)/isteps)/leftfs_fsc(i/isteps),
      slices=slices
     )
      rotate([0,0,-leftfs_ftw(i/isteps)])
       scale(leftfs_fsc(i/isteps))
        obj2D_leftfs();
   }
  }
}
// This function defines the scale function
// - Function name must not be modified
// - Modify the contents/return value to define the function
function leftfs_fsc(x)=
  let(scale=3,span=140,start=20)
  scale*sin(x*span+start);
// This function defines the twist function
// - Function name must not be modified
// - Modify the contents/return value to define the function
function leftfs_ftw(x)=
  let(twist=30,span=360,start=0)
  twist*sin(x*span+start);
// This module defines the base 2D object to be extruded
// - Function name must not be modified
// - Modify the contents to define the base 2D object
module obj2D_leftfs(){
   square([12,9]);
}
//Left rendered objects demonstrating the steps effect
translate([0,-50,-60])
rotate([0,0,90])
linear_extrude_ftfs(height=50,isteps=3);

translate([0,-50,0])
linear_extrude_ftfs(height=50,isteps=3);
//Center rendered objects demonstrating the slices effect
translate([0,0,-60])
rotate([0,0,90])
linear_extrude_ftfs(height=50,isteps=3,slices=20);

linear_extrude_ftfs(height=50,isteps=3,slices=20);
//Right rendered objects with default parameters
translate([0,50,-60])
rotate([0,0,90])
linear_extrude_ftfs(height=50);

translate([0,50,0])
linear_extrude_ftfs(height=50);
