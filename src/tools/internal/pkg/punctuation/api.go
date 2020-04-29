package punctuation

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

//Request is punctuation service json input
type Request struct {
	Text string `json:"text"`
}

//RequestArray is punctuation service json input
type RequestArray struct {
	Words []string `json:"input"`
}

//Response is json result
type Response struct {
	Original       []string `json:"original"`
	Punctuated     []string `json:"punctuated"`
	PunctuatedText string   `json:"punctuatedText"`
}

//Punctuator is punctuate url wrapper
type Punctuator struct {
	url string
}

//NewPunctuator creates instance
func NewPunctuator(url string) *Punctuator {
	return &Punctuator{url: url}
}

//Punctuate invokes word punctuate service
func (p *Punctuator) Punctuate(words []string) (*Response, error) {
	inp := RequestArray{Words: words}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(inp)
	resp, err := http.Post(p.url, "application/json; charset=utf-8", b)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't invoke post to %s", p.url)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return nil, errors.Errorf("Wrong response code from server. Code: %d", resp.StatusCode)
	}
	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, errors.Wrap(err, "Can't decode json")
	}
	return &res, nil
}
