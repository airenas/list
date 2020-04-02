package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type workingData struct {
	inPipe  string
	outPipe string
	timeout time.Duration
	args    []string
}

func main() {
	inputPipeStrPtr := flag.String("i", "", "Input pipe name")
	outputPipeStrPtr := flag.String("o", "", "Output pipe name")
	timeoutPtr := flag.String("t", "3m", "Timeout to read initial value from pipe (eg. 30s or 2m)")
	flag.Parse()

	if strings.TrimSpace(*inputPipeStrPtr) == "" || strings.TrimSpace(*outputPipeStrPtr) == "" ||
		len(flag.Args()) < 2 {
		panic(errors.New("Usage: pipe.runner -i <inputPipe> -o <outputPipe> -t <timeoutSec> <args to real command>\n\t"))
	}
	var err error
	data := workingData{}
	data.inPipe = *inputPipeStrPtr
	data.outPipe = *outputPipeStrPtr
	data.timeout, err = time.ParseDuration(*timeoutPtr)
	if err != nil {
		panic(errors.Errorf("Can't parse duration '%s' (eg. 30s or 2m)", *timeoutPtr))
	}
	data.args = flag.Args()

	fmt.Fprint(os.Stdout, "Starting pipe runner\n")
	fmt.Fprintf(os.Stdout, "In: 		%s\n", data.inPipe)
	fmt.Fprintf(os.Stdout, "Out: 		%s\n", data.outPipe)
	fmt.Fprintf(os.Stdout, "Timeout: 	%v\n", data.timeout)
	fmt.Fprintf(os.Stdout, "Args: 		%v\n", data.args)

	fmt.Fprint(os.Stdout, "Running\n")
	err = run(&data)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(os.Stdout, "Done\n")
	os.Exit(0)
}

type fOpenResult struct {
	file *os.File
	err  error
}

func openPipe(name string, mode int, timeout time.Duration) (*os.File, error) {
	fmt.Fprintf(os.Stdout, "Opening pipe %s\n", name)
	c := make(chan fOpenResult, 1)
	go func() {
		f, err := os.OpenFile(name, mode, os.ModeNamedPipe)
		if err != nil {
			c <- fOpenResult{file: nil, err: err}
		} else {
			c <- fOpenResult{file: f, err: nil}
		}
	}()

	select {
	case res := <-c:
		return res.file, res.err
	case <-time.After(timeout):
		return nil, errors.Errorf("Timeout opening pipe %s", name)
	}
}

func writeArgs(data *workingData) error {
	pIn, err := openPipe(data.inPipe, os.O_WRONLY, data.timeout)
	if err != nil {
		return errors.Wrapf(err, "Can't open pipe '%s'", data.inPipe)
	}
	defer pIn.Close()

	_, err = pIn.WriteString(getWriteData(data))
	if err != nil {
		return errors.Wrapf(err, "Can't write to pipe '%s'", data.outPipe)
	}
	return nil
}

func readPipe(data *workingData) error {
	pOut, err := openPipe(data.outPipe, os.O_RDONLY, data.timeout)
	if err != nil {
		return errors.Wrapf(err, "Can't open pipe '%s'", data.outPipe)
	}
	defer pOut.Close()

	scanner := bufio.NewScanner(pOut)
	line := ""
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		fmt.Fprintf(os.Stdout, "%s\n", line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if line != "0" {
		return errors.Errorf("Program did not indicate '0' exit")
	}
	return nil
}

func run(data *workingData) error {
	err := writeArgs(data)
	if err != nil {
		return err
	}
	err = readPipe(data)
	if err != nil {
		return err
	}
	return nil
}

func getWriteData(data *workingData) string {
	res := ""
	sep := ""
	for _, s := range data.args {
		res = res + sep + s
		sep = " "
	}
	res = res + sep + data.outPipe
	return res
}
