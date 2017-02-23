package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shaftoe/godevsum"
)

const version = "0.1.0"

type gitRemote struct {
	url, regexpPrefix string
}

var projects = map[string]gitRemote{
	"consul":         gitRemote{"git://github.com/hashicorp/consul", "refs/tags/v"},
	"docker":         gitRemote{"git://github.com/docker/docker", "refs/tags/v"},
	"fabric":         gitRemote{"git://github.com/fabric/fabric", "refs/tags/"},
	"git":            gitRemote{"git://github.com/git/git", "refs/tags/v"},
	"go":             gitRemote{"git://github.com/golang/go", "refs/tags/go"},
	"home-assistant": gitRemote{"git://github.com/home-assistant/home-assistant", "refs/tags/"},
	"kubernetes":     gitRemote{"git://github.com/kubernetes/kubernetes", "refs/tags/v"},
	"prometheus":     gitRemote{"git://github.com/prometheus/prometheus", "refs/tags/v"},
	"terraform":      gitRemote{"git://github.com/hashicorp/terraform.git", "refs/tags/v"},
	"vault":          gitRemote{"git://github.com/hashicorp/vault", "refs/tags/v"},
}

func main() {
	msg := make(map[string]string)
	defer func() {
		resp, _ := json.Marshal(msg)
		fmt.Println(string(resp))
		os.Exit(0)
	}()

	// native actions receive one argument, the JSON object as a string
	arg := os.Args[1]

	// unmarshal the string to a JSON object
	var obj map[string]interface{}
	json.Unmarshal([]byte(arg), &obj)
	proj, ok := obj["project"].(string)
	if !ok {
		msg["err"] = "empty request"
		return
	}
	msg["request"] = proj

	remote, ok := projects[proj]
	if !ok {
		msg["err"] = "project not available"
		return
	}

	// we need to force use of the statically linked git binary present in the
	// OpenWhisk docker container. Binary built with:
	// $ make "CFLAGS=${CFLAGS} -static" NO_OPENSSL=1 NO_CURL=1
	if err := godevsum.SetGitPath("/action/git", true); err != nil {
		msg["err"] = err.Error()
		return
	}

	// building GIT statically has few limitations, for example
	// we can't use neither http transport nor DNS resolution,
	// so we need to replace the host with the IP address in the URL
	url, err := godevsum.ReplaceHostWithIP(remote.url)
	if err != nil {
		msg["err"] = err.Error()
		return
	}

	result, err := godevsum.LatestTaggedVersion(url, remote.regexpPrefix)
	if err != nil {
		msg["err"] = err.Error()
	}
	msg["result"] = result
}
