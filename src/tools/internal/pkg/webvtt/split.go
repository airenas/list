package webvtt

import (
	"math"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
)

type splitParamsType struct {
	wantedCharCount int
	exceededPenalty float64
	dotPenalty      float64 //applies for ?!
	commaPenalty    float64
	pausePenalty    float64
}

var splitParams splitParamsType

func init() {
	splitParams = splitParamsType{
		wantedCharCount: 30,
		exceededPenalty: 4,
		dotPenalty:      15,
		commaPenalty:    7,
		pausePenalty:    10,
	}
}

func splitText(words []*lattice.Word) [][]int {
	l := len(words)
	splitInd := make([]bool, l)
	c := calc(words, splitInd, l)
	nc, i := findBestSplit(words, splitInd)
	for nc < c {
		c = nc
		splitInd[i] = !splitInd[i]
		nc, i = findBestSplit(words, splitInd)
	}
	res := make([][]int, 0)
	pr := 0
	for i, v := range splitInd {
		if v {
			res = append(res, []int{pr, i + 1})
			pr = i + 1
		}
	}
	if pr < l {
		res = append(res, []int{pr, l})
	}
	return res
}

func calc(words []*lattice.Word, splitInd []bool, ch int) float64 {
	pi := 0
	v := float64(0)
	for i := 0; i < (len(words) - 1); i++ {
		if (splitInd[i] && i != ch) || (!splitInd[i] && i == ch) {
			v += calcSegment(words[pi:i+1]) + pauseCost(words[i], words[i+1])
			pi = i + 1
		}
	}
	if pi != len(words) {
		v += calcSegment(words[pi:])
	}
	return v
}

func calcSegment(words []*lattice.Word) float64 {
	lc := splitParams.wantedCharCount - lenText(words)
	res := float64(lc)
	if lc < 0 {
		res = -res * splitParams.exceededPenalty
	}

	return res + splitParams.dotPenalty*float64(dotCount(words)) + splitParams.commaPenalty*float64(commaCount(words))
}

func lenText(words []*lattice.Word) int {
	res := 0
	for _, w := range words {
		res = res + len(w.Words) - 1
		for _, wt := range w.Words {
			res = res + len(wt)
		}
		if w.Punct != "" {
			res = res + 1
		}
	}
	return res
}

func pauseCost(w1, w2 *lattice.Word) float64 {
	d := lattice.Duration(w2.From) - lattice.Duration(w1.To)
	return math.Abs(0.5-d.Seconds()) * splitParams.pausePenalty
}

func dotCount(words []*lattice.Word) int {
	res := 0
	for _, w := range words[:len(words)-1] {
		if w.Punct == "." || w.Punct == "!" || w.Punct == "?" {
			res = res + 1
		}
	}
	return res
}

func commaCount(words []*lattice.Word) int {
	res := 0
	for _, w := range words[:len(words)-1] {
		if w.Punct == "," {
			res = res + 1
		}
	}
	return res
}

func findBestSplit(words []*lattice.Word, splitInd []bool) (float64, int) {
	rc, ri := calc(words, splitInd, 0), 0
	for i := 1; i < (len(words) - 1); i++ {
		c := calc(words, splitInd, i)
		if rc > c {
			rc, ri = c, i
		}
	}
	return rc, ri
}
