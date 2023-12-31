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
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw w2\n", text)
}

func TestGetWebVTT_Underscore(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w_x_y
1 1.02 2.02 w2_a
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw x y w2 a\n", text)
}

func TestGetWebVTT_Several(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2

# 2 S2
1 5.01 6.02 w3
1 6.02 7.02 w4
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw w2\n\n00:05.010 --> 00:07.020\nw3 w4\n", text)
}

func TestGetSkipNonMain(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
0 1.02 2.02 w2
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:01.020\nw\n", text)
}

func TestGetSkipSil(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 2.54 w
1 2.55 2.02 <eps>
1 2.02 3.02 w2
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:03.020\nw w2\n", text)
}

func TestGetChopOnLongSil(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 <eps>
1 2.02 3.02 w2
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:01.020\nw\n\n00:02.020 --> 00:03.020\nw2\n", text)
}

func TestGetMinutes(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 62.02 w2
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 01:02.020\nw w2\n", text)
}

func TestGetHours(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 3602.02 w2
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 01:00:02.020\nw w2\n", text)
}

func TestGetLongHours(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 3600182.52 w2
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 1000:03:02.520\nw w2\n", text)
}

func TestGetPuntc(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 W ,
1 1.02 2.02 w2 .
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nW, w2.\n", text)
}

func TestGetPuntcDash(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w -
1 1.02 2.02 w2 .
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw - w2.\n", text)
}

func Test_takeParams(t *testing.T) {
	type args struct {
		params []string
		env string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
		wantHeader string
	}{
		{name:"No params", args: args{params: []string{}, env: ""}, wantErr: false, wantHeader: ""},
		{name:"Params", args: args{params: []string{"--header", "olia olia"}, env: ""}, wantErr: false, wantHeader: "olia olia"},
		{name:"Env", args: args{params: []string{}, env: "aaa"}, wantErr: false, wantHeader: "aaa"},
		{name:"Params first", args: args{params: []string{"--header", "olia olia"}, env: "aaa"}, wantErr: false, wantHeader: "olia olia"},
		{name:"Fail", args: args{params: []string{"--xxx"}}, wantErr: true, wantHeader: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("WEBVTT_HEADER", tt.args.env)
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			header:= ""
			takeParams(fs, &header)
			err := fs.Parse(tt.args.params)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantHeader, header)
		})
	}
}
