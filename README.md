# go-csg

[![GoDoc](https://godoc.org/github.com/gmlewis/go-csg?status.svg)](https://godoc.org/github.com/gmlewis/go-csg)
[![Test Status](https://github.com/gmlewis/go-csg/workflows/tests/badge.svg)](https://github.com/gmlewis/go-csg/actions?query=workflow%3Atests)

An experimental CSG interpreter and converter written in Go.

CSG is an intermediate file format that is output by CAD tools
such as OpenSCAD.
See: https://github.com/openscad/openscad/wiki/CSG-File-Format
for more information.

This project is based on the "Monkey language" from the book
`Writing An Interpreter in Go` by Thorsten Ball:
https://interpreterbook.com/

The goal of this project is to parse all valid CSG files
and convert them to [IRMF](http://irmf.io) files
in order to make IRMF accessible to more people.

## Dependencies

- [Go](https://golang.org) version 1.13.4 or later

## Running an example

You can use the `run.sh` bash script and the example number found in
the [examples directory](/examples) to convert a CSG example file
to IRMF.

```sh
$ ./run.sh 11
(ctrl-c to quit)
```

## CSG Supported Features:

- [x] circle
- [x] cube
- [x] cylinder
- [x] difference
- [ ] hull
- [x] intersection
- [x] linear_extrude
- [ ] minkowski
- [x] multmatrix
- [x] polygon
- [ ] polyhedron
- [ ] projection
- [x] rotate_extrude
- [x] sphere
- [x] square
- [ ] text
- [x] union

PRs are welcome! :smile:

---

Enjoy!

---

# License

Copyright 2019 Glenn M. Lewis. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
