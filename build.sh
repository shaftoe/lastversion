#!/bin/bash
# A few wrappers around the OpenWhisk platform.

set -e
readonly action='lastversion'
readonly command="${1:-create}"
readonly api="/${action}"
readonly version="v1"
readonly wskCli="wsk -i" # remove -i if developing on Bluemix platform

cleanup() {
    rm -f action.zip exec
}
trap cleanup ERR EXIT

usage() {
    echo "Usage: $0 [create|update|run|delete]"
}

build() {
    rm -f action.zip exec
    env GOOS=linux GOARCH=amd64 go build -o exec lastversion.go
    zip action.zip exec git
}

createupdate() {
    $wskCli action $command $action action.zip --docker
}

run() {
    $wskCli action invoke $action --blocking --result --param project $2
}

createapigateway() {
    $wskCli api-experimental create $api /$version get $action
}

delete() {
    $wskCli api-experimental delete $api
    $wskCli action delete $action
}

case $1 in
    create) set -x; build; createupdate; createapigateway;;
    update) set -x; build; createupdate;;
    run) set -x; run $*;;
    delete) set -x; delete;;
    *) usage; exit 1;;
esac
