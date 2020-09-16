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
	len float64
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

	_, err = fmt.Fprint(destination, makeEmptyLattice(params.len))
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
}

func validateParams(data *params) error {
	if data.len <= 0 {
		return errors.New("Wrong audio len specified")
	}
	return nil
}

func makeEmptyLattice(l float64) string {
	res := strings.Builder{}
	res.WriteString("# 1 S0000\n")
	res.WriteString(fmt.Sprintf("1 0 %s <eps>\n", round(l)))
	return res.String()
}

func round(l float64) string {
	return fmt.Sprintf("%.2f", l)
}
