package main

import (
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	assert.Equal(t, "WEBVTT\n", getHeader())
}

func TestGetWebVTT(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2
`))
	text := getWebVTT(lat)
	assert.Equal(t, "\n00:00:00.010 --> 00:00:00.020\nw w2\n", text)
}
