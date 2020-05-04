package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTakeEnv(t *testing.T) {
	data := &params{}
	os.Setenv("MODEL", "model")
	os.Setenv("METRICS_URL", "url")
	takeEnvValues(data)
	assert.Equal(t, "url", data.url)
	assert.Equal(t, "model", data.model)
}

func TestParseParams(t *testing.T) {
	data := &params{}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	takeParams(fs, data)
	fs.Parse([]string{"-u", "url", "-m", "model", "-i", "id", "-t", "task", "-w", "worker", "-s"})
	assert.Equal(t, "url", data.url)
	assert.Equal(t, "model", data.model)
	assert.Equal(t, "id", data.id)
	assert.Equal(t, "task", data.task)
	assert.Equal(t, "worker", data.worker)
	assert.True(t, data.start)
}

func TestParseStart(t *testing.T) {
	data := &params{}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	takeParams(fs, data)
	fs.Parse([]string{"-s"})
	assert.True(t, data.start)
	fs = flag.NewFlagSet("test", flag.ContinueOnError)
	takeParams(fs, data)
	fs.Parse([]string{})
	assert.False(t, data.start)
}

func TestValidate(t *testing.T) {
	data := &params{url: "url", id: "id", model: "model", worker: "worker", task: "task"}
	err := validateParams(data)
	assert.Nil(t, err)
	data = &params{id: "id", model: "model", worker: "worker", task: "task"}
	err = validateParams(data)
	assert.NotNil(t, err)
	data = &params{url: "url", model: "model", worker: "worker", task: "task"}
	err = validateParams(data)
	assert.NotNil(t, err)
	data = &params{url: "url", id: "id", worker: "worker", task: "task"}
	err = validateParams(data)
	assert.NotNil(t, err)
	data = &params{url: "url", id: "id", model: "model", task: "task"}
	err = validateParams(data)
	assert.NotNil(t, err)
	data = &params{url: "url", id: "id", model: "model", worker: "worker"}
	err = validateParams(data)
	assert.NotNil(t, err)
}

func TestMap(t *testing.T) {
	data := &params{url: "url", id: "id", model: "model", worker: "worker", task: "task"}
	n := time.Now()
	req := mapRequest(data, n)
	assert.Equal(t, "id", req.ID)
	assert.Equal(t, "model", req.Model)
	assert.Equal(t, "worker", req.Worker)
	assert.Equal(t, "task", req.Task)
	assert.Equal(t, "end", req.Type)
	assert.Equal(t, n.UnixNano(), req.Timestap)
}

func TestMapStart(t *testing.T) {
	data := &params{start: true}
	n := time.Now()
	req := mapRequest(data, n)
	assert.Equal(t, "start", req.Type)
}
