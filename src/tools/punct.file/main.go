package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/airenas/list/src/tools/internal/pkg/punctuation"
	"github.com/pkg/errors"
)

func main() {
	filePtr := flag.String("f", "", "file in")
	urlPtr := flag.String("u", "", "punctuation URL")
	outPtr := flag.String("o", "", "file for output")
	flag.Parse()

	if *filePtr == "" || *urlPtr == "" || *outPtr == "" {
		panic(errors.New("Usage: ./punct.file -f <fileIn> -u <punctuation URL> -o <output file>"))
	}

	b, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		panic(errors.Wrapf(err, "Can't read file %s ", *filePtr))
	}

	destination, err := os.Create(*outPtr)
	if err != nil {
		panic(errors.Wrap(err, "Can't create file "+*outPtr))
	}
	defer destination.Close()

	result, err := punctuate(string(b), *urlPtr)
	if err != nil {
		panic(errors.Wrap(err, "Can't punctuate"))
	}
	destination.WriteString(result)
}

func punctuate(str string, url string) (string, error) {
	if strings.TrimSpace(str) == "" {
		return "", nil
	}
	inp := punctuation.Request{Text: str}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(inp)
	resp, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		return "", errors.Wrapf(err, "Can't invoke post to %s", url)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return "", errors.Errorf("Wrong response code from server. Code: %d", resp.StatusCode)
	}
	var res punctuation.Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", errors.Wrap(err, "Can't decode json")
	}
	return res.PunctuatedText, nil
}
