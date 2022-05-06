package webvtt

import (
	"strings"
	"testing"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	assert.True(t, strings.HasPrefix(Header(""), "WEBVTT"))
	assert.Equal(t, "WEBVTT - olia\n", Header("olia"))
	assert.Equal(t, "WEBVTT - loooooooong olia\n", Header("loooooooong olia"))
}

func TestGetWebVTT(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
 1 0.01 1.02 w
 1 1.02 2.02 w2
 `))
	text := extract(lat)

	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw w2\n", text)
}

func TestGetWebVTT_Underscore(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w_x_y
1 1.02 2.02 w2_a
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw x y w2 a\n", text)
}

func TestGetWebVTT_Several(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 w2

# 2 S2
1 5.01 6.02 w3
1 6.02 7.02 w4
`))
	text := extract(lat)

	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw w2\n\n00:05.010 --> 00:07.020\nw3 w4\n", text)
}

func extract(lat []*lattice.Part) string {
	lines := Extract(lat)
	sb := strings.Builder{}
	WriteTo(&sb, lines)
	return sb.String()
}

func TestGetSkipNonMain(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
0 1.02 2.02 w2
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:01.020\nw\n", text)
}

func TestGetSkipSil(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 2.54 w
1 2.55 2.02 <eps>
1 2.02 3.02 w2
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:03.020\nw w2\n", text)
}

func TestGetChopOnLongSil(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 2.02 <eps>
1 2.02 3.02 w2
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:01.020\nw\n\n00:02.020 --> 00:03.020\nw2\n", text)
}

func TestGetMinutes(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 62.02 w2
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 01:02.020\nw w2\n", text)
}

func TestGetHours(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 3602.02 w2
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 01:00:02.020\nw w2\n", text)
}

func TestGetLongHours(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w
1 1.02 3600182.52 w2
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 1000:03:02.520\nw w2\n", text)
}

func TestGetPuntc(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 W ,
1 1.02 2.02 w2 .
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nW, w2.\n", text)
}

func TestGetPuntcDash(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.02 w -
1 1.02 2.02 w2 .
`))
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:02.020\nw - w2.\n", text)
}

func TestSplit(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.00 w1
1 1.00 2.00 w2
1 2.00 3.00 w3
`))
	splitFunc = testSplitFunction
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:01.000\nw1\n\n00:01.000 --> 00:02.000\nw2\n\n00:02.000 --> 00:03.000\nw3\n", text)
}

func TestSplit2(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.00 w1
1 1.00 2.00 w2
1 2.00 3.00 w3
`))
	splitFunc = func(words []*lattice.Word) [][]int { return [][]int{{0, 1}, {1, 3}} }
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:01.000\nw1\n\n00:01.000 --> 00:03.000\nw2 w3\n", text)
}

func TestReplacesGreater(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
1 0.01 1.00 <w1>
`))
	splitFunc = splitText
	text := extract(lat)
	assert.Equal(t, "\n00:00.010 --> 00:01.000\n&lt;w1&gt;\n", text)
}

func TestWriteSpeaker(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
 1 0.01 1.02 w
 1 1.02 2.02 w2
 `))
	lines := Extract(lat)
	sb := strings.Builder{}
	WriteLineTo(&sb, lines[0], "speaker")
	text := sb.String()

	assert.Equal(t, "\n00:00.010 --> 00:02.020\n<v speaker>w w2\n", text)
}

func TestWriteSpeakerHTML(t *testing.T) {
	lat, _ := lattice.Read(strings.NewReader(
		`# 1 S1
 1 0.01 1.02 w
 1 1.02 2.02 w2
 `))
	lines := Extract(lat)
	sb := strings.Builder{}
	WriteLineTo(&sb, lines[0], "speaker><")
	text := sb.String()

	assert.Equal(t, "\n00:00.010 --> 00:02.020\n<v speaker&gt;&lt;>w w2\n", text)
}

func testSplitFunction(words []*lattice.Word) [][]int {
	res := make([][]int, 0)
	for i := range words {
		res = append(res, []int{i, i + 1})
	}
	return res
}

func Test_asHTML(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: "No changes", args: "olia", want: "olia"},
		{name: "Changes <>", args: "olia<a>", want: "olia&lt;a&gt;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := asHTML(tt.args); got != tt.want {
				t.Errorf("asHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sanitize(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: "No changes", args: "olia", want: "olia"},
		{name: "Changes <>", args: "olia<a>", want: "olia&lt;a&gt;"},
		{name: "Changes spaces", args: "olia < > ", want: "olia_&lt;_&gt;_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitize(tt.args); got != tt.want {
				t.Errorf("sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStyles(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "One sp", args: []string{"sp1"}, want: "STYLE\n::cue(v[voice=" + `"` + "sp1" + `"` + "]) { color: purple }\n"},
		{name: "Two sp", args: []string{"sp1", "sp2 2"}, want: "STYLE\n::cue(v[voice=" + `"` + "sp1" + `"` + "]) { color: purple }\n" +
			"::cue(v[voice=" + `"` + "sp2_2" + `"` + "]) { color: lightblue }\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStyles(tt.args); got != tt.want {
				t.Errorf("GetStyles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStyles(t *testing.T) {
	type args struct {
		speakers []string
		colors   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "One sp", args: args{speakers: []string{"sp1"}, colors: []string{"blue"}},
			want: "::cue(v[voice=" + `"` + "sp1" + `"` + "]) { color: blue }\n"},
		{name: "Two sp", args: args{speakers: []string{"sp1", "sp2"}, colors: []string{"blue"}},
			want: "::cue(v[voice=" + `"` + "sp1" + `"` + "]) { color: blue }\n" +
				"::cue(v[voice=" + `"` + "sp2" + `"` + "]) { color: blue }\n"},
		{name: "Resets collors", args: args{speakers: []string{"sp1", "sp2", "sp3"}, colors: []string{"blue", "red"}},
			want: "::cue(v[voice=" + `"` + "sp1" + `"` + "]) { color: blue }\n" +
				"::cue(v[voice=" + `"` + "sp2" + `"` + "]) { color: red }\n" +
				"::cue(v[voice=" + `"` + "sp3" + `"` + "]) { color: blue }\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStyles(tt.args.speakers, tt.args.colors); got != tt.want {
				t.Errorf("getStyles() = %v, want %v", got, tt.want)
			}
		})
	}
}
