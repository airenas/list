package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

type params struct {
	c     int
	vocab string
}

func main() {
	//defer profile.Start().Stop()
	log.SetOutput(os.Stderr)
	params := &params{}
	fs := flag.CommandLine
	takeParams(fs, params)
	fs.Parse(os.Args[1:])
	params.c-- // make zero based
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

	log.Printf("Open vocab " + params.vocab)
	vFile, err := os.Open(params.vocab)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Can't open vocab %s ", params.vocab))
	}
	defer vFile.Close()
	vocab, err := readVocab(vFile)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Can't read vocab %s", params.vocab))
	}

	destination, err := util.NewWriteWrapper(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	scanner := bufio.NewScanner(f)
	ln := 0
	for scanner.Scan() {
		ln++
		line := scanner.Text()
		nLine, err := mapLine(line, vocab, params.c)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "Error on line %d", ln))
		}
		_, err = fmt.Fprintf(destination, "%s\n", nLine)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Cant write file"))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done")
}

func takeParams(fs *flag.FlagSet, data *params) {
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s: <params> [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.StringVar(&data.vocab, "v", "", "Vocabulary")
	fs.IntVar(&data.c, "c", 0, "Column to change")
}

func validateParams(data *params) error {
	if data.vocab == "" {
		return errors.New("No vocab")
	}
	if data.c < 0 {
		return errors.New("No column specified")
	}
	return nil
}

func getSize(data []byte) (int, error) {
	l, err := getLastLine(data)
	if err != nil {
		return 0, err
	}
	_, n, err := parseLine(l)
	return n, err
}

func getLastLine(data []byte) (string, error) {
	l := len(data)
	res := ""
	l--
	for l--; l > 0 && res == ""; l-- {
		if data[l] == '\n' {
			res = strings.TrimSpace(string(data[l:]))
		}
	}
	return res, nil
}

func getTo(data []byte, start int) int {
	l := len(data)
	for ; start < l; start++ {
		if data[start] == '\n' {
			return start + 1
		}
	}
	return start
}

func parseLine(line string) (string, int, error) {
	ind := strings.Index(line, " ")
	if ind > -1 {
		num, err := strconv.Atoi(strings.TrimSpace(line[ind+1:]))
		if err != nil {
			return "", -1, errors.Wrapf(err, "Wrong number in %s", line)
		}
		return line[:ind], num, nil
	}
	return "", -1, errors.New("No space")
}

func readVocab(src io.Reader) ([]string, error) {
	return readVocabInt(src, 4)
}

func readVocabInt(src io.Reader, pJobs int) ([]string, error) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	c := len(data)
	n, err := getSize(data)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't get vocab size")
	}
	res := &results{res: make([]string, n+1), lock: &sync.Mutex{}}
	wg := sync.WaitGroup{}
	cd := c / pJobs
	toC := cd
	for i := 0; i < c; i = toC {
		toC = i + cd
		if toC > c {
			toC = c
		} else {
			toC = getTo(data, toC)
		}
		wg.Add(1)
		go read(data[i:toC], res, &wg)
	}
	wg.Wait()
	return res.res, res.err
}

type results struct {
	err  error
	lock *sync.Mutex
	res  []string
}

func setError(res *results, err error) {
	res.lock.Lock()
	defer res.lock.Unlock()
	if res.err == nil {
		res.err = err
	}
}

func read(data []byte, res *results, wg *sync.WaitGroup) {
	defer wg.Done()

	buf := bytes.NewBuffer(data)
	var err error
	var line string
	for err != io.EOF {
		line, err = buf.ReadString('\n')
		if err != nil && err != io.EOF {
			setError(res, errors.Wrapf(err, "Can't read"))
			return
		}
		line = strings.TrimSpace(line)
		if line != "" {
			w, num, err := parseLine(line)
			if err != nil {
				setError(res, errors.Wrapf(err, "Can't parse %s", line))
				return
			}
			for num >= len(res.res) {
				setError(res, errors.Wrapf(err, "Too small vocab %d. Wanted %d. Line: %s", len(res.res), num, line))
				return
			}
			res.res[num] = w
		}
	}
}

func mapLine(line string, vocab []string, c int) (string, error) {
	strs := strings.Split(line, " ")
	var w string
	if c < len(strs) {
		num, err := strconv.Atoi(strs[c])
		if err != nil {
			return "", errors.Wrapf(err, "Not a number %s", strs[c])
		}

		if num < len(vocab) {
			w = vocab[num]
		} else {
			w = ""
		}
		if w == "" {
			return "", errors.Errorf("Not found word by id %d", num)
		}
		strs[c] = w
		return toString(strs), nil
	}
	return line, nil
}

func toString(strs []string) string {
	res := strings.Builder{}
	for _, s := range strs {
		if res.Len() > 0 {
			res.WriteString(" ")
		}
		res.WriteString(s)
	}
	return res.String()
}
