package swagsampler

import (
	"errors"
	"fmt"

	"github.com/breise/rstack"
	"github.com/breise/swagexpander"
	"gopkg.in/yaml.v2"
)

func MkSample(specBytes []byte, endpoint string, method string) (interface{}, error) {
	var spec interface{}
	if err := yaml.Unmarshal(specBytes, &spec); err != nil {
		return nil, fmt.Errorf("cannot unmarshal expanded spec: %s", err)
	}
	specExpanded, err := swagexpander.CopyAndExpand(spec)
	if err != nil {
		return nil, fmt.Errorf("cannot expand spec: %s", err)
	}
	specExpandedMap, ok := specExpanded.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("cannot cast expanded spec to a map[interface{}]interface{}")
	}

	paths, ok := specExpandedMap["paths"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("paths is not a map!: %s", paths)
	}
	endpointNode, ok := paths[endpoint].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("endpoint is not a map")
	}
	methodNode, ok := endpointNode[method].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("method is not a map")
	}
	parameters, ok := methodNode["parameters"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("parameters is not a map")
	}

	var bodyMap map[interface{}]interface{}
	for i, p := range parameters {
		m := p.(map[interface{}]interface{})
		in, ok := m["in"]
		if !ok {
			return nil, fmt.Errorf("parameter map #%d is missing the `in` key", i)
		}
		if in == "body" {
			if bodyMap != nil {
				return nil, fmt.Errorf("found more than one 'in: body' parameter map: %d", i)
			}
			bodyMap = m
		}
	}
	schema := bodyMap["schema"].(map[interface{}]interface{})
	breadcrumbs := rstack.New().Push("schema")
	rv, err := NewSampleBuilder().RenderSample(breadcrumbs, schema)
	if err != nil {
		return nil, fmt.Errorf("cannot mkSample() for %s: %s.\nNode: %v", "schema", err, schema)
	}
	return rv, nil
}
