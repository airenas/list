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
	p[0].Words[0] = &Word{Main: "1", From: "fr", To: "to", Word: "word"}
	p[0].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Word: "word2"}
	var b bytes.Buffer
	Write(p, &b)
	assert.Equal(t, "# 1 S1\n1 fr to word\n1 fr1 to2 word2\n\n", string(b.Bytes()))
}

func TestWriteSeveral(t *testing.T) {
	p := make([]*Part, 2)
	p[0] = &Part{Num: 1, Speaker: "S1"}
	p[0].Words = make([]*Word, 2)
	p[0].Words[0] = &Word{Main: "1", From: "fr", To: "to", Word: "word"}
	p[0].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Word: "word2"}
	p[1] = &Part{Num: 2, Speaker: "S2"}
	p[1].Words = make([]*Word, 2)
	p[1].Words[0] = &Word{Main: "1", From: "fr", To: "to", Word: "word"}
	p[1].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Word: "word2"}
	var b bytes.Buffer
	Write(p, &b)
	assert.Equal(t, "# 1 S1\n1 fr to word\n1 fr1 to2 word2\n\n# 2 S2\n1 fr to word\n1 fr1 to2 word2\n\n", string(b.Bytes()))
}

func TestWritePunc(t *testing.T) {
	p := make([]*Part, 1)
	p[0] = &Part{Num: 1, Speaker: "S1"}
	p[0].Words = make([]*Word, 2)
	p[0].Words[0] = &Word{Main: "1", From: "fr", To: "to", Word: "word", Punct: ","}
	p[0].Words[1] = &Word{Main: "1", From: "fr1", To: "to2", Word: "word2", Punct: "."}
	var b bytes.Buffer
	Write(p, &b)
	assert.Equal(t, "# 1 S1\n1 fr to word ,\n1 fr1 to2 word2 .\n\n", string(b.Bytes()))
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
	assert.Equal(t, "word", p[0].Words[0].Word)
}

func TestReadPunct(t *testing.T) {
	p, _ := Read(strings.NewReader("# 1 S1\n1 fr to word .\n\n"))

	assert.Equal(t, 1, len(p))
	assert.Equal(t, 1, len(p[0].Words))
	assert.Equal(t, ".", p[0].Words[0].Punct)
	assert.Equal(t, "word", p[0].Words[0].Word)
}

func TestRead_Fail(t *testing.T) {
	p, err := Read(strings.NewReader("1 fr to word .\n\n"))

	assert.NotNil(t, err)
	assert.Nil(t, p)
}

func TestDuration(t *testing.T) {
	w := &Word{From: "0.00", To: "10.00"}
	d := Duration(w)
	assert.Equal(t, 10*time.Second, d)
}

func TestDurationParse(t *testing.T) {
	d := duration("10.05")
	assert.Equal(t, 10*time.Second+50*time.Millisecond, d)
}

func TestDurationDefault(t *testing.T) {
	d := duration("apoa")
	assert.Equal(t, 0*time.Second, d)
}
