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
	segmentName   string
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
	err = lattice.Write(data, destination)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't write lattice"))
	}
	log.Print("Done")
}

func takeParams(fs *flag.FlagSet, data *params) {
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:[params] [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.Float64Var(&data.len, "l", 0, "Len of audio file. Eg.: 1.23")
	fs.StringVar(&data.silenceWord, "s", "<tyla>", "Silence word symbol")
	fs.StringVar(&data.segmentName, "sn", "TYLA", "Segment name")
	fs.Float64Var(&data.minSilenceLen, "ms", 0.10, "Min missing segment len in sec to insert")
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

func fixTime(data []*lattice.Part, p *params) ([]*lattice.Part, error) {
	l := len(data)
	audioDuration := lattice.ToDuration(p.len)
	minLen := lattice.ToDuration(p.minSilenceLen)
	res := make([]*lattice.Part, 0)
	ct := time.Duration(0)
	for i := 0; i < l; i++ {
		from, to, err := getSegmentTimes(data[i])
		if err != nil {
			return nil, errors.Wrap(err, "Can't get times")
		}
		if (from - ct) > minLen {
			log.Printf("Insert segment from %s, to %s", lattice.DurationToText(ct), lattice.DurationToText(from))
			res = append(res, newSegment(ct, from, p))
		}
		res = append(res, data[i])
		ct = to
	}
	if (audioDuration - ct) > minLen {
		log.Printf("Insert segment from %s, to %s", lattice.DurationToText(ct), lattice.DurationToText(audioDuration))
		res = append(res, newSegment(ct, audioDuration, p))
	}
	// fix segment num
	for i := 0; i < len(res); i++ {
		res[i].Num = i + 1
	}
	return res, nil
}

func getSegmentTimes(part *lattice.Part) (time.Duration, time.Duration, error) {
	from, err := getSegmentTimeFrom(part)
	if err != nil {
		return time.Duration(0), time.Duration(0), errors.Wrap(err, "Can't get from time")
	}
	to, err := getSegmentTimeTo(part)
	if err != nil {
		return time.Duration(0), time.Duration(0), errors.Wrap(err, "Can't get to time")
	}
	return from, to, nil
}

func getSegmentTimeFrom(part *lattice.Part) (time.Duration, error) {
	for i := 0; i < len(part.Words); i++ {
		cw := part.Words[i]
		if cw.Main == lattice.MainInd {
			return lattice.Duration(cw.From), nil
		}
	}
	return time.Duration(0), errors.New("No time")
}

func getSegmentTimeTo(part *lattice.Part) (time.Duration, error) {
	for i := len(part.Words) - 1; i >= 0; i-- {
		cw := part.Words[i]
		if cw.Main == lattice.MainInd {
			return lattice.Duration(cw.To), nil
		}
	}
	return time.Duration(0), errors.New("No time")
}

func newSegment(from, to time.Duration, p *params) *lattice.Part {
	res := &lattice.Part{Speaker: p.segmentName}
	res.Words = append(res.Words, &lattice.Word{From: lattice.DurationToText(from),
		To: lattice.DurationToText(to), Main: lattice.MainInd, Words: []string{p.silenceWord}})
	return res
}
