package main

import (
	"fmt"
	"log"
	"math"

	"github.com/breise/rstack"
)

type SampleBuilder struct {
	nextIntVal int
}

func NewSampleBuilder() *SampleBuilder {
	return &SampleBuilder{}
}

func (s *SampleBuilder) renderSample(rs *rstack.RStack, node map[interface{}]interface{}) (interface{}, error) {
	var rv interface{}
	tp, ok := node["type"]
	if !ok {
		return nil, fmt.Errorf("no 'type' key found in node for %s", rs.Join(`/`))
	}
	example, haveExample := node["example"]
	enum, haveEnum := node["enum"]
	if haveExample {
		rv = example
	} else if haveEnum {
		enumSlice := enum.([]interface{})
		rv = enumSlice[0]
	} else if tp == "object" {
		tmp := map[string]interface{}{}
		properties, ok := node["properties"]
		if !ok {
			log.Fatalf("object missing 'properties' key in node for %s", rs.Join(`/`))
		}
		propertiesMap, ok := properties.(map[interface{}]interface{})
		if !ok {
			log.Fatalf("properties is not a map in node for %s!", rs.Join(`/`))
		}
		for name, node := range propertiesMap {
			var err error
			nameString, ok := name.(string)
			if !ok {
				log.Fatalf("cannot cast %s as a string", name)
			}
			tmp[nameString], err = s.renderSample(rs.Push(name), node.(map[interface{}]interface{}))
			if err != nil {
				log.Fatalf("cannot mkSample() for %s in node for %s: %s", name, rs.Join(`/`), err)
			}
		}
		rv = tmp
	} else if tp == "array" {
		items, ok := node["items"]
		if !ok {
			log.Fatalf("array missing 'items' key in node for %s", rs.Join(`/`))
		}
		itemsMap, ok := items.(map[interface{}]interface{})
		if !ok {
			log.Fatalf("items is not a map in node for %s", rs.Join(`/`))
		}
		thing, err := s.renderSample(rs.Push("0"), itemsMap)
		if err != nil {
			return nil, fmt.Errorf("cannot mkSample() in node for %s", rs.Join(`/`))
		}
		rv = []interface{}{thing}
	} else if tp == "integer" {
		rv = s.nextInt()
	} else if tp == "number" {
		rv = fmt.Sprintf("%0.2f", s.nextFloat())
	} else if tp == "string" {
		rv = s.nextString()
	} else if tp == "boolean" {
		rv = "false"
	}
	return rv, nil

}

func (s *SampleBuilder) nextInt() int {
	s.nextIntVal++
	return s.nextIntVal
}

func (s *SampleBuilder) nextFloat() float64 {
	s.nextIntVal++
	whole := float64(s.nextIntVal)
	s.nextIntVal++
	fraction := float64(s.nextIntVal)
	places := int(math.Floor(math.Log10(fraction)) + 1)
	decimal := fraction / math.Pow10(places)
	rv := whole + decimal
	return rv
}

func (s *SampleBuilder) nextString() string {
	s.nextIntVal++
	return fmt.Sprintf("string%04d", s.nextIntVal)
}
