package punctuation

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

//Punctuator server call wrapper
type Punctuator interface {
	Punctuate(words []string) (*Response, error)
}
