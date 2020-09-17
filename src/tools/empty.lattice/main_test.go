package main

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	assert.Equal(t, "1.23", round(1.23))
	assert.Equal(t, "1.23", round(1.2344546))
	assert.Equal(t, "1.24", round(1.236768))
	assert.Equal(t, "1.20", round(1.20))
	assert.Equal(t, "1.00", round(1.001212))
	assert.Equal(t, "100.00", round(100.001212))
	assert.Equal(t, "0.00", round(0.00))
}

func TestMakeEmptyLattice(t *testing.T) {
	params := &params{len: 1.23, segmentName: "SN", silenceWord: "<ttt>"}
	assert.Equal(t, "# 1 SN\n1 0 1.23 <ttt>\n", makeEmptyLattice(params))
}

func TestParseParams(t *testing.T) {
	params := &params{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	takeParams(fs, params)
	err := fs.Parse([]string{"-l", "10"})
	assert.Nil(t, err)
	assert.InDelta(t, 10, params.len, 0.0001)
	err = fs.Parse([]string{"-l", "10.123"})
	assert.Nil(t, err)
	assert.InDelta(t, 10.123, params.len, 0.0001)
	err = fs.Parse([]string{"-l", "500"})
	assert.Nil(t, err)
	assert.InDelta(t, 500, params.len, 0.0001)

	fs = flag.NewFlagSet("", flag.ContinueOnError)
	takeParams(fs, params)
	err = fs.Parse([]string{"-l", "50", "-s", "<ttt>", "-sn", "TT"})
	assert.Nil(t, err)
	assert.InDelta(t, 50, params.len, 0.0001)
	assert.Equal(t, "<ttt>", params.silenceWord)
	assert.Equal(t, "TT", params.segmentName)
}
