package lattice

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//Word is one line in lattice file
type Word struct {
	Main  string
	From  string
	To    string
	Text  string
	Words []string // may be several words separated by '_'
	Punct string
}

//Part is some lattice data
type Part struct {
	Speaker string
	Num     int
	Words   []*Word
}

//SilWord indicate sil word in lattice
var SilWord = "<eps>"

//MainInd ord line indicator
var MainInd = "1"

//Read reads lattice file
func Read(src io.Reader) ([]*Part, error) {
	res := make([]*Part, 0)
	scanner := bufio.NewScanner(src)
	line := ""
	var cPart *Part
	var err error
	ln := 0
	for scanner.Scan() {
		ln++
		line = strings.TrimSpace(scanner.Text())
		strs := strings.Split(line, " ")
		if strings.HasPrefix(line, "#") {
			cPart = &Part{}
			if len(strs) < 3 {
				return nil, errors.Errorf("Line %d. Wrong part start: %s", ln, line)
			}
			cPart.Speaker = strs[2]
			cPart.Num, err = strconv.Atoi(strs[1])
			if err != nil {
				return nil, errors.Wrapf(err, "Line %d. Wrong part start: %s", ln, line)
			}
			res = append(res, cPart)
		} else if line == "" {
			cPart = nil
		} else {
			if cPart == nil {
				return nil, errors.Errorf("Line %d. No init part", ln)
			}
			if len(strs) < 4 {
				return nil, errors.Errorf("Line %d. Wrong line start: %s", ln, line)
			}
			w := &Word{}
			w.Main = strs[0]
			w.From = strs[1]
			w.To = strs[2]
			w.Text = strs[3]
			w.Words = splitWord(w.Text)
			if len(strs) > 4 {
				w.Punct = strs[4]
			}
			cPart.Words = append(cPart.Words, w)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

//Write writes lattice file
func Write(data []*Part, writer io.Writer) error {
	for _, p := range data {
		_, err := fmt.Fprintf(writer, "# %d %s\n", p.Num, p.Speaker)
		if err != nil {
			return err
		}
		for _, w := range p.Words {
			punct := ""
			if w.Punct != "" {
				punct = " " + w.Punct
			}
			_, err = fmt.Fprintf(writer, "%s %s %s %s%s\n", w.Main, w.From, w.To, getText(w.Words), punct)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(writer, "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

//WordDuration return duration from word line
func WordDuration(w *Word) time.Duration {
	return Duration(w.To) - Duration(w.From)
}

//Duration return duration from string as seconds (eg.: "2.50")
func Duration(str string) time.Duration {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return time.Duration(math.Round(res*1000)) * time.Millisecond
}

//SilDuration return duration between main last word of data[i-1] and first main word of data[i]
func SilDuration(data []*Part, i int) time.Duration {
	if i == 0 {
		return 0
	}
	tf := getLastWordDuration(data[i-1].Words)
	if tf == 0 {
		return 0
	}
	tt := getFirstWordDuration(data[i].Words)
	if tt == 0 {
		return 0
	}
	return tt - tf
}

func getLastWordDuration(data []*Word) time.Duration {
	for i := len(data) - 1; i >= 0; i-- {
		w := data[i]
		if w.Main == MainInd {
			if w.Text != SilWord {
				return Duration(w.To)
			}
		}
	}
	return 0
}

func getFirstWordDuration(data []*Word) time.Duration {
	for _, w := range data {
		if w.Main == MainInd {
			if w.Text != SilWord {
				return Duration(w.From)
			}
		}
	}
	return 0
}

func splitWord(s string) []string {
	if s == "_" {
		return []string{s}
	}
	return strings.Split(s, "_")
}

func getText(words []string) string {
	return strings.Join(words, "_")
}
