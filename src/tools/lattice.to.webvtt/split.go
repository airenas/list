package main

import (
	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
)

func splitText(words []*lattice.Word) [][]int {
	res := make([][]int, 0)
	res = append(res, []int{0, len(words)})
	return res
}
