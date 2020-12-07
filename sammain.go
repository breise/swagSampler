package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/breise/rstack"
)

func main() {
	endpointP := flag.String("endpoint", "", "endpoint, including leading slash")
	methodP := flag.String("method", "", "method, in lowercase (e.g. get post)")
	flag.Parse()
	endpoint := *endpointP
	method := *methodP
	if len(flag.Args()) != 1 {
		log.Fatal("usage: sam -endpoint {endpoint} -method {method} {file path}")
	}
	file := flag.Arg(0)
	spec, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("cannot open file '%s' for reading: %s", file, err)
	}
	sample := mkSample(spec, endpoint, method)
	b, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		log.Fatalf("cannot marshal to json: %s", err)
	}
	fmt.Printf("%s\n", b)
}

func mkSample(spec []byte, endpoint string, method string) interface{} {
	var thing map[interface{}]interface{}
	if err := yaml.Unmarshal(spec, &thing); err != nil {
		log.Fatalf("cannot unmarshal spec: %s", err)
	}
	paths, ok := thing["paths"].(map[interface{}]interface{})
	if !ok {
		log.Fatalf("paths is not a map!: %s", paths)
	}
	endpointNode, ok := paths[endpoint].(map[interface{}]interface{})
	if !ok {
		log.Fatalf("endpoint is not a map!")
	}
	methodNode, ok := endpointNode[method].(map[interface{}]interface{})
	if !ok {
		log.Fatalf("method is not a map!")
	}
	parameters, ok := methodNode["parameters"].([]interface{})
	if !ok {
		log.Fatalf("parameters is not a map!")
	}

	var bodyMap map[interface{}]interface{}
	for i, p := range parameters {
		m := p.(map[interface{}]interface{})
		in, ok := m["in"]
		if !ok {
			log.Fatalf("parameter map #%d is missing the `in` key", i)
		}
		if in == "body" {
			if bodyMap != nil {
				log.Fatalf("found more than one 'in: body' parameter map!: %d", i)
			}
			bodyMap = m
		}
	}
	schema := bodyMap["schema"].(map[interface{}]interface{})
	rs := rstack.New().Push("schema")
	sample, err := NewSampleBuilder().renderSample(rs, schema)
	if err != nil {
		log.Fatalf("cannot mkSample() for %s: %s", "schema", err)
	}
	return sample
}
