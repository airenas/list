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
1 fr to w
1 fr1 to2 w2
`))
	text := getText(lat)
	assert.Equal(t, "w w2", text)
}

func TestGetText_SkipSil(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w
1 fr1 to2 w2
1 fr to <eps>
1 fr to <eps>
`))
	text := getText(lat)
	assert.Equal(t, "w w2", text)
}

func TestGetText_SkipNonMain(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w
0 fr1 to2 w2
`))
	text := getText(lat)
	assert.Equal(t, "w", text)
}

func TestGetText_SeveralSpeaker(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w
0 fr1 to2 w2

# 1 S1
1 fr to w
0 fr1 to2 w2
`))
	text := getText(lat)
	assert.Equal(t, "w w", text)
}

func TestGetText_NewLineSpeaker(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w

# 1 S2
1 fr to w2
`))
	text := getText(lat)
	assert.Equal(t, "w\nw2", text)
}

func TestGetText_NewLineSil(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w
1 fr to <eps>
1 fr1 to2 w2
`))
	text := getText(lat)
	assert.Equal(t, "w\nw2", text)
}

func TestGetText_Sep(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w ,
1 fr to w2`))
	text := getText(lat)
	assert.Equal(t, "w, w2", text)
}

func TestGetText_SepDash(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 fr to w -
1 fr1 to2 w2
`))
	text := getText(lat)
	assert.Equal(t, "w - w2", text)
}
