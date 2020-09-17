package main

import (
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

func TestAddInitialSilence(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 10.00 10.12 w9
 1 10.12 10.25 w10
 1 10.25 12.25 <eps>

 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 w12
 `))
	lat, _ = fixTime(lat, &params{len: 12.80})
	testSil(t, lat[0].Words[0], "0.00", "10.00")
}

func TestAddFinalSilence(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 10.00 10.12 w9
 1 10.12 10.25 w10
 1 10.25 12.25 <eps>

 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 w12
 `))
	lat, _ = fixTime(lat, &params{len: 13.80})
	testSil(t, lat[1].Words[2], "12.80", "13.80")
}

func TestAppendSilence(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 10.00 10.12 <eps>
 1 10.12 10.25 w10
 1 10.25 12.25 <eps>

 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 <eps>
 `))
	lat, _ = fixTime(lat, &params{len: 13.80})
	testSil(t, lat[0].Words[0], "0.00", "10.12")
	testSil(t, lat[1].Words[1], "12.50", "13.80")
}

func TestInsertInside(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 10.12 10.25 w10
 
 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 <eps>
 `))
	lat, _ = fixTime(lat, &params{len: 13.80})
	testSil(t, lat[1].Words[0], "10.25", "12.25")
}

func TestAppendInside(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 10.12 10.25 w10
 1 10.25 10.50 <eps>
 
 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 <eps>
 `))
	lat, _ = fixTime(lat, &params{len: 13.80})
	testSil(t, lat[0].Words[2], "10.25", "12.25")
}

func testSil(t *testing.T, w *lattice.Word, from, to string) {
	assert.Equal(t, "<eps>", w.Words[0])
	assert.Equal(t, "1", w.Main)
	assert.Equal(t, from, w.From)
	assert.Equal(t, to, w.To)
}
