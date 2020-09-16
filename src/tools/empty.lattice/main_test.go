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
	assert.Equal(t, "# 1 S0000\n1 0 1.23 <eps>\n", makeEmptyLattice(1.23))
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
}
