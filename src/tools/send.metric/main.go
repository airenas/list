package main

import (
	"bytes"
	"context"
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
		log.Print(err.Error())
		fs.Usage()
		return
	}
	req := mapRequest(data, time.Now())
	log.Printf("Sending metric: %s, %s, %s, %s, %d", req.Type, req.Worker, req.Task, req.ID, req.Timestap/1000000)
	err = post(req, data.url)
	if err != nil {
		log.Printf("can't send metric: %s", err.Error())
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
		return errors.New("no URL")
	}
	if data.worker == "" {
		return errors.New("no Worker")
	}
	if data.task == "" {
		return errors.New("no Task")
	}
	if data.id == "" {
		return errors.New("no ID")
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
	ctx, cf := context.WithTimeout(context.Background(), time.Second*10)
	defer cf()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, b)
	if err != nil {
		return errors.Wrapf(err, "can't prepare request", url)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "can't invoke post to %s", url)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return errors.Errorf("wrong response code from server, code: %d", resp.StatusCode)
	}
	return nil
}
