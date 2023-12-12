package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintBuildInfo(t *testing.T) {
	cases := []struct {
		name   string
		ver    string
		date   string
		commit string
		want   string
	}{
		{
			name: "nothing set",
			want: `Build version: N/A
Build date: N/A
Build commit: N/A`,
		},
		{
			name: "build version set",
			ver:  "1.0.0",
			want: `Build version: 1.0.0
Build date: N/A
Build commit: N/A`,
		},
		{
			name: "build date set",
			date: "30.11.2023",
			want: `Build version: N/A
Build date: 30.11.2023
Build commit: N/A`,
		},
		{
			name:   "build commit set",
			commit: "aabbcc",
			want: `Build version: N/A
Build date: N/A
Build commit: aabbcc`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buildVersion = tc.ver
			buildDate = tc.date
			buildCommit = tc.commit

			sb := strings.Builder{}
			printBuildInfo(&sb)

			assert.Equal(t, tc.want, sb.String())
		})
	}
}
