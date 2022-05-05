package lattice

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	p := make([]*Part, 1)
	p[0] = &Part{Num: 1, Speaker: "S1"}
	p[0].Words = make([]*Word, 2)
	p[0].Words[0] = &Word{Main: "1", From: "fr", To: "to", Text: "word", Words: []string{"word"}}
	p[0].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Text: "word2", Words: []string{"word2"}}
	var b bytes.Buffer
	Write(p, &b)
	assert.Equal(t, "# 1 S1\n1 fr to word\n1 fr1 to2 word2\n\n", string(b.Bytes()))
}

func TestWriteSeveral(t *testing.T) {
	p := make([]*Part, 2)
	p[0] = &Part{Num: 1, Speaker: "S1"}
	p[0].Words = make([]*Word, 2)
	p[0].Words[0] = &Word{Main: "1", From: "fr", To: "to", Text: "word", Words: []string{"word"}}
	p[0].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Text: "word2", Words: []string{"word2"}}
	p[1] = &Part{Num: 2, Speaker: "S2"}
	p[1].Words = make([]*Word, 2)
	p[1].Words[0] = &Word{Main: "1", From: "fr", To: "to", Text: "word", Words: []string{"word"}}
	p[1].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Text: "word2", Words: []string{"word2"}}
	var b bytes.Buffer
	Write(p, &b)
	assert.Equal(t, "# 1 S1\n1 fr to word\n1 fr1 to2 word2\n\n# 2 S2\n1 fr to word\n1 fr1 to2 word2\n\n", string(b.Bytes()))
}

func TestWritePunc(t *testing.T) {
	p := make([]*Part, 1)
	p[0] = &Part{Num: 1, Speaker: "S1"}
	p[0].Words = make([]*Word, 2)
	p[0].Words[0] = &Word{Main: "1", From: "fr", To: "to", Text: "word", Punct: ",", Words: []string{"word"}}
	p[0].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Text: "word2", Punct: ".", Words: []string{"word2"}}
	var b bytes.Buffer
	Write(p, &b)
	assert.Equal(t, "# 1 S1\n1 fr to word ,\n1 fr1 to2 word2 .\n\n", string(b.Bytes()))
}

func TestWriteUnderscore(t *testing.T) {
	p := make([]*Part, 1)
	p[0] = &Part{Num: 1, Speaker: "S1"}
	p[0].Words = make([]*Word, 1)
	p[0].Words[0] = &Word{Main: "1", From: "fr", To: "to", Text: "word_x", Words: []string{"word", "x"}}
	var b bytes.Buffer
	Write(p, &b)
	assert.Equal(t, "# 2 S1\n1 fr to word_x\n\n", string(b.Bytes()))
}

func TestRead(t *testing.T) {
	p, _ := Read(strings.NewReader("# 1 S1\n1 fr to word\n1 fr1 to2 word2\n\n# 2 S2\n1 fr to word\n1 fr1 to2 word2\n\n"))

	assert.Equal(t, 2, len(p))
	assert.Equal(t, 1, p[0].Num)
	assert.Equal(t, "S1", p[0].Speaker)
	assert.Equal(t, 2, len(p[0].Words))
	assert.Equal(t, "1", p[0].Words[0].Main)
	assert.Equal(t, "fr", p[0].Words[0].From)
	assert.Equal(t, "to", p[0].Words[0].To)
	assert.Equal(t, "word", p[0].Words[0].Text)
	assert.Equal(t, []string{"word"}, p[0].Words[0].Words)
}

func TestReadUnderscoreWords(t *testing.T) {
	p, _ := Read(strings.NewReader("# 1 S1\n1 fr to word_x\n"))

	assert.Equal(t, 1, len(p))
	assert.Equal(t, 1, p[0].Num)
	assert.Equal(t, "S1", p[0].Speaker)
	assert.Equal(t, "1", p[0].Words[0].Main)
	assert.Equal(t, "fr", p[0].Words[0].From)
	assert.Equal(t, "to", p[0].Words[0].To)
	assert.Equal(t, "word_x", p[0].Words[0].Text)
	assert.Equal(t, []string{"word", "x"}, p[0].Words[0].Words)
}

func TestReadUnderscoreWords2(t *testing.T) {
	p, _ := Read(strings.NewReader("# 1 S1\n1 fr to word_x_y\n"))

	assert.Equal(t, 1, len(p))
	assert.Equal(t, 1, p[0].Num)
	assert.Equal(t, "word_x_y", p[0].Words[0].Text)
	assert.Equal(t, []string{"word", "x", "y"}, p[0].Words[0].Words)
}

func TestReadUnderscore(t *testing.T) {
	p, _ := Read(strings.NewReader("# 1 S1\n1 fr to _\n"))

	assert.Equal(t, 1, len(p))
	assert.Equal(t, 1, p[0].Num)
	assert.Equal(t, "_", p[0].Words[0].Text)
	assert.Equal(t, []string{"_"}, p[0].Words[0].Words)
}

func TestReadPunct(t *testing.T) {
	p, _ := Read(strings.NewReader("# 1 S1\n1 fr to word .\n\n"))

	assert.Equal(t, 1, len(p))
	assert.Equal(t, 1, len(p[0].Words))
	assert.Equal(t, ".", p[0].Words[0].Punct)
	assert.Equal(t, "word", p[0].Words[0].Text)
}

func TestRead_Fail(t *testing.T) {
	p, err := Read(strings.NewReader("1 fr to word .\n\n"))

	assert.NotNil(t, err)
	assert.Nil(t, p)
}

func TestWordDuration(t *testing.T) {
	w := &Word{From: "0.00", To: "10.00"}
	d := WordDuration(w)
	assert.Equal(t, 10*time.Second, d)
}

func TestDurationParse(t *testing.T) {
	assert.Equal(t, 10*time.Second+50*time.Millisecond, Duration("10.05"))
	assert.Equal(t, 10*time.Millisecond, Duration("0.01"))
	assert.Equal(t, 100*time.Millisecond, Duration("0.10"))
}

func TestDurationDefault(t *testing.T) {
	d := Duration("apoa")
	assert.Equal(t, 0*time.Second, d)
}

func TestIsSilence(t *testing.T) {
	assert.True(t, IsSilence(&Word{Words: []string{"<eps>"}}))
	assert.False(t, IsSilence(&Word{Words: []string{}}))
	assert.False(t, IsSilence(&Word{Words: []string{""}}))
	assert.False(t, IsSilence(&Word{Words: []string{"olia", "bet"}}))
	assert.False(t, IsSilence(&Word{Words: []string{"<eps>", "bet"}}))
}

func TestDurationToText(t *testing.T) {
	assert.Equal(t, "2.13", DurationToText(ToDuration(2.13)))
	assert.Equal(t, "2.00", DurationToText(ToDuration(2.001)))
	assert.Equal(t, "0.00", DurationToText(ToDuration(0.00)))
	assert.Equal(t, "100.54", DurationToText(ToDuration(100.539)))
}
