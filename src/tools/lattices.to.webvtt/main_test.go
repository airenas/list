package main

import (
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
	assert.Equal(t, "\n00:00.010 --> 00:02.020\n<v Kalbėtojas 1>w w2\n", text)
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
