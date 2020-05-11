package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/lattice"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/punctuation"
	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

func main() {
	log.SetOutput(os.Stderr)
	urlPtr := flag.String("u", "", "punctuation URL")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:[params] [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *urlPtr == "" {
		flag.Usage()
		log.Fatal("No punctuation URL provided")
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
	log.Printf("Punctuating")
	data, err = punctuate(data, punctuation.NewWrapper(*urlPtr))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't punctuate"))
	}
	err = lattice.Write(data, destination)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't write lattice"))
	}
	log.Print("Done punctuation")
}

func punctuate(data []*lattice.Part, p punctuation.Punctuator) ([]*lattice.Part, error) {
	l := len(data)
	i := 0
	var wg sync.WaitGroup
	errCh := make(chan error)
	for i < l {
		ni := getNextPartIndex(data, i)
		wg.Add(1)
		go invokePunc(data, i, ni, p, &wg, errCh)
		i = ni
	}
	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()
	select {
	case err := <-errCh:
		return nil, err
	case <-waitCh:
	}
	return data, nil
}

func invokePunc(data []*lattice.Part, i, ni int, p punctuation.Punctuator, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()
	words := getWords(data, i, ni)
	if (len(words)) > 0 {
		pWords, err := p.Punctuate(words)
		if err != nil {
			errCh <- errors.Wrap(err, "Can't punctuate")
		}
		err = addPunctuatioData(data, i, ni, pWords)
		if err != nil {
			errCh <- errors.Wrap(err, "Can't add punctuate result")
		}
	}
}

func getNextPartIndex(data []*lattice.Part, i int) int {
	l := len(data)
	if i >= l {
		return i
	}
	iPrev := i
	for sp := data[i].Speaker; i < l && sp == data[i].Speaker && noPause(data, iPrev, i); i++ {
		iPrev = i
	}
	return i
}

func noPause(data []*lattice.Part, from, to int) bool {
	if from == to {
		return true
	}
	tf := getLastWordTime(data[from].Words)
	if tf == 0 {
		return true
	}
	tt := getFirstWordTime(data[to].Words)
	if tt == 0 {
		return true
	}
	return (tt - tf) < (2 * time.Second)
}

func getLastWordTime(data []*lattice.Word) time.Duration {
	for i := len(data) - 1; i >= 0; i-- {
		w := data[i]
		if w.Main == lattice.MainInd {
			if w.Word != lattice.SilWord {
				return lattice.Duration(w.To)
			}
		}
	}
	return 0
}

func getFirstWordTime(data []*lattice.Word) time.Duration {
	for _, w := range data {
		if w.Main == lattice.MainInd {
			if w.Word != lattice.SilWord {
				return lattice.Duration(w.From)
			}
		}
	}
	return 0
}

func getWords(data []*lattice.Part, i int, to int) []string {
	res := make([]string, 0)
	for ; i < to; i++ {
		for _, w := range data[i].Words {
			if w.Main == lattice.MainInd {
				if w.Word != lattice.SilWord {
					res = append(res, trimWord(w.Word))
				}
			}
		}
	}
	return res
}

func trimWord(w string) string {
	w = strings.TrimLeft(w, "<")
	return strings.TrimRight(strings.TrimLeft(w, "<"), ">")
}

func addPunctuatioData(data []*lattice.Part, i int, to int, pResp *punctuation.Response) error {
	crI := 0
	for ; i < to; i++ {
		for _, w := range data[i].Words {
			if w.Main == lattice.MainInd {
				if w.Word != lattice.SilWord {
					pw, pp, err := getPunctuation(pResp, crI)
					if err != nil {
						return errors.Wrapf(err, "can't get punctuation for index %d", crI)
					}
					w.Punct = pp
					if strings.HasPrefix(w.Word, "<") && !strings.HasPrefix(pw, "<") {
						w.Word = "<" + pw + ">"
					} else {
						w.Word = pw
					}
					crI++
				}
			}
		}
	}
	return nil
}

func getPunctuation(pResp *punctuation.Response, i int) (string, string, error) {
	if pResp == nil {
		return "", "", errors.New("No punctuate result got")
	}
	if i >= len(pResp.Original) {
		return "", "", errors.Errorf("To short result array. Result len = %d", len(pResp.Original))
	}
	if i >= len(pResp.Punctuated) {
		return "", "", errors.Errorf("To short result array. Result len = %d", len(pResp.Punctuated))
	}
	l := len(pResp.Original[i])
	if l > len(pResp.Punctuated[i]) {
		return "", "", errors.Errorf("Punctuate result does not match %s vs %s ", pResp.Original[i], pResp.Punctuated[i])
	}
	if l == len(pResp.Punctuated[i]) {
		return pResp.Punctuated[i], "", nil
	}
	return pResp.Punctuated[i][:l], strings.TrimSpace(pResp.Punctuated[i][l:]), nil
}
