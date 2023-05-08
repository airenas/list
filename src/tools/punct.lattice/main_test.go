package main

import (
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/punctuation"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/test/mocks"
	"github.com/petergtz/pegomock/v3"
	"github.com/stretchr/testify/assert"
)

var punctuatorMock *mocks.MockPunctuator

func initTest(t *testing.T) {
	mocks.AttachMockToTest(t)
	punctuatorMock = mocks.NewMockPunctuator()
}

func TestAll(t *testing.T) {
	initTest(t)
	pegomock.When(punctuatorMock.Punctuate(pegomock.AnyStringSlice())).ThenReturn(makeTestPResp("w w2 w3 w4", "W, w2 w3 w4."), nil)
	d, err := punctuate(makeTestLattice()[:2], punctuatorMock)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	pr := punctuatorMock.VerifyWasCalled(pegomock.Once()).Punctuate(pegomock.AnyStringSlice()).GetCapturedArguments()
	assert.Equal(t, []string{"w", "w2", "w3", "w4"}, pr)
	assert.Equal(t, "W", d[0].Words[0].Words[0])
	assert.Equal(t, ",", d[0].Words[0].Punct)
	assert.Equal(t, "w4", d[1].Words[1].Text)
	assert.Equal(t, ".", d[1].Words[1].Punct)
}

func TestGetNext(t *testing.T) {
	initTest(t)
	ni := getNextPartIndex(makeTestLattice(), 0)
	assert.Equal(t, 2, ni)
	ni = getNextPartIndex(makeTestLattice(), 1)
	assert.Equal(t, 2, ni)
	ni = getNextPartIndex(makeTestLattice(), 2)
	assert.Equal(t, 4, ni)
}

func TestGetNext_SplitOnPause(t *testing.T) {
	initTest(t)
	tl, _ := lattice.Read(strings.NewReader(
		`# 6 S4
1 10.00 10.12 w9
1 10.12 10.25 w10
1 10.25 12.25 <eps>

# 7 S4
1 12.25 12.50 w11
1 12.50 12.80 w12
`))

	ni := getNextPartIndex(tl, 0)
	assert.Equal(t, 1, ni)
}

func TestGetNext_IgnoreNonMain(t *testing.T) {
	initTest(t)
	tl, _ := lattice.Read(strings.NewReader(
		`# 6 S4
1 10.00 10.12 w9
1 10.12 10.25 w10
1 10.25 12.25 <eps>
0 10.25 12.25 olia

# 7 S4
1 12.25 12.50 w11
1 12.50 12.80 w12
`))

	ni := getNextPartIndex(tl, 0)
	assert.Equal(t, 1, ni)
}

func TestGetNext_SplitOnShortPause(t *testing.T) {
	initTest(t)
	tl, _ := lattice.Read(strings.NewReader(
		`# 6 S4
1 10.00 10.12 w9
1 10.12 10.25 w10
1 10.25 11.25 <eps>

# 7 S4
1 11.25 12.50 w11
1 12.50 12.80 w12
`))

	ni := getNextPartIndex(tl, 0)
	assert.Equal(t, 2, ni)
}

func TestGetNext_IncreaseIndex(t *testing.T) {
	initTest(t)
	tl, _ := lattice.Read(strings.NewReader(
		`# 6 S4
1 1.00 1.12 w9
1 1.12 2 w10

# 7 S4
1 5 6 w11
1 7 8 w12

# 8 S4
1 13 14 w11
1 15 16 w12
`))

	ni := getNextPartIndex(tl, 1)
	assert.Equal(t, 2, ni)
}

func TestGetWords(t *testing.T) {
	initTest(t)
	wrds := getWords(makeTestLattice(), 0, 1)
	assert.Equal(t, []string{"w", "w2"}, wrds)
}

func TestGetWords_Several(t *testing.T) {
	initTest(t)
	wrds := getWords(makeTestLattice(), 0, 2)
	assert.Equal(t, []string{"w", "w2", "w3", "w4"}, wrds)
}

func TestGetWords_IgnoreNonMain(t *testing.T) {
	initTest(t)
	wrds := getWords(makeTestLattice(), 2, 3)
	assert.Equal(t, []string{"word5", "word6"}, wrds)
}

func TestGetWords_IgnoreSil(t *testing.T) {
	initTest(t)
	wrds := getWords(makeTestLattice(), 3, 4)
	assert.Equal(t, []string{"w5", "w6"}, wrds)
}

func TestGetWords_DropBrakets(t *testing.T) {
	initTest(t)
	wrds := getWords(makeTestLattice(), 4, 5)
	assert.Equal(t, []string{"w7", "w8"}, wrds)
}

func TestAddPunctuation(t *testing.T) {
	initTest(t)
	d := makeTestLattice()
	_ = addPunctuatioData(d, 3, 4, makeTestPResp("w5 w6", "W5, w6."))
	assert.Equal(t, "W5", d[3].Words[1].Words[0])
	assert.Equal(t, ",", d[3].Words[1].Punct)
	assert.Equal(t, "w6", d[3].Words[3].Words[0])
	assert.Equal(t, ".", d[3].Words[3].Punct)
}

func TestAddPunctuation_RestoreBracket(t *testing.T) {
	initTest(t)
	d := makeTestLattice()
	_ = addPunctuatioData(d, 4, 5, makeTestPResp("w7 w8", "W7, w8."))
	assert.Equal(t, "<W7>", d[4].Words[0].Words[0])
	assert.Equal(t, ",", d[4].Words[0].Punct)
}

func makeTestLattice() []*lattice.Part {
	res, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w
1 fr1 to2 w2

# 2 S1
1 fr to w3
1 fr1 to2 w4

# 3 S2
1 fr to word5
0 fr to word
1 fr1 to2 word6

# 4 S2
1 fr to <eps>
1 fr to w5
0 fr to word
1 fr1 to2 w6

# 5 S3
1 fr to <w7>
1 fr to w8
0 fr to word
`))
	return res
}

func TestAddPunctuation_Underscore(t *testing.T) {
	initTest(t)
	d, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w_x
1 fr1 to2 w2_x_y

`))
	_ = addPunctuatioData(d, 0, 1, makeTestPResp("w x w2 x y", "w x w2 x y,"))
	assert.Equal(t, []string{"w", "x"}, d[0].Words[0].Words)
	assert.Equal(t, "", d[0].Words[0].Punct)
	assert.Equal(t, []string{"w2", "x", "y"}, d[0].Words[1].Words)
	assert.Equal(t, ",", d[0].Words[1].Punct)
	_ = addPunctuatioData(d, 0, 1, makeTestPResp("w x w2 x y", "w X W2. X, y,"))
	assert.Equal(t, []string{"w", "x"}, d[0].Words[0].Words)
	assert.Equal(t, "", d[0].Words[0].Punct)
	assert.Equal(t, []string{"W2", "x", "y"}, d[0].Words[1].Words)
	assert.Equal(t, ",", d[0].Words[1].Punct)
}

func makeTestPResp(orig, punct string) *punctuation.Response {
	res := &punctuation.Response{}
	res.Original = strings.Split(orig, " ")
	res.Punctuated = strings.Split(punct, " ")
	return res
}
