#!/bin/bash -ex
go run cmd/csg2irmf/main.go examples/$@*/$@*.csg
