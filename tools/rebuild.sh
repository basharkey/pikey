#!/bin/bash
go build type-byte.go
sudo cp type-byte /usr/local/bin
sudo go run destroy.go
sudo go run initialize.go
