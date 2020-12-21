package swagsampler

import (
	"errors"
	"fmt"
	"math"

	"github.com/breise/rstack"
	"github.com/breise/swagexpander"
	"gopkg.in/yaml.v2"
)

type SwagSampler struct {
	defaultExclusiveMaximum bool
	defaultExclusiveMinimum bool
	defaultMaximum          float64
	defaultMaxItems         int
	defaultMaxLength        int
	defaultMaxProperties    int
	defaultMinimum          float64
	defaultMinItems         int
	defaultMinLength        int
	defaultMinProperties    int
	defaultPattern          string
	defaultUniqueItems      bool
	useExample              bool
}

func New() *SwagSampler {
	return &SwagSampler{
		defaultMaximum:   float64(math.MaxInt32),
		defaultMinItems:  1,
		defaultMaxItems:  2,
		defaultMinLength: 6,
		defaultMaxLength: 16,
		defaultPattern:   `[A-Za-z0-9_]`,
		useExample:       false,
	}
}

func (s *SwagSampler) MkSample(specBytes []byte, endpoint string, method string) (interface{}, error) {
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
	rv, err := s.RenderSample(breadcrumbs, schema)
	if err != nil {
		return nil, fmt.Errorf("cannot RenderSample() for %s: %s.\nNode: %v", "schema", err, schema)
	}
	return rv, nil
}

