package main

import (
	"flag"
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

var param *params

func initTest(t *testing.T) {
	param = &params{}
	param.segmentName = "SN"
	param.silenceWord = "tyla"
	param.len = 20
	param.minSilenceLen = 0.1
}
func TestAddInitialSilence(t *testing.T) {
	initTest(t)
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 10.00 10.12 w9
 1 10.12 10.25 w10
 1 10.25 12.25 <eps>

 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 w12
 `))
	lat, _ = fixTime(lat, param)
	testSil(t, lat[0], "0.00", "10.00")
}

func TestNoInitialSilence(t *testing.T) {
	initTest(t)
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00 10.12 <eps>
 1 10.12 10.25 w10
 `))
	lat, _ = fixTime(lat, param)
	assert.Equal(t, "S4", lat[0].Speaker)
}

func TestNoEnd(t *testing.T) {
	initTest(t)
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00 10.12 <eps>
 1 10.12 20.0 w10
 `))
	lat, _ = fixTime(lat, param)
	assert.Equal(t, 1, len(lat))
}

func TestAddFinalSilence(t *testing.T) {
	initTest(t)
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 10.00 10.12 w9
 1 10.12 10.25 w10
 1 10.25 12.25 <eps>

 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 w12
 `))
	lat, _ = fixTime(lat, param)
	testSil(t, lat[3], "12.80", "20.00")
}

func TestInsertSilence(t *testing.T) {
	initTest(t)
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00 10.12 <eps>
 1 10.12 10.25 w10
 1 10.25 11.25 <eps>

 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 <eps>
 `))
	lat, _ = fixTime(lat, param)
	testSil(t, lat[1], "11.25", "12.25")
}

func TestIgnoreNonMain(t *testing.T) {
	initTest(t)
	lat, _ := lattice.Read(strings.NewReader(
		`# 6 S4
 1 0.00  10.25 w10		
 1 10.12 10.25 w10
 1 10.25 10.50 <eps>
 0 10.50 12.25 w10		
 
 # 7 S4
 1 12.25 12.50 w11
 1 12.50 12.80 aaa
 0 12.80 13.80 olia
 `))
	lat, _ = fixTime(lat, param)
	testSil(t, lat[1], "10.50", "12.25")
}

func testSil(t *testing.T, p *lattice.Part, from, to string) {
	assert.Equal(t, param.segmentName, p.Speaker)
	assert.Equal(t, 1, len(p.Words))
	assert.Equal(t, param.silenceWord, p.Words[0].Words[0])
	assert.Equal(t, from, p.Words[0].From)
	assert.Equal(t, to, p.Words[0].To)
}

func TestParseParams(t *testing.T) {
	params := &params{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	takeParams(fs, params)
	err := fs.Parse([]string{"-l", "10.123"})
	assert.Nil(t, err)
	assert.InDelta(t, 10.123, params.len, 0.0001)
	assert.Equal(t, "<tyla>", params.silenceWord)
	assert.Equal(t, "TYLA", params.segmentName)
	fs = flag.NewFlagSet("", flag.ContinueOnError)
	takeParams(fs, params)
	err = fs.Parse([]string{"-l", "50", "-s", "<ttt>", "-sn", "TT", "-ms", "0.2"})
	assert.Nil(t, err)
	assert.InDelta(t, 50, params.len, 0.0001)
	assert.Equal(t, "<ttt>", params.silenceWord)
	assert.Equal(t, "TT", params.segmentName)
	assert.InDelta(t, 0.2, params.minSilenceLen, 0.0001)
}

func TestNumFixed(t *testing.T) {
	initTest(t)
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
	lat, _ = fixTime(lat, param)
	for i := 0; i < len(lat); i++ {

		assert.Equal(t, i+1, lat[i].Num)
	}
}
