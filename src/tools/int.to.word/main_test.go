package main

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	w, n, err := parseLine("olia 10")
	assert.Nil(t, err)
	assert.Equal(t, "olia", w)
	assert.Equal(t, 10, n)

	w, n, err = parseLine("olia 101010 ")
	assert.Nil(t, err)
	assert.Equal(t, "olia", w)
	assert.Equal(t, 101010, n)

	w, n, err = parseLine("olia 101010\n")
	assert.Nil(t, err)
	assert.Equal(t, "olia", w)
	assert.Equal(t, 101010, n)

	w, n, err = parseLine("olia101010")
	assert.NotNil(t, err)
}

func TestToString(t *testing.T) {
	str := toString([]string{"olia", "10"})
	assert.Equal(t, "olia 10", str)

	str = toString([]string{"olia"})
	assert.Equal(t, "olia", str)

	str = toString([]string{"olia", "qqq", "aaa"})
	assert.Equal(t, "olia qqq aaa", str)
}

func TestReadVocab(t *testing.T) {
	rd := strings.NewReader("olia 0\nolia1 1")
	v, err := readVocab(rd)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(v), 2)
	assert.Equal(t, "olia", v[0])
	assert.Equal(t, "olia1", v[1])
	assert.True(t, len(v) < 3 || v[2] == "")
}

func TestReadVocabFail(t *testing.T) {
	rd := strings.NewReader("olia 0\nolia1 aa")
	_, err := readVocab(rd)
	assert.NotNil(t, err)
}

func TestReadVocabFailNoSpace(t *testing.T) {
	rd := strings.NewReader("olia 0\nolia12")
	_, err := readVocab(rd)
	assert.NotNil(t, err)
}

func TestReadVocabFailNonLast(t *testing.T) {
	rd := strings.NewReader("olia 0\nolia2\nopa 2")
	_, err := readVocab(rd)
	assert.NotNil(t, err)
	rd = strings.NewReader("olia aaa\nolia 1\nopa 2")
	_, err = readVocab(rd)
	assert.NotNil(t, err)
}

func TestReadLarger(t *testing.T) {
	n := 20000
	rd := generate(n)
	v, err := readVocab(rd)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(v), 2)
	for i := 0; i < n; i++ {
		assert.Equal(t, fmt.Sprintf("w_%d", i), v[i])
	}
}

func TestMap(t *testing.T) {
	rd := generate(10)
	v, err := readVocab(rd)
	assert.Nil(t, err)
	nl, err := mapLine("a 1 b", v, 1)
	assert.Nil(t, err)
	assert.Equal(t, "a w_1 b", nl)

	nl, err = mapLine("1 a", v, 0)
	assert.Nil(t, err)
	assert.Equal(t, "w_1 a", nl)

	nl, err = mapLine("a b 1", v, 2)
	assert.Nil(t, err)
	assert.Equal(t, "a b w_1", nl)

	nl, err = mapLine("a b 1 1", v, 2)
	assert.Nil(t, err)
	assert.Equal(t, "a b w_1 1", nl)
}

func TestMapSkip(t *testing.T) {
	rd := generate(10)
	v, err := readVocab(rd)
	assert.Nil(t, err)
	nl, err := mapLine("a 1 b", v, 5)
	assert.Nil(t, err)
	assert.Equal(t, "a 1 b", nl)
}

func TestMap_FailNotWord(t *testing.T) {
	rd := generate(10)
	v, err := readVocab(rd)
	assert.Nil(t, err)
	_, err = mapLine("a bbb", v, 1)
	assert.NotNil(t, err)
}

func TestMap_FailNoWord(t *testing.T) {
	rd := generate(10)
	v, err := readVocab(rd)
	assert.Nil(t, err)
	_, err = mapLine("a 10", v, 1)
	assert.NotNil(t, err)
}

func benchmarkParse(b *testing.B, l string) {
	for n := 0; n < b.N; n++ {
		parseLine(l)
	}
}

func BenchmarkParse1(b *testing.B) { benchmarkParse(b, "olia 15") }
func BenchmarkParse2(b *testing.B) { benchmarkParse(b, "oliaaaaaaaaaaaa 15\n") }

func benchmarkToString(b *testing.B, str []string) {
	for n := 0; n < b.N; n++ {
		toString(str)
	}
}

func BenchmarkToString1(b *testing.B) { benchmarkToString(b, []string{"olia"}) }
func BenchmarkToString2(b *testing.B) { benchmarkToString(b, []string{"olia", "olia"}) }
func BenchmarkToString3(b *testing.B) { benchmarkToString(b, []string{"olia", "olia", "olia"}) }
func BenchmarkToString5(b *testing.B) {
	benchmarkToString(b, []string{"olia", "olia", "olia", "olia", "olia"})
}

func benchmarkReadVocab(b *testing.B, n int, pJobs int) {
	str := generateStr(n)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		readVocabInt(strings.NewReader(str), pJobs)
	}
}

func BenchmarkReadVocab10_4(b *testing.B)        { benchmarkReadVocab(b, 10, 4) }
func BenchmarkReadVocab100_4(b *testing.B)       { benchmarkReadVocab(b, 100, 4) }
func BenchmarkReadVocab1000_4(b *testing.B)      { benchmarkReadVocab(b, 1000, 4) }
func BenchmarkReadVocab10000000_1(b *testing.B)  { benchmarkReadVocab(b, 1000000, 1) }
func BenchmarkReadVocab10000000_2(b *testing.B)  { benchmarkReadVocab(b, 1000000, 2) }
func BenchmarkReadVocab10000000_3(b *testing.B)  { benchmarkReadVocab(b, 1000000, 3) }
func BenchmarkReadVocab10000000_4(b *testing.B)  { benchmarkReadVocab(b, 1000000, 4) }
func BenchmarkReadVocab10000000_5(b *testing.B)  { benchmarkReadVocab(b, 1000000, 5) }
func BenchmarkReadVocab10000000_6(b *testing.B)  { benchmarkReadVocab(b, 1000000, 6) }
func BenchmarkReadVocab10000000_20(b *testing.B) { benchmarkReadVocab(b, 1000000, 20) }
func BenchmarkReadVocab20000000_4(b *testing.B)  { benchmarkReadVocab(b, 2000000, 4) }

func generate(size int) io.Reader {
	return strings.NewReader(generateStr(size))
}

func generateStr(size int) string {
	res := strings.Builder{}
	for i := 0; i < size; i++ {
		res.WriteString(fmt.Sprintf("w_%d %d\n", i, i))
	}
	return res.String()
}
