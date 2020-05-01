package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

type request struct {
	ID       string `json:"id"`
	Timestap int64  `json:"timestamp"`
	Type     string `json:"type"`
	Worker   string `json:"worker"`
	Task     string `json:"task"`
}

func main() {
	log.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:<start | end> <worker> <task> <ID>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	mURL := os.Getenv("METRICS_URL")
	if mURL == "" {
		flag.Usage()
		log.Printf("Can't send metric: No METRICS_URL environment variable")
		return
	}
	if len(flag.Args()) != 4 {
		flag.Usage()
		log.Printf("Wrong params")
		return
	}
	data := &request{}
	data.Type = flag.Arg(0)
	data.Worker = flag.Arg(1)
	data.Task = flag.Arg(2)
	data.ID = flag.Arg(3)
	data.Timestap = time.Now().UnixNano()
	log.Printf("Sending metric: %s, %s, %s, %s, %d", data.Type, data.Worker, data.Task, data.ID, data.Timestap/1000000)
	err := post(data, mURL)
	if err != nil {
		log.Printf("Can't send metric: %s", err.Error())
	}
}

func post(data *request, url string) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)
	resp, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		return errors.Wrapf(err, "Can't invoke post to %s", url)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return errors.Errorf("Wrong response code from server. Code: %d", resp.StatusCode)
	}
	return nil
}
