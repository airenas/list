package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"container/heap"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/webvtt"
	"github.com/pkg/errors"
)

var version string

func main() {
	log.SetOutput(os.Stderr)

	fs := flag.CommandLine
	fnMap := ""
	strHeader := ""
	takeParams(fs, &fnMap, &strHeader)
	fs.Parse(os.Args[1:])

	log.Printf("lattices.to.webtt: %s\n", version)

	idsSpkMap := util.ParseSpeakers(fnMap)
	data, err := readFiles(flag.Args(), idsSpkMap)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't read lattices"))
	}

	destination := os.Stdout

	_, err = destination.WriteString(webvtt.Header(strHeader))
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
	n int // next index of fdata.vttData
	d *fdata
}

func getWebVTT(data []*fdata) string {
	sb := &strings.Builder{}
	for i, d := range data {
		d.vttData = webvtt.Extract(d.data)
		if d.speaker == "" {
			d.speaker = fmt.Sprintf("Kalbėtojas_%d", i+1)
		}
	}

	if len(data) > 1 {
		var names []string
		for _, d := range data {
			names = append(names, d.speaker)
		}
		sb.WriteString("\n")
		sb.WriteString(webvtt.GetStyles(names))
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

func readFiles(fns []string, idSp map[string]string) ([]*fdata, error) {
	if len(fns) < 1 {
		return nil, errors.New("no lattice files")
	}
	res := make([]*fdata, 0)
	for _, f := range fns {
		fd := &fdata{}
		var err error
		fd.data, err = readLattice(f)
		fd.speaker = util.GetSpeakerByPath(idSp, f)
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
	return di.From < dj.From || (di.From == dj.From && di.Speaker < dj.Speaker)
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *pqueue) Push(x interface{}) {
	item := x.(*pdata)
	*pq = append(*pq, item)
}

func (pq *pqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func takeParams(fs *flag.FlagSet, fnMap, header *string) {
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: lattice-input-file0 [lattice-input-file1] ... [lattice-input-fileN] > webvtt-file-out\n", os.Args[0])
		flag.PrintDefaults()
	}
	fs.StringVar(fnMap, "namesMap", "", "Map for ids to file names")
	fs.StringVar(header, "header", os.Getenv("WEBVTT_HEADER"), "WebVTT header string")
}
