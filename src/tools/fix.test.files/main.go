package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
)

func main() {
	inPtr := flag.String("i", "", "dir in")
	outPtr := flag.String("o", "", "dir out")
	flag.Parse()

	files, err := os.ReadDir(*inPtr)
	if err != nil {
		panic(err)
	}

	i := 401
	var name string
	for _, f := range files {
		name = fmt.Sprintf("t_aft_%d.wav", i)

		source, err := os.Open(path.Join(*inPtr, f.Name()))
		if err != nil {
			panic(err)
		}
		defer source.Close()

		destination, err := os.Create(path.Join(*outPtr, name))
		if err != nil {
			panic(err)
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s %s\n", f.Name(), name)
		i++
	}
}