func (s *SwagSampler) RenderSample(breadcrumbs *rstack.RStack, node map[interface{}]interface{}) (interface{}, error) {
	// fmt.Printf("breadcrumbs: %s, node: %v\n", breadcrumbs.Join(`/`), node)
	var rv interface{}
	tp, ok := node["type"]
	if !ok {
		return nil, fmt.Errorf("no 'type' key found in node for %s", breadcrumbs.Join(`/`))
	}

	pattern, havePattern := node["pattern"]
	enum, haveEnum := node["enum"]
	example, haveExample := node["example"]
	format, haveFormat := node["format"]
	// exclusiveMaximum, haveExclusiveMaximum := node["exclusiveMaximum"]
	// exclusiveMinimum, haveExclusiveMinimum := node["exclusiveMinimum"]
	// maximum, haveMaximum := node["maximum"]
	// maxItems, haveMaxItems := node["maxItems"]
	// maxLength, haveMaxLength := node["maxLength"]
	// maxProperties, haveMaxProperties := node["maxProperties"]
	// minimum, haveMinimum := node["minimum"]
	// minItems, haveMinItems := node["minItems"]
	// minLength, haveMinLength := node["minLength"]
	// minProperties, haveMinProperties := node["minProperties"]
	// uniqueItems, haveUniqueItems := node["uniqueItems"]

	if s.useExample && haveExample {
		rv = example
	} else if haveEnum {
		var err error
		rv, err = genEnum(enum)
		if err != nil {
			return nil, fmt.Errorf("in node for %s, cannot genEnum(): %s", breadcrumbs.Join(`/`), err)
		}
	} else if havePattern {
		var err error
		var pat string
		var ok bool
		if pat, ok = pattern.(string); !ok {
			return nil, fmt.Errorf("in node for %s, cannot cast pattern as a string: %v", breadcrumbs.Join(`/`), pattern)
		}
		rv, err = genSampleFromPattern(pat)
		if err != nil {
			return nil, fmt.Errorf("cannot RenderSample() in node for %s, type %s: %s", breadcrumbs.Join(`/`), tp, err)
		}
	} else if tp == "object" {
		tmp := map[string]interface{}{}
		properties, ok := node["properties"]
		if !ok {
			return nil, fmt.Errorf("object missing 'properties' key in node for %s", breadcrumbs.Join(`/`))
		}
		propertiesMap, ok := properties.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("properties is not a map in node for %s!", breadcrumbs.Join(`/`))
		}
		for name, node := range propertiesMap {
			var err error
			nameString, ok := name.(string)
			if !ok {
				return nil, fmt.Errorf("cannot cast %s as a string", name)
			}
			tmp[nameString], err = s.RenderSample(breadcrumbs.Push(name), node.(map[interface{}]interface{}))
			if err != nil {
				return nil, fmt.Errorf("cannot RenderSample() for %s in node for %s: %s", name, breadcrumbs.Join(`/`), err)
			}
		}
		rv = tmp
	} else if tp == "array" {
		// TODO: minItems, maxItems, uniqueItems(?)
		items, ok := node["items"]
		if !ok {
			return nil, fmt.Errorf("array missing 'items' key in node for %s", breadcrumbs.Join(`/`))
		}
		itemsMap, ok := items.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("items is not a map in node for %s", breadcrumbs.Join(`/`))
		}
		thing, err := s.RenderSample(breadcrumbs.Push("0"), itemsMap)
		if err != nil {
			return nil, fmt.Errorf("cannot RenderSample() in node for %s: %s", breadcrumbs.Join(`/`), err)
		}
		rv = []interface{}{thing}
	} else if tp == "integer" {
		var err error
		if haveFormat {
			switch format {
			case "int32":
				rv, err = s.genInt32(node)
			case "int64":
				rv, err = s.genInt64(node)
			default:
				rv, err = s.genInt(node)
			}
		} else {
			rv, err = s.genInt(node)
		}
		if err != nil {
			return nil, fmt.Errorf("cannot RenderSample() in node for %s, type %s: %s", breadcrumbs.Join(`/`), tp, err)
		}
	} else if tp == "number" {
		var err error
		if haveFormat {
			switch format {
			case "float":
				rv, err = s.genFloat32(node)
			case "double":
				rv, err = s.genFloat64(node)
			default:
				rv, err = s.genFloat32(node)
			}
		} else {
			rv, err = s.genFloat32(node)
		}
		if err != nil {
			return nil, fmt.Errorf("cannot RenderSample() in node for %s, type %s: %s", breadcrumbs.Join(`/`), tp, err)
		}
	} else if tp == "string" {
		var err error
		if haveFormat {
			switch format {
			// case "byte":
			// 	rv, err = genBase64(node)
			// case "binary":
			// 	rv, err = genOctets(node)
			// case "date":
			// 	rv, err = genDate()
			// case "date-time":
			// 	rv, err = genDateTime()
			// case "password":
			// 	rv, err = genPassword()
			default:
				rv, err = s.genString(node)
			}
		} else {
			rv, err = s.genString(node)
		}
		if err != nil {
			return nil, fmt.Errorf("cannot RenderSample() in node for %s, type %s: %s", breadcrumbs.Join(`/`), tp, err)
		}
	} else if tp == "boolean" {
		var err error
		rv = genBool()
		if err != nil {
			return nil, fmt.Errorf("cannot RenderSample() in node for %s, type %s: %s", breadcrumbs.Join(`/`), tp, err)
		}
	}
	return rv, nil
	// type		format
	// ----		------
	// boolean
	// integer	int32	//signed 32 bits
	// integer	int64	//signed 64 bits
	// number	float
	// number	double
	// string
	// string	byte	//base64 encoded characters
	// string	binary	//any sequence of octets
	// string	date	//As defined by full-date - RFC3339
	// string	date-time	// As defined by date-time - RFC3339
	// string	password	// Used to hint UIs the input needs to be obscured.
}

// func (s *SampleBuilder) nextInt() int {
// 	s.nextIntVal++
// 	return s.nextIntVal
// }

// func (s *SampleBuilder) nextFloat() float64 {
// 	s.nextIntVal++
// 	whole := float64(s.nextIntVal)
// 	s.nextIntVal++
// 	fraction := float64(s.nextIntVal)
// 	places := int(math.Floor(math.Log10(fraction)) + 1)
// 	decimal := fraction / math.Pow10(places)
// 	rv := whole + decimal
// 	return rv
// }

// func (s *SampleBuilder) nextString() string {
// 	s.nextIntVal++
// 	return fmt.Sprintf("string%04d", s.nextIntVal)
// }
