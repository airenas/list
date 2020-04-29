package lattice

import (
	"bytes"
	"testing"

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
