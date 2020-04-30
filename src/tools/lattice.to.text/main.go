package main

import (
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
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	f, err := util.NewReadWrapper(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	destination, err := util.NewWriteWrapper(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	data, err := lattice.Read(f)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't read lattice"))
	}
	text := getText(data)
	_, err = destination.WriteString(text)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't write text"))
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
					if lattice.WordDuration(w) > (time.Second * 2) {
						sep = newLine(&res, sep)
					}
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
