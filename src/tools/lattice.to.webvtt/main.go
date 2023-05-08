package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/webvtt"
	"github.com/pkg/errors"
)

func main() {
	log.SetOutput(os.Stderr)
	fs := flag.CommandLine
	strHeader := ""
	takeParams(fs, &strHeader)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	
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
	_, err = destination.WriteString(webvtt.Header(strHeader))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't write result lattice"))
	}
	text := getWebVTT(data)
	_, err = destination.WriteString(text)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't write result lattice"))
	}
	log.Print("Done generation")
}

func takeParams(fs *flag.FlagSet, header *string) {
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		flag.PrintDefaults()
	}
	fs.StringVar(header, "header", os.Getenv("WEBVTT_HEADER"), "WEbVTT header string")
}

func getWebVTT(data []*lattice.Part) string {
	vttData := webvtt.Extract(data)
	sb := &strings.Builder{}
	webvtt.WriteTo(sb, vttData)
	return sb.String()
}
