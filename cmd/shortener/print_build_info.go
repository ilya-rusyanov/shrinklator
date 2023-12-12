package main

import (
	"fmt"
	"io"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo(w io.Writer) {
	ver := set(buildVersion)
	date := set(buildDate)
	commit := set(buildCommit)

	bi := fmt.Sprintf(`Build version: %s
Build date: %s
Build commit: %s`, ver, date, commit)

	w.Write([]byte(bi))
}

func set(inp string) string {
	res := inp

	if res == "" {
		res = "N/A"
	}

	return res
}
