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
	Model    string `json:"model"`
	Task     string `json:"task"`
}

type params struct {
	url    string
	start  bool
	id     string
	worker string
	model  string
	task   string
}

func main() {
	log.SetOutput(os.Stderr)
	data := &params{}
	takeEnvValues(data)
	fs := flag.CommandLine
	takeParams(fs, data)
	fs.Parse(os.Args[1:])
	err := validateParams(data)
	if err != nil {
		log.Printf(err.Error())
		fs.Usage()
		return
	}
	req := mapRequest(data, time.Now())
	log.Printf("Sending metric: %s, %s, %s, %s, %d", req.Type, req.Worker, req.Task, req.ID, req.Timestap/1000000)
	err = post(req, data.url)
	if err != nil {
		log.Printf("Can't send metric: %s", err.Error())
	}
}

func takeEnvValues(data *params) {
	data.url = os.Getenv("METRICS_URL")
	data.model = os.Getenv("MODEL")
}

func takeParams(fs *flag.FlagSet, data *params) {
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s: <params>\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.StringVar(&data.id, "i", "", "Transcription ID")
	fs.StringVar(&data.worker, "w", "", "Worker")
	fs.StringVar(&data.task, "t", "", "Task")
	fs.StringVar(&data.model, "m", data.model, "Model")
	fs.StringVar(&data.url, "u", data.url, "Metrics URL")
	fs.BoolVar(&data.start, "s", false, "Start if set, else end")
}

func validateParams(data *params) error {
	if data.url == "" {
		return errors.New("No URL")
	}
	if data.model == "" {
		return errors.New("No Model")
	}
	if data.id == "" {
		return errors.New("No ID")
	}
	if data.task == "" {
		return errors.New("No Task")
	}
	if data.worker == "" {
		return errors.New("No Worker")
	}
	return nil
}

func mapRequest(data *params, t time.Time) *request {
	req := &request{}
	req.ID = data.id
	req.Model = data.model
	req.Task = data.task
	req.Worker = data.worker
	req.Timestap = t.UnixNano()
	req.Type = "end"
	if data.start {
		req.Type = "start"
	}
	return req
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
