package punctuation

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

//Wrapper is punctuate url wrapper
type Wrapper struct {
	url string
}

//NewWrapper creates instance
func NewWrapper(url string) *Wrapper {
	return &Wrapper{url: url}
}

//Punctuate invokes word punctuate service
func (p *Wrapper) Punctuate(words []string) (*Response, error) {
	inp := RequestArray{Words: words}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(inp)
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
