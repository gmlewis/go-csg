//Linear Extrude with Scale as an interpolated function
// This module does not need to be modified,
// - unless default parameters want to be changed
// - or additional parameters want to be forwarded (e.g. slices,...)
module linear_extrude_fs(height=1,isteps=20,twist=0){
 //union of piecewise generated extrudes
 union(){
   for(i = [ 0: 1: isteps-1]){
     //each new piece needs to be adjusted for height
     translate([0,0,i*height/isteps])
      linear_extrude(
       height=height/isteps,
       twist=twist/isteps,
       scale=f_lefs((i+1)/isteps)/f_lefs(i/isteps)
      )
       // if a twist constant is defined it is split into pieces
       rotate([0,0,-(i/isteps)*twist])
        // each new piece starts where the last ended
        scale(f_lefs(i/isteps))
         obj2D_lefs();
   }
 }
}
// This function defines the scale function
// - Function name must not be modified
// - Modify the contents/return value to define the function
function f_lefs(x) =
 let(span=150,start=20,normpos=45)
 sin(x*span+start)/sin(normpos);
// This module defines the base 2D object to be extruded
// - Function name must not be modified
// - Modify the contents to define the base 2D object
module obj2D_lefs(){
 translate([-4,-3])
  square([9,12]);
}
//Top rendered object demonstrating the interpolation steps
translate([0,0,25])
linear_extrude_fs(height=20,isteps=4);
linear_extrude_fs(height=20);
//Bottom rendered object demonstrating the inclusion of a twist
translate([0,0,-25])
linear_extrude_fs(height=20,twist=90,isteps=30);
