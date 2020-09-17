package main

import (
	"flag"
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

func TestNoInitialSilence(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00 10.12 <eps>
 1 10.12 10.25 w10
 `))
	lat, _ = fixTime(lat, &params{len: 12.80})
	testSil(t, lat[0].Words[0], "0.00", "10.12")
}

func TestNoEnd(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00 10.12 <eps>
 1 10.12 10.25 w10
 `))
	lat, _ = fixTime(lat, &params{len: 10.25})
	assert.Equal(t, 2, len(lat[0].Words))
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

func TestIgnoreNonMain(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 0 0.00  10.25 w10		
 1 10.12 10.25 w10
 1 10.25 10.50 <eps>
 0 10.50 12.25 w10		
 
 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 aaa
 0 12.80 13.80 olia
 `))
	lat, _ = fixTime(lat, &params{len: 13.80})
	assert.Equal(t, 5, len(lat[0].Words))
	testSil(t, lat[0].Words[1], "0.00", "10.12")
	assert.Equal(t, 4, len(lat[1].Words))
	testSil(t, lat[1].Words[3], "12.80", "13.80")
}

func testSil(t *testing.T, w *lattice.Word, from, to string) {
	assert.Equal(t, "<eps>", w.Words[0])
	assert.Equal(t, "1", w.Main)
	assert.Equal(t, from, w.From)
	assert.Equal(t, to, w.To)
}

func TestParseParams(t *testing.T) {
	params := &params{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	takeParams(fs, params)
	err := fs.Parse([]string{"-l", "10.123"})
	assert.Nil(t, err)
	assert.InDelta(t, 10.123, params.len, 0.0001)
	assert.Equal(t, "", params.silenceWord)
	fs = flag.NewFlagSet("", flag.ExitOnError)
	takeParams(fs, params)
	err = fs.Parse([]string{"-l", "50", "-s", "<tyla>"})
	assert.Nil(t, err)
	assert.InDelta(t, 50, params.len, 0.0001)
	assert.Equal(t, "<tyla>", params.silenceWord)
}

func TestChangeSilence(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00 10.12 <eps>		
 0 0.00 10.25 w10		
 1 10.12 10.25 w10
 1 10.25 10.75 <eps>
 0 10.50 12.25 w10	
 1 10.75 11.00 <eps>	
 
 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 aaa
 `))
	lat = changeSilenceWord(lat, &params{silenceWord: "<tttt>", minSilenceLen: 0.5})
	assert.Equal(t, "<tttt>", lat[0].Words[0].Words[0])
	assert.Equal(t, "<tttt>", lat[0].Words[3].Words[0])
	assert.Equal(t, "<eps>", lat[0].Words[5].Words[0])
}

func TestChangeSilence_LeaveShort(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00 10.12 <eps>		
 0 0.00 10.25 w10		
 1 10.25 10.75 <eps>
 `))
	lat = changeSilenceWord(lat, &params{silenceWord: "<tttt>", minSilenceLen: 15})
	assert.Equal(t, "<eps>", lat[0].Words[0].Words[0])
}

func TestChangeSilence_LeaveNonMain(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 0 0.00 10.12 <eps>		
 0 0.00 10.25 w10		
 1 10.25 10.75 <eps>
 `))
	lat = changeSilenceWord(lat, &params{silenceWord: "<tttt>", minSilenceLen: 0.5})
	assert.Equal(t, "<eps>", lat[0].Words[0].Words[0])
}
