//Linear Extrude with Twist as an interpolated function
// This module does not need to be modified,
// - unless default parameters want to be changed
// - or additional parameters want to be forwarded (e.g. slices,...)
module linear_extrude_ft(height=1,isteps=20,scale=1){
  //union of piecewise generated extrudes
  union(){
    for(i = [ 0: 1: isteps-1]){
      //each new piece needs to be adjusted for height
      translate([0,0,i*height/isteps])
       linear_extrude(
        height=height/isteps,
        twist=f_left((i+1)/isteps)-f_left((i)/isteps),
        scale=(1-(1-scale)*(i+1)/isteps)/(1-(1-scale)*i/isteps)
       )
        //Rotate to next start point
        rotate([0,0,-f_left(i/isteps)])
         //Scale to end of last piece size
         scale(1-(1-scale)*(i/isteps))
          obj2D_left();
    }
  }
}
// This function defines the twist function
// - Function name must not be modified
// - Modify the contents/return value to define the function
function f_left(x) =
  let(twist=90,span=180,start=0)
  twist*sin(x*span+start);
// This module defines the base 2D object to be extruded
// - Function name must not be modified
// - Modify the contents to define the base 2D object
module obj2D_left(){
  translate([-4,-3])
   square([12,9]);
}
//Left rendered object demonstrating the interpolation steps
translate([-20,0])
linear_extrude_ft(height=30,isteps=5);
linear_extrude_ft(height=30);
//Right rendered object demonstrating the scale inclusion
translate([25,0])
linear_extrude_ft(height=30,scale=3);
