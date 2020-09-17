package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

type params struct {
	len         float64
	silenceWord string
	segmentName string
}

func main() {
	log.SetOutput(os.Stderr)
	params := &params{}
	fs := flag.CommandLine
	takeParams(fs, params)
	fs.Parse(os.Args[1:])
	err := validateParams(params)
	if err != nil {
		log.Printf(err.Error())
		fs.Usage()
		return
	}

	destination, err := util.NewWriteWrapper(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	_, err = fmt.Fprint(destination, makeEmptyLattice(params))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Can't write to output"))
	}
	log.Printf("Done")
}

func takeParams(fs *flag.FlagSet, data *params) {
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s: <params> [output-file | stdout]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.Float64Var(&data.len, "l", 0, "Len of audio file. Eg.: 1.23")
	fs.StringVar(&data.silenceWord, "s", "<tyla>", "Silence word symbol")
	fs.StringVar(&data.segmentName, "sn", "TYLA", "Segment name")
}

func validateParams(data *params) error {
	if data.len <= 0 {
		return errors.New("Wrong audio len specified")
	}
	if data.silenceWord == "" {
		return errors.New("No silence word symbol specified")
	}
	if data.segmentName == "" {
		return errors.New("No segment name specified")
	}

	return nil
}

func makeEmptyLattice(p *params) string {
	res := strings.Builder{}
	res.WriteString(fmt.Sprintf("# 1 %s\n", p.segmentName))
	res.WriteString(fmt.Sprintf("1 0 %s %s\n", round(p.len), p.silenceWord))
	return res.String()
}

func round(l float64) string {
	return fmt.Sprintf("%.2f", l)
}
