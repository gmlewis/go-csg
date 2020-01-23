difference()
{
  rotate_extrude($fn=200) polygon( points=[[0,0],[50,0],[47,3],[6,6],[6,66],[99,159],[96,159],[0,69]] );
  scale([-1,1,1]) text("Glenn", halign="center", valign="center", size=27);
}
