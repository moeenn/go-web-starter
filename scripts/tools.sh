#! /bin/sh

# installation commands for all tool binaries.
go install -v github.com/a-h/templ/cmd/templ@v0.3.906
go install -v github.com/joho/godotenv/cmd/godotenv@latest
go install -v github.com/pressly/goose/v3/cmd/goose@v3.24.3

# dependency installation command.
go mod tidy
