package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/pkg/errors"
)

func main() {
	filePtr := flag.String("f", "", "file in")
	outPtr := flag.String("o", "", "file for output")
	flag.Parse()

	if *filePtr == "" || *outPtr == "" {
		panic(errors.New("Usage: ./lattice.to.webvtt -f <fileIn> -o <output file>"))
	}

	f, err := os.Open(*filePtr)
	if err != nil {
		panic(errors.Wrapf(err, "Can't read file %s ", *filePtr))
	}
	defer f.Close()

	destination, err := os.Create(*outPtr)
	if err != nil {
		panic(errors.Wrap(err, "Can't create file "+*outPtr))
	}
	defer destination.Close()

	log.Printf("Reading file %s", *filePtr)
	data, err := lattice.Read(f)
	if err != nil {
		panic(errors.Wrap(err, "Can't read file "+*filePtr))
	}
	_, err = destination.WriteString(getHeader())
	text := getWebVTT(data)
	_, err = destination.WriteString(text)
	if err != nil {
		panic(errors.Wrap(err, "Can't write file "+*outPtr))
	}
	log.Print("Done convertion")
}

type tdata struct {
	sb       strings.Builder
	speaker  string
	from, to *lattice.Word
	sbLine   strings.Builder
	sep      string
}

func getHeader() string {
	return "WEBVTT\n"
}

func getWebVTT(data []*lattice.Part) string {
	td := &tdata{}
	for _, p := range data {
		if p.Speaker != td.speaker {
			write(td)
			td.speaker = p.Speaker
		}
		for _, w := range p.Words {
			if w.Main == lattice.MainInd {
				if w.Word != lattice.SilWord {
					if td.to != nil && (lattice.Duration(w.From)-lattice.Duration(td.to.To)) > 500*time.Millisecond {
						write(td)
					}
					add(td, w)
				}
			}
		}
	}
	write(td)
	return td.sb.String()
}

func write(td *tdata) {
	if td.from != nil {
		td.sb.WriteString("\n")
		td.sb.WriteString(asString(lattice.Duration(td.from.From)))
		td.sb.WriteString(" --> ")
		td.sb.WriteString(asString(lattice.Duration(td.to.To)))
		td.sb.WriteString("\n")
		td.sb.WriteString(td.sbLine.String())
		td.sb.WriteString("\n")
	}
	td.sbLine.Reset()
	td.from = nil
	td.to = nil
	td.sep = ""
}

func add(td *tdata, w *lattice.Word) {
	if td.from == nil {
		td.from = w
	}
	td.to = w
	td.sep = writeWord(&td.sbLine, w.Word, td.sep)
	td.sep = writePunct(&td.sbLine, w.Punct)
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

func asString(d time.Duration) string {
	res := ""
	ms := d.Milliseconds()
	tv := int64(ms) / time.Hour.Milliseconds()
	if tv > 0 {
		res += padString(tv, 2) + ":"
	}
	ms = ms % time.Hour.Milliseconds()
	tv = ms / time.Minute.Milliseconds()
	res += padString(tv, 2) + ":"
	ms = ms % time.Minute.Milliseconds()
	tv = ms / time.Second.Milliseconds()
	res += padString(tv, 2) + "."
	ms = ms % time.Second.Milliseconds()
	res += padString(ms, 3)
	return res
}

func padString(v int64, l int) string {
	res := strconv.Itoa(int(v))
	for len(res) < l {
		res = "0" + res
	}
	return res
}
