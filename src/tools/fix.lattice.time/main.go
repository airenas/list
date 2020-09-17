package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

type params struct {
	len           float64
	silenceWord   string
	minSilenceLen float64 // in seconds
}

func main() {
	log.SetOutput(os.Stderr)
	params := &params{minSilenceLen: 0.5}
	fs := flag.CommandLine
	takeParams(fs, params)
	fs.Parse(os.Args[1:])
	err := validateParams(params)
	if err != nil {
		log.Printf(err.Error())
		fs.Usage()
		return
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
	log.Printf("Fixing")
	data, err = fixTime(data, params)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't fix time"))
	}
	if params.silenceWord != "" {
		data = changeSilenceWord(data, params)
	}
	err = lattice.Write(data, destination)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't write lattice"))
	}
	log.Print("Done")
}

func takeParams(fs *flag.FlagSet, data *params) {
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:[params] [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		flag.PrintDefaults()
	}
	fs.Float64Var(&data.len, "l", 0, "Len of audio file. Eg.: 1.23")
	fs.StringVar(&data.silenceWord, "s", "<tyla>", "Silence word symbol")
}

func validateParams(data *params) error {
	if data.len <= 0 {
		return errors.New("Wrong audio len specified")
	}
	if data.silenceWord == "" {
		return errors.New("No silence word symbol specified")
	}
	return nil
}

func fixTime(data []*lattice.Part, p *params) ([]*lattice.Part, error) {
	l := len(data)
	audioDuration := lattice.ToDuration(p.len)
	res := make([]*lattice.Part, len(data))
	ct := time.Duration(0)
	var lastWord *lattice.Word
	for i := 0; i < l; i++ {
		res[i] = &lattice.Part{Speaker: data[i].Speaker, Num: data[i].Num}
		for j := 0; j < len(data[i].Words); j++ {
			cw := data[i].Words[j]
			if cw.Main == lattice.MainInd {
				wtFrom := lattice.Duration(cw.From)
				wtTo := lattice.Duration(cw.To)
				if ct < wtFrom {
					if lastWord != nil && lattice.IsSilence(lastWord) {
						lastWord.To = cw.From
					} else if lattice.IsSilence(cw) {
						cw.From = lattice.DurationToText(ct)
					} else {
						res[i].Words = append(res[i].Words, &lattice.Word{From: lattice.DurationToText(ct),
							To: cw.From, Words: []string{lattice.SilWord},
							Main: lattice.MainInd})
					}
				}
				lastWord = cw
				ct = wtTo
			}
			res[i].Words = append(res[i].Words, cw)
		}
		if i == l-1 {
			if ct < audioDuration {
				if lastWord != nil && lattice.IsSilence(lastWord) {
					lastWord.To = lattice.DurationToText(audioDuration)
				} else {
					res[i].Words = append(res[i].Words, &lattice.Word{From: lattice.DurationToText(ct),
						To: lattice.DurationToText(audioDuration), Words: []string{lattice.SilWord},
						Main: lattice.MainInd})
				}
			}
		}
	}
	return res, nil
}

func changeSilenceWord(data []*lattice.Part, p *params) []*lattice.Part {
	minLen := lattice.ToDuration(p.minSilenceLen)
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i].Words); j++ {
			cw := data[i].Words[j]
			if cw.Main == lattice.MainInd && lattice.IsSilence(cw) &&
				lattice.WordDuration(cw) > minLen {
				cw.Words = []string{p.silenceWord}
			}
		}
	}
	return data
}
