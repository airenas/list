package main

import (
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

func TestSplitText(t *testing.T) {
	words := makeTestWords(testWord(40, "", 0, 1), testWord(24, "", 1.0, 2))
	res := splitText(words)
	assert.Equal(t, [][]int{{0, 1}, {1, 2}}, res)
}

func TestSplitText_Several(t *testing.T) {
	words := makeTestWords(testWord(20, "", 0, 1), testWord(20, "", 1, 2),
		testWord(20, "", 2, 3), testWord(20, "", 3, 4))
	res := splitText(words)
	assert.Equal(t, [][]int{{0, 2}, {2, 4}}, res)
}

func TestSplitText_Dot(t *testing.T) {
	words := makeTestWords(testWord(35, ".", 0, 1), testWord(5, ",", 1, 2),
		testWord(20, "", 2, 3), testWord(10, "", 3, 4))
	res := splitText(words)
	assert.Equal(t, [][]int{{0, 1}, {1, 4}}, res)
}

func TestSplitText_Dot2(t *testing.T) {
	words := makeTestWords(testWord(35, ",", 0, 1), testWord(5, ".", 1, 2),
		testWord(20, "", 2, 3), testWord(10, "", 3, 4))
	res := splitText(words)
	assert.Equal(t, [][]int{{0, 2}, {2, 4}}, res)
}

func TestSplitText_Pause(t *testing.T) {
	words := makeTestWords(testWord(30, ",", 0, 1), testWord(5, ",", 1, 2),
		testWord(20, "", 2.2, 3), testWord(10, "", 3, 4))
	res := splitText(words)
	assert.Equal(t, [][]int{{0, 2}, {2, 4}}, res)
}

func TestSplitText_SeveralSplit(t *testing.T) {
	words := makeTestWords(testWord(20, "", 0, 1), testWord(20, "", 1, 2),
		testWord(20, "", 2, 3), testWord(20, "", 3, 4),
		testWord(20, "", 4, 5), testWord(20, "", 5, 6),
		testWord(20, "", 6, 7), testWord(20, "", 7, 8))
	res := splitText(words)
	assert.Equal(t, [][]int{{0, 2}, {2, 4}, {4, 6}, {6, 8}}, res)
}

func TestSplitText_SeveralSplit2(t *testing.T) {
	words := makeTestWords(testWord(10, "", 0, 1), testWord(10, "", 1, 2),
		testWord(20, "", 2, 3), testWord(10, "", 3, 4),
		testWord(20, "", 4, 5), testWord(10, "", 5, 6),
		testWord(20, "", 6.1, 7), testWord(10, "", 7, 8))
	res := splitText(words)
	assert.Equal(t, [][]int{{0, 3}, {3, 6}, {6, 8}}, res)
}

func makeTestWords(words ...*lattice.Word) []*lattice.Word {
	res := make([]*lattice.Word, 0)
	res = append(res, words...)
	return res
}

func testWord(wc int, sep string, from, to float64) *lattice.Word {
	res := &lattice.Word{Punct: sep, From: fmt.Sprintf("%.2f", from), To: fmt.Sprintf("%.2f", to),
		Words: []string{strings.Repeat("a", wc)}}
	return res
}
