//go:generate go-winres make --arch=386,amd64,arm,arm64 --in=../../cli/winresources/docker-proxy/winres.json --out=../../cli/winresources/docker-proxy/resource

package main

import _ "github.com/harness-community/docker-v23/cli/winresources/docker-proxy"
