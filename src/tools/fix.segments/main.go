package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

type params struct {
	minMillis int
	speakers  string
}

type line struct {
	fields  []string
	from    int
	len     int
	rFields []string
	speaker string
}

func main() {
	log.SetOutput(os.Stderr)
	params := &params{}
	fs := flag.CommandLine
	takeParams(fs, params)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	err := validateParams(params)
	if err != nil {
		log.Print(err.Error())
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

	lns := make([]*line, 0)
	scanner := bufio.NewScanner(f)
	ln := 0
	for scanner.Scan() {
		ln++
		line := scanner.Text()
		if line != "" {
			l, err := parseLine(line)
			if err != nil {
				log.Fatal(errors.Wrapf(err, "Error on line %d", ln))
			}
			lns = append(lns, l)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Slice(lns, func(i, j int) bool { return lns[i].from < lns[j].from })

	log.Printf("Join segments shorter than %d ms", params.minMillis)
	lnsF := joinLines(lns, params.minMillis)
	lnsF = oneSpeaker(lnsF, params.speakers)

	for _, l := range lnsF {
		_, err := fmt.Fprintf(destination, "%s\n", toStr(l))
		if err != nil {
			log.Fatal(errors.Wrapf(err, "Can't write to output"))
		}
	}

	log.Printf("Done")
}

func takeParams(fs *flag.FlagSet, data *params) {
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s: <params> [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.IntVar(&data.minMillis, "m", 110, "Minimum ms to keep separate line")
	fs.StringVar(&data.speakers, "speakers", "", "Expected speakers count, if = '1' then all segments will have one speaker")
}

func validateParams(data *params) error {
	if data.minMillis < 0 {
		return errors.New("Wrong minimum ms specified")
	}
	return nil
}

func parseLine(str string) (*line, error) {
	res := &line{}
	res.fields = strings.Split(str, " ")
	if len(res.fields) < 8 {
		return nil, errors.Errorf("Wrong line %s", str)
	}
	var err error
	res.from, err = strconv.Atoi(res.fields[2])
	if err != nil {
		return nil, errors.Wrapf(err, "Wrong number in %s", str)
	}
	res.len, err = strconv.Atoi(res.fields[3])
	if err != nil {
		return nil, errors.Wrapf(err, "Wrong number in %s", str)
	}
	res.rFields = res.fields[4:7]
	res.speaker = res.fields[7]
	return res, nil
}

func joinLines(lns []*line, minMillis int) []*line {
	res := make([]*line, 0)
	var last *line
	for _, l := range lns {
		if last != nil && (l.len*10 <= minMillis || last.len*10 <= minMillis) {
			jl := l
			if last.len < l.len {
				last.rFields = l.rFields
				last.speaker = l.speaker
				jl = last
			}
			log.Printf("Join segment at %d-%d ms", jl.from*10, (jl.from+jl.len)*10)
			last.len = l.len + l.from - last.from
		} else {
			res = append(res, l)
			last = l
		}
	}
	return res
}

func oneSpeaker(lns []*line, spCount string) []*line {
	if spCount != "1" {
		return lns
	}
	if len(lns) < 1 {
		return lns
	}
	log.Printf("Fix to contain one speaker. Speakers = %s", spCount)

	sp := lns[0].speaker
	for _, l := range lns[1:] {
		l.speaker = sp
	}
	return lns
}

func toStr(l *line) string {
	res := strings.Builder{}
	res.WriteString(l.fields[0])
	res.WriteString(" ")
	res.WriteString(l.fields[1])
	res.WriteString(" ")
	res.WriteString(strconv.Itoa(l.from))
	res.WriteString(" ")
	res.WriteString(strconv.Itoa(l.len))
	for _, s := range l.rFields {
		res.WriteString(" ")
		res.WriteString(s)
	}
	res.WriteString(" ")
	res.WriteString(l.speaker)
	return res.String()
}
