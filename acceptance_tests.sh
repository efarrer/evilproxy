#!/bin/bash -e

cd $(dirname "$0")

cd simulation;

# Code formatting
if [ "" != "$(go fmt)" ]; then
    echo "Code is not formatted correctly. Please correct and commit"
    exit 1
fi

# Code vetting
go vet

# Race detection
go test -race

# Code coverage
go test -coverprofile cover.out > /dev/null
if go tool cover -func=cover.out | grep -v "100.0%"; then
    echo "Unit test code coverage failed"
    exit 1
fi
