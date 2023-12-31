package main

import (
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	l, err := parseLine("olia 1 20 10 a b c d")
	assert.Nil(t, err)
	assert.Equal(t, strings.Split("olia 1 20 10 a b c d", " "), l.fields)
	assert.Equal(t, 20, l.from)
	assert.Equal(t, 10, l.len)
	assert.Equal(t, []string{"a", "b", "c"}, l.rFields)
}

func TestFails(t *testing.T) {
	_, err := parseLine("olia 1 20 ")
	assert.NotNil(t, err)
	_, err = parseLine("olia 1 20a 10 a b c")
	assert.NotNil(t, err)
	_, err = parseLine("olia 1 20 c10 a b c")
	assert.NotNil(t, err)
	_, err = parseLine("olia 1 20 10 a b c")
	assert.NotNil(t, err)
}
func TestParse_Speaker(t *testing.T) {
	l, err := parseLine("olia 1 20 10 a b c d")
	assert.Nil(t, err)
	assert.Equal(t, strings.Split("olia 1 20 10 a b c d", " "), l.fields)
	assert.Equal(t, 20, l.from)
	assert.Equal(t, 10, l.len)
	assert.Equal(t, []string{"a", "b", "c"}, l.rFields)
	assert.Equal(t, "d", l.speaker)
}

func TestToString(t *testing.T) {
	l, err := parseLine("olia 1 20 10 a b c d")
	assert.Nil(t, err)
	assert.Equal(t, "olia 1 20 10 a b c d", toStr(l))
	l.from = 200
	assert.Equal(t, "olia 1 200 10 a b c d", toStr(l))
	l.len = 100
	assert.Equal(t, "olia 1 200 100 a b c d", toStr(l))
}

func TestJoin_Leave(t *testing.T) {
	l1, _ := parseLine("olia 1 0 10 a b c d")
	l2, _ := parseLine("olia 1 10 10 a b c d")
	l3, _ := parseLine("olia 1 20 10 a b c d")
	lns := joinLines([]*line{l1, l2, l3}, 50)
	assert.Equal(t, 3, len(lns))
}

func TestJoin_LeaveOnLimit(t *testing.T) {
	l1, _ := parseLine("olia 1 0 5 a b c d")
	l2, _ := parseLine("olia 1 5 5 a b c d")
	l3, _ := parseLine("olia 1 10 5 a b c d")
	lns := joinLines([]*line{l1, l2, l3}, 49)
	assert.Equal(t, 3, len(lns))
}

func TestJoin(t *testing.T) {
	l1, _ := parseLine("olia 1 0 100 a b c d")
	l2, _ := parseLine("olia 1 100 10 a b c c")
	l3, _ := parseLine("olia 1 210 100 a b c d")
	lns := joinLines([]*line{l1, l2, l3}, 100)
	assert.Equal(t, 2, len(lns))
	assert.Equal(t, "olia 1 0 110 a b c d", toStr(lns[0]))
	assert.Equal(t, "olia 1 210 100 a b c d", toStr(lns[1]))
}

func TestSeveral(t *testing.T) {
	l1, _ := parseLine("olia 1 0 100 a b c d")
	l2, _ := parseLine("olia 1 100 10 a b c c")
	l3, _ := parseLine("olia 1 110 10 a b c c")
	l4, _ := parseLine("olia 1 120 100 a b c d")
	lns := joinLines([]*line{l1, l2, l3, l4}, 100)
	assert.Equal(t, 2, len(lns))
	assert.Equal(t, "olia 1 0 120 a b c d", toStr(lns[0]))
	assert.Equal(t, "olia 1 120 100 a b c d", toStr(lns[1]))
}

func TestSpeakers(t *testing.T) {
	l1, _ := parseLine("olia 1 0 100 a b c a")
	l2, _ := parseLine("olia 1 100 10 a b c b")
	l3, _ := parseLine("olia 1 110 10 a b c c")
	l4, _ := parseLine("olia 1 120 100 a b c d")
	lns := oneSpeaker([]*line{l1, l2, l3, l4}, "")
	assert.Equal(t, 4, spTestCount(lns))
	lns = oneSpeaker([]*line{l1, l2, l3, l4}, "2")
	assert.Equal(t, 4, spTestCount(lns))
	lns = oneSpeaker([]*line{l1, l2, l3, l4}, "null")
	assert.Equal(t, 4, spTestCount(lns))
	lns = oneSpeaker([]*line{l1, l2, l3, l4}, "0")
	assert.Equal(t, 4, spTestCount(lns))
	lns = oneSpeaker([]*line{l1, l2, l3, l4}, "1")
	assert.Equal(t, 1, spTestCount(lns))
}

func spTestCount(lns []*line) int {
	s := make(map[string]bool)
	for _, l := range lns {
		s[l.speaker] = true
	}
	return len(s)
}

func TestJoinLast(t *testing.T) {
	l1, _ := parseLine("olia 1 0 100 a b c d")
	l2, _ := parseLine("olia 1 100 10 a b c c")
	lns := joinLines([]*line{l1, l2}, 100)
	assert.Equal(t, 1, len(lns))
	assert.Equal(t, "olia 1 0 110 a b c d", toStr(lns[0]))
}

func TestJoinFirst(t *testing.T) {
	l1, _ := parseLine("olia 1 0 10 a b c d")
	l2, _ := parseLine("olia 1 10 100 a b c c")
	lns := joinLines([]*line{l1, l2}, 100)
	assert.Equal(t, 1, len(lns))
	assert.Equal(t, "olia 1 0 110 a b c c", toStr(lns[0]))
}

func TestParseParams(t *testing.T) {
	params := &params{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	takeParams(fs, params)
	err := fs.Parse([]string{"-m", "10", "-speakers", "1"})
	assert.Nil(t, err)
	assert.Equal(t, 10, params.minMillis)
	assert.Equal(t, "1", params.speakers)
}

func TestParseParams_Speaker(t *testing.T) {
	params := &params{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	takeParams(fs, params)
	err := fs.Parse([]string{"-speakers", ""})
	assert.Nil(t, err)
	assert.Equal(t, "", params.speakers)
	err = fs.Parse([]string{"-speakers", "null"})
	assert.Nil(t, err)
	assert.Equal(t, "null", params.speakers)
	err = fs.Parse([]string{"-speakers", "2"})
	assert.Nil(t, err)
	assert.Equal(t, "2", params.speakers)
}
