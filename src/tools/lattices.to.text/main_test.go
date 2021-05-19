package main

import (
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

func TestGetText(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2
`))
	fd := &fdata{data: lat}
	text := getText([]*fdata{fd})
	assert.Equal(t, "Kalbėtojas 1: w w2", text)
}

func TestGetText_Speaker(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	text := getText([]*fdata{fd})
	assert.Equal(t, "sp: w w2", text)
}

func TestGetText_Several(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2
1 2.02 3.02 <eps>
1 10.02 20.02 w3
`))
	lat2, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 2.01 3.02 w
1 3.02 4.02 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	fd2 := &fdata{data: lat2, speaker: "sp2"}
	text := getText([]*fdata{fd, fd2})
	assert.Equal(t, "sp: w w2\nsp2: w w2\nsp: w3", text)
}

func TestGetTextUnderscore(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w_x_y
1 fr1 to2 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	text := getText([]*fdata{fd})
	assert.Equal(t, "sp: w x y w2", text)
}

func TestGetText_SkipSil(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w
1 fr1 to2 w2
1 fr to <eps>
1 fr to <eps>
`))
	fd := &fdata{data: lat, speaker: "sp"}
	text := getText([]*fdata{fd})
	assert.Equal(t, "sp: w w2", text)
}

func TestGetText_SkipNonMain(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w
0 fr1 to2 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	text := getText([]*fdata{fd})
	assert.Equal(t, "sp: w", text)
}

func TestGetText_Sep(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w ,
1 fr to w2`))
	fd := &fdata{data: lat, speaker: "sp"}
	text := getText([]*fdata{fd})
	assert.Equal(t, "sp: w, w2", text)
}

func TestGetText_SepDash(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w -
1 fr1 to2 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	text := getText([]*fdata{fd})
	assert.Equal(t, "sp: w - w2", text)
}

func TestGetText_SplitOnSep(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w .
1 1.02 2.02 w2 ?
1 2.02 3.02 w2 !
1 10.02 20.02 w3
`))
	lat2, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 1.00 2.00 w ,
1 2.00 4.00 w2 .
1 4.00 5.00 w2
`))
	fd := &fdata{data: lat, speaker: "sp"}
	fd2 := &fdata{data: lat2, speaker: "sp2"}
	text := getText([]*fdata{fd, fd2})
	assert.Equal(t, "sp: w.\nsp2: w, w2.\nsp: w2? w2!\nsp2: w2\nsp: w3", text)
}
