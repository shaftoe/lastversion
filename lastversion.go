package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shaftoe/godevsum"
)

const version = "0.2.1"

type gitRemote struct {
	url, regexpPrefix string
}

var projects = map[string]gitRemote{
	"chef":           {"git://github.com/chef/chef", "refs/tags/v"},
	"consul":         {"git://github.com/hashicorp/consul", "refs/tags/v"},
	"cpython":        {"git://github.com/python/cpython", "refs/tags/v"},
	"docker":         {"git://github.com/docker/docker", "refs/tags/v"},
	"fabric":         {"git://github.com/fabric/fabric", "refs/tags/"},
	"git":            {"git://github.com/git/git", "refs/tags/v"},
	"go":             {"git://github.com/golang/go", "refs/tags/go"},
	"home-assistant": {"git://github.com/home-assistant/home-assistant", "refs/tags/"},
	"kubernetes":     {"git://github.com/kubernetes/kubernetes", "refs/tags/v"},
	"lastversion":    {"git://github.com/shaftoe/lastversion", "refs/tags/v"},
	"packer":         {"git://github.com/mitchellh/packer", "refs/tags/v"},
	"prometheus":     {"git://github.com/prometheus/prometheus", "refs/tags/v"},
	"puppet":         {"git://github.com/puppetlabs/puppet", "refs/tags/"},
	"react":          {"git://github.com/facebook/react", "refs/tags/v"},
	"salt":           {"git://github.com/saltstack/salt", "refs/tags/v"},
	"sslexpired":     {"git://github.com/shaftoe/sslexpired", "refs/tags/"},
	"terraform":      {"git://github.com/hashicorp/terraform.git", "refs/tags/v"},
	"vault":          {"git://github.com/hashicorp/vault", "refs/tags/v"},
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
	if !ok || proj == "" {
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
	gf, err := godevsum.NewGitFetcher("/action/git", true)
	if err != nil {
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

	result, err := godevsum.LatestTaggedVersion(url, remote.regexpPrefix, gf)
	if err != nil {
		msg["err"] = err.Error()
	}
	msg["result"] = result
}
