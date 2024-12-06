#! /bin/sh

# installation commands for all tool binaries goes here
go install github.com/a-h/templ/cmd/templ@latest

# dependency installation commands go here
go mod tidy