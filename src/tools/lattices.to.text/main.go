package main

import (
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

func main() {
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: lattice-input-file0 [lattice-input-file1] ... [lattice-input-fileN] > text-file-out\n", os.Args[0])
		flag.PrintDefaults()
	}
	fnMap := ""
	flag.StringVar(&fnMap, "namesMap", "", "Map for ids to file names")
	flag.Parse()

	idsSpkMap := util.ParseSpeakers(fnMap)
	data, err := readFiles(flag.Args(), idsSpkMap)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't read lattices"))
	}

	destination := os.Stdout

	text := getText(data)
	_, err = destination.WriteString(text)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't write text"))
	}
	log.Print("Done convertion")
}

func extract(data []*lattice.Part) []*tdata {
	res := []*tdata{}
	c := &tdata{}
	ps := ""
	for i, p := range data {
		if ps != p.Speaker || lattice.SilDuration(data, i) >= (time.Second*2) {
			res, c = add(res, c)
		}
		ps = p.Speaker
		for _, w := range p.Words {
			if w.Main == lattice.MainInd {
				if w.Text != lattice.SilWord {
					c.words = append(c.words, w)
					c.to = lattice.Duration(w.To)
					if w.Punct == "." || w.Punct == "!" || w.Punct == "?" {
						res, c = add(res, c)
					}
				} else if len(c.words) > 0 && (lattice.Duration(w.To)-c.to) >= (time.Millisecond*200) {
					res, c = add(res, c)
				}
			}
		}
	}
	res, _ = add(res, c)
	return res
}

func add(res []*tdata, c *tdata) ([]*tdata, *tdata) {
	if len(c.words) > 0 {
		c.from = lattice.Duration(c.words[0].From)
		c.to = lattice.Duration(c.words[len(c.words)-1].To)
		res = append(res, c)
		c = &tdata{}
	}
	return res, c
}

func writeWord(res *strings.Builder, word, sep string) string {
	res.WriteString(sep + word)
	return " "
}

func writePunct(res *strings.Builder, punct string) string {
	if punct == "-" {
		res.WriteString(" ")
	}
	res.WriteString(punct)
	return " "
}

type fdata struct {
	speaker  string
	data     []*lattice.Part
	textData []*tdata
}

type tdata struct {
	from, to time.Duration
	words    []*lattice.Word
}

type pdata struct {
	n int // next index for fdata.textData
	d *fdata
}

func getText(data []*fdata) string {
	sb := &strings.Builder{}
	for i, d := range data {
		d.textData = extract(d.data)
		if d.speaker == "" {
			d.speaker = fmt.Sprintf("Kalbėtojas %d", i+1)
		}
	}

	pq := make(pqueue, 0)
	for _, d := range data {
		if len(d.textData) > 0 {
			heap.Push(&pq, &pdata{d: d, n: 0})
		}
	}

	ls := ""
	for len(pq) > 0 {
		item := heap.Pop(&pq).(*pdata)
		ls = writeLineTo(sb, item.d.textData[item.n].words, item.d.speaker, ls)
		if len(item.d.textData) > (item.n + 1) {
			item.n++
			heap.Push(&pq, item)
		}
	}
	return sb.String()
}

func writeLineTo(res *strings.Builder, words []*lattice.Word, speaker, lastSpeaker string) string {
	if res.Len() > 0 && lastSpeaker != speaker {
		res.WriteString("\n")
	}
	if lastSpeaker != speaker {
		res.WriteString(speaker)
		res.WriteString(":")
	}
	sep := " "
	for _, w := range words {
		writeWord(res, strings.Join(w.Words, " "), sep)
		sep = writePunct(res, w.Punct)
	}
	return speaker
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
	di := pq[i].d.textData[pq[i].n]
	dj := pq[j].d.textData[pq[j].n]
	return di.from < dj.from || (di.from == dj.from && pq[i].d.speaker < pq[j].d.speaker)
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
