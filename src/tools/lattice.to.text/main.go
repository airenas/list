package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"github.com/pkg/errors"
)

func main() {
	filePtr := flag.String("f", "", "file in")
	outPtr := flag.String("o", "", "file for output")
	flag.Parse()

	if *filePtr == "" || *outPtr == "" {
		panic(errors.New("Usage: ./lattice.to.text -f <fileIn> -o <output file>"))
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
	text := getText(data)
	_, err = destination.WriteString(text)
	if err != nil {
		panic(errors.Wrap(err, "Can't write file "+*outPtr))
	}
	log.Print("Done convertion")
}

func getText(data []*lattice.Part) string {
	var res strings.Builder
	sep := ""
	pp := ""
	for _, p := range data {
		if pp != p.Speaker && sep != "" {
			sep = newLine(&res, sep)
		}
		pp = p.Speaker
		for _, w := range p.Words {
			if w.Main == lattice.MainInd {
				if w.Word == lattice.SilWord {
					sep = newLine(&res, sep)
				} else {
					sep = writeWord(&res, w.Word, sep)
					sep = writePunct(&res, w.Punct)
				}
			}

		}
	}
	return res.String()
}

func newLine(res *strings.Builder, sep string) string {
	if sep == "" {
		return ""
	}
	return "\n"
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
