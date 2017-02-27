# Lastversion: an HTTP service to fetch last stable version of OpenSource projects
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/shaftoe/lastversion)](https://goreportcard.com/report/github.com/shaftoe/lastversion)
[![Issue Count](https://codeclimate.com/github/shaftoe/lastversion/badges/issue_count.svg)](https://codeclimate.com/github/shaftoe/lastversion)

_Lastversion_ is a small [serverless][6] project that I built [to teach myself][7] a bit of [Go][1] and [OpenWhisk][2] along the way. It's making use of [OpenWhisk built-in Docker support][4] to run the Go application (lastversion.go) in a container, which is forking the provided (statically linked) `git` binary to fetch tags from public Git repositories, then selects the last stable version and returns it in json format

## Usage

    $ curl lastversion.info/?project=go
    {
        "request": "go",
        "result": "1.8"
    }
    $ curl lastversion.info/?project=kubernetes
    {
        "request": "kubernetes",
        "result": "1.5.3"
    }
    $ curl lastversion.info/?project=skynet
    {
        "err": "project not available",
        "request": "skynet"
    }

## Current limitations

- Lastversion supports only a very limited list of projects, but if you're interested in having more supported, just send (me a messsage|a pull request) and I'll be happy to add. One obvious next step could be to let the client provide the (GitHub) url and regexp prefix needed to devise the last stable version to make it more general
- the docker image running the app has no git binary available, so we ship a statically linked git binary which has some limitations, so for example `git://` is the only transport protocol supported
- the [OpenWhisk API gateway][3] support is still experimental and hence very limited in features, so it's not possible to support a url scheme like `http://lastversion.info/docker` or similar, nor is possible to return anything other then `Content-Type: application/json` as response

## Develop on OpenWhisk

Fetch the development environment installing [OpenWhisk development Vagrant box][5] and set up credentials for `wsk` tool as suggested

### Deploy lastversion action and api gateway

    $ ./build.sh create
    # ... edit code ...
    $ ./build.sh update

### Invoke lastversion action

    $ ./build.sh run <project_name>

or fetch the api URL and use an http client

    $ wsk -i api-experimental list /lastversion
    $ curl <https url>/?project=<project_name>

### Destroy lastversion

    $ ./build.sh delete

### Git binary

Compiled from sources (v2.11.1) with those flags: `$ make "CFLAGS=${CFLAGS} -static" NO_OPENSSL=1 NO_CURL=1`

[1]: https://golang.org/ "Go"
[2]: http://openwhisk.org/ "OpenWhisk"
[3]: https://github.com/openwhisk/openwhisk/blob/master/docs/apigateway.md "API gateway"
[4]: https://www.ibm.com/blogs/bluemix/2017/01/docker-bluemix-openwhisk/ "Docker support"
[5]: https://github.com/openwhisk/openwhisk#quick-start "OpenWhisk devel quick start"
[6]: https://en.wikipedia.org/wiki/Serverless_computing
[7]: https://github.com/shaftoe/godevsum