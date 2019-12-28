# go-openscad

An experimental OpenSCAD interpreter and converter.

This project is based on the "Monkey language" from the book
`Writing An Interpreter in Go` by Thorsten Ball:
https://interpreterbook.com/

The goal of this project is to parse all valid OpenSCAD
programs and convert them to [IRMF](http://irmf.io) files
in the hope of making IRMF more accessible to more people.

## Dependencies

- [Go](https://golang.org) version 1.13.4 or later

## Running an example

You can use the `run.sh` bash script and the example number found in
the [examples directory](/examples) to convert an OpenSCAD example
to IRMF.

```sh
$ ./run.sh 11
(ctrl-c to quit)
```

----------------------------------------------------------------------

Enjoy!

----------------------------------------------------------------------

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
