package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/breise/swagsampler"
)

func main() {
	endpointP := flag.String("endpoint", "", "endpoint, including leading slash")
	methodP := flag.String("method", "", "method, in lowercase (e.g. get post)")
	flag.Parse()
	endpoint := *endpointP
	method := *methodP
	if len(flag.Args()) != 1 {
		log.Fatal("usage: swagSampler -endpoint {endpoint} -method {method} {file path}")
	}
	file := flag.Arg(0)
	specBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("cannot open file '%s' for reading: %s", file, err)
	}

	sample, err := swagsampler.MkSample(specBytes, endpoint, method)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		log.Fatalf("cannot marshal to json: %s", err)
	}
	fmt.Printf("%s\n", b)
}