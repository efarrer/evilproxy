#!/bin/bash -e

cd $(dirname "$0")

START_PORT=8080
END_PORT=8081

function startServer
{
    nc -l localhost $END_PORT
}

function startEvilProxy
{
    ./evilproxy --server=:$START_PORT --client=:$END_PORT --connections=1 --debug
}

function startClient
{
    output="$1"
    echo "$output" | nc localhost $START_PORT > /dev/null
}


function testCommon
{
    # Code formatting
    if [ "" != "$(go fmt)" ]; then
        echo "Code is not formatted correctly. Please correct and commit"
    fi

    # Code vetting
    go vet

}

function ensureServerAndProxyShutdownWhenClientQuits
{
    echo "Ensure server and proxy shutdown when client quits."
    startEvilProxy &
    startServer &
    sleep 0.5
    startClient "Hi" &
    wait
}

function testMain
{
    testCommon

    go build -race

    ensureServerAndProxyShutdownWhenClientQuits
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
