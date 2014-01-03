#!/bin/bash -e

cd $(dirname "$0")

function testCommon
{
    # Code formatting
    if [ "" != "$(go fmt)" ]; then
        echo "Code is not formatted correctly. Please correct and commit"
    fi

    # Code vetting
    go vet

}

function testMain
{
    testCommon

    go build -race
}

function testLibrary
{
    library="$1"

    pushd "$library"
    echo "Test $library"

    testCommon

    # Race detection
    go test -race

    # Code coverage
    go test -coverprofile cover.out > /dev/null
    if [ -e cover.out ]; then
        if go tool cover -func=cover.out | grep -v "100.0%"; then
            echo "Unit test code coverage failed"
            exit 1
        fi
    fi

    popd
}

testLibrary debug
testLibrary testing_utils
testLibrary connection
testLibrary packet
testLibrary pipe
testLibrary parser
echo "Test main"
testMain
