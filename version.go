package main

import (
	"fmt"
	"runtime/debug"
)

const Version = "0.2.3"

func version() string {
	v := fmt.Sprintf("v%s", Version)
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return v
	}
	for _, kv := range info.Settings {
		if kv.Key == "vcs.revision" {
			v = fmt.Sprintf("%s+%s", v, kv.Value[0:10])
		}
	}

	return v
}
