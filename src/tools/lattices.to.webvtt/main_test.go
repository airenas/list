package main

import (
	"flag"
	"os"
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

func TestGetWebVTT(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2
`))
	fd := &fdata{data: lat}
	text := getWebVTT([]*fdata{fd})
	assert.Equal(t, "\n00:00.010 --> 00:02.020\n<v Kalbėtojas_1>w w2\n", text)
}

func TestGetWebVTT_Speaker(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	text := getWebVTT([]*fdata{fd})
	assert.Equal(t, "\n00:00.010 --> 00:02.020\n<v sp>w w2\n", text)
}

func TestGetWebVTT_Several(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2
1 10.02 20.02 w3
`))
	lat2, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 2.01 3.02 w
1 3.02 4.02 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	fd2 := &fdata{data: lat2, speaker: "sp2"}
	text := getWebVTT([]*fdata{fd, fd2})
	assert.Equal(t, "\n00:00.010 --> 00:02.020\n<v sp>w w2\n\n00:02.010 --> 00:04.020\n<v sp2>w w2\n\n00:10.020 --> 00:20.020\n<v sp>w3\n", text)
}

func Test_takeParams(t *testing.T) {
	type args struct {
		params []string
		env    string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantHeader string
		wantMap    string
	}{
		{name: "No params", args: args{params: []string{}, env: ""}, wantErr: false, wantHeader: ""},
		{name: "Params", args: args{params: []string{"--header", "olia olia", "--namesMap", "file"}, env: ""}, wantErr: false,
			wantHeader: "olia olia", wantMap: "file"},
		{name: "Env", args: args{params: []string{}, env: "aaa"}, wantErr: false, wantHeader: "aaa"},
		{name: "Params first", args: args{params: []string{"--header", "olia olia"}, env: "aaa"}, wantErr: false, wantHeader: "olia olia"},
		{name: "Fail", args: args{params: []string{"--xxx"}}, wantErr: true, wantHeader: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("WEBVTT_HEADER", tt.args.env)
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			fnMap := ""
			header := ""
			takeParams(fs, &fnMap, &header)
			err := fs.Parse(tt.args.params)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantHeader, header)
			assert.Equal(t, tt.wantMap, fnMap)
		})
	}
}
