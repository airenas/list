package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"container/heap"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/webvtt"
	"github.com/pkg/errors"
)

func main() {
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: input-file0 [input-file1] ... [input-fileN] > file_out\n", os.Args[0])
		flag.PrintDefaults()
	}
	fnMap := ""
	flag.StringVar(&fnMap, "namesMap", "", "Map for ids to file namas")
	flag.Parse()

	data, err := readFiles(flag.Args())
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't read lattices"))
	}

	destination := os.Stdout

	_, err = destination.WriteString(webvtt.Header())
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't write result lattice"))
	}
	text := getWebVTT(data)
	_, err = destination.WriteString(text)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't write result lattice"))
	}
	log.Print("Done generation")
}

type fdata struct {
	speaker string
	data    []*lattice.Part
	vttData []*webvtt.Line
}

type pdata struct {
	n int
	d *fdata
	i int //index in priority queue
}

func getWebVTT(data []*fdata) string {
	sb := &strings.Builder{}
	for i, d := range data {
		d.vttData = webvtt.Extract(d.data)
		d.speaker = fmt.Sprintf("Kalbėtojas %d", i+1)
	}

	pq := make(pqueue, 0)
	for _, d := range data {
		if len(d.vttData) > 0 {
			heap.Push(&pq, &pdata{d: d, n: 0})
		}
	}

	for len(pq) > 0 {
		item := heap.Pop(&pq).(*pdata)
		webvtt.WriteLineTo(sb, item.d.vttData[item.n], item.d.speaker)
		if len(item.d.vttData) > (item.n + 1) {
			item.n++
			heap.Push(&pq, item)
		}
	}

	return sb.String()
}

func readFiles(fns []string) ([]*fdata, error) {
	if len(fns) < 1 {
		return nil, errors.New("no lattice files")
	}
	res := make([]*fdata, 0)
	for _, f := range fns {
		fd := &fdata{}
		var err error
		fd.data, err = readLattice(f)
		if err != nil {
			return nil, errors.Wrapf(err, "can't read file %s ", f)
		}
		res = append(res, fd)
	}
	return res, nil
}

func readLattice(fn string) ([]*lattice.Part, error) {
	log.Printf("Open file " + fn)
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open file %s ", fn)
	}
	defer f.Close()
	return lattice.Read(f)
}

type pqueue []*pdata

func (pq pqueue) Len() int { return len(pq) }

func (pq pqueue) Less(i, j int) bool {
	di := pq[i].d.vttData[pq[i].n]
	dj := pq[j].d.vttData[pq[j].n]
	return di.From < dj.From || di.Speaker < dj.Speaker
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].i = i
	pq[j].i = j
}

func (pq *pqueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pdata)
	item.i = n
	*pq = append(*pq, item)
}

func (pq *pqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.i = -1    // for safety
	*pq = old[0 : n-1]
	return item
}
