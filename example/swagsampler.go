package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/breise/swagsampler"
)

const wordcharsRe = `[A-Za-z0-9_]{6,12}`
func main() {
	endpointP := flag.String("endpoint", "", "endpoint, including leading slash")
	methodP := flag.String("method", "", "method, in lowercase (e.g. get post)")
	wordcharsP := flag.Bool("wordchars", false, fmt.Sprintf("Generate strings where no pattern is specified with pattern %s", wordcharsRe))
	flag.Parse()
	endpoint := *endpointP
	method := *methodP
	wordchars := *wordcharsP
	if endpoint == "" || method == "" || len(flag.Args()) != 1 {
		flag.Usage()
		log.Fatal("usage: swagSampler -endpoint {endpoint} -method {method} {file path}")
	}
	file := flag.Arg(0)
	specBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("cannot open file '%s' for reading: %s", file, err)
	}

	sample, err := swagsampler.MkSample(specBytes, endpoint, method, wordchars)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		log.Fatalf("cannot marshal to json: %s", err)
	}
	fmt.Printf("%s\n", b)
}