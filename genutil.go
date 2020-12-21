package swagsampler

import (
	"fmt"
	"math/rand"
	"time"

	goregen "github.com/zach-klippenstein/goregen"
)

var asciiMinChar = int32(' ')
var asciiMaxChar = int32('~')
var asciiRangeSize = 1 + asciiMaxChar - asciiMinChar

var randSource = rand.NewSource(time.Now().UnixNano())

var randgen = rand.New(randSource)

var regenArgs = &goregen.GeneratorArgs{
	RngSource: randSource,
}

func genSampleFromPattern(pattern string) (string, error) {
	generator, err := goregen.NewGenerator(pattern, regenArgs)
	if err != nil {
		return "", fmt.Errorf("cannot goregen.NewGenerator(): %s", err)
	}
	return generator.Generate(), nil
}

func (s *SwagSampler) genString(node map[interface{}]interface{}) (string, error) {
	maxLength, haveMaxLength := node["maxLength"]
	minLength, haveMinLength := node["minLength"]

	var minLen int = s.defaultMinLength
	var maxLen int = s.defaultMaxLength
	if haveMinLength {
		var ok bool
		if minLen, ok = minLength.(int); !ok {
			return "", fmt.Errorf("genString(): cannot cast minLength to an int: '%v' of type %T", minLength, minLength)
		}
	}
	if haveMaxLength {
		var ok bool
		if maxLen, ok = maxLength.(int); !ok {
			return "", fmt.Errorf("genString(): cannot cast maxLength to an int: '%v' of type %T", maxLength, maxLength)
		}
	}
	if s.defaultPattern != "" {
		// pat := fmt.Sprintf("%s{%d,%d}", s.defaultPattern, minLen, maxLen)
		return genSampleFromPattern(s.defaultPattern)
	}
	// determine length
	rangeSize := 1 + maxLen - minLen
	len := randgen.Int31n(int32(rangeSize)) + int32(minLen)
	// make a string of len pseudorandom printable ascii characters
	rv := make([]byte, len)
	for i := int32(0); i < len; i++ {
		rv[i] = byte(randgen.Int31n(asciiRangeSize) + asciiMinChar)
	}
	return string(rv), nil
}

func genBool() bool {
	n := randgen.Int31n(2) // 0 or 1
	return n > 0
}

func (s *SwagSampler) genInt(node map[interface{}]interface{}) (int, error) {
	exclusiveMaximum, haveExclusiveMaximum := node["exclusiveMaximum"]
	exclusiveMinimum, haveExclusiveMinimum := node["exclusiveMinimum"]
	maximum, haveMaximum := node["maximum"]
	minimum, haveMinimum := node["minimum"]
	min := int(s.defaultMinimum)
	max := int(s.defaultMaximum)
	if haveMinimum {
		var ok bool
		if min, ok = minimum.(int); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast minimum to an int: '%v' of type %T", minimum, minimum)
		}
	}
	if haveMaximum {
		var ok bool
		if max, ok = maximum.(int); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast maximum to an int: '%v' of type %T", maximum, maximum)
		}
	}
	if haveExclusiveMinimum {
		var xm bool
		var ok bool
		if xm, ok = exclusiveMinimum.(bool); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast exclusiveMinimum to a bool: '%v' of type %T", exclusiveMinimum, exclusiveMinimum)
		}
		if xm {
			min++
		}
	}
	if haveExclusiveMaximum {
		var xm bool
		var ok bool
		if xm, ok = exclusiveMaximum.(bool); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast exclusiveMaximum to a bool: '%v' of type %T", exclusiveMaximum, exclusiveMaximum)
		}
		if xm {
			max++
		}
	}
	return int(randgen.Int31n(int32(max))) + min, nil
}

// Being lazy
func (s *SwagSampler) genInt32(node map[interface{}]interface{}) (int32, error) {
	x, err := s.genInt(node)
	return int32(x), err
}

// Being lazy
func (s *SwagSampler) genInt64(node map[interface{}]interface{}) (int64, error) {
	x, err := s.genInt(node)
	return int64(x), err
}

func (s *SwagSampler) genFloat32(node map[interface{}]interface{}) (float32, error) {
	exclusiveMaximum, haveExclusiveMaximum := node["exclusiveMaximum"]
	exclusiveMinimum, haveExclusiveMinimum := node["exclusiveMinimum"]
	maximum, haveMaximum := node["maximum"]
	minimum, haveMinimum := node["minimum"]
	var min = float32(s.defaultMinimum)
	var max = float32(s.defaultMaximum)
	if haveMinimum {
		var ok bool
		if min, ok = minimum.(float32); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast minimum to an float32: '%v' of type %T", minimum, minimum)
		}
	}
	if haveMaximum {
		var ok bool
		if max, ok = maximum.(float32); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast maximum to an float32: '%v' of type %T", maximum, maximum)
		}
	}
	if haveExclusiveMinimum {
		var xm bool
		var ok bool
		if xm, ok = exclusiveMinimum.(bool); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast exclusiveMinimum to a bool: '%v' of type %T", exclusiveMinimum, exclusiveMinimum)
		}
		if xm {
			min++
		}
	}
	if haveExclusiveMaximum {
		var xm bool
		var ok bool
		if xm, ok = exclusiveMaximum.(bool); !ok {
			return 0, fmt.Errorf("getInt(): cannot cast exclusiveMaximum to a bool: '%v' of type %T", exclusiveMaximum, exclusiveMaximum)
		}
		if xm {
			max++
		}
	}
	return randgen.Float32()*(max-min) + min, nil
}

// Being lazy
func (s *SwagSampler) genFloat64(node map[interface{}]interface{}) (float64, error) {
	x, err := s.genFloat32(node)
	return float64(x), err
}

func genEnum(enum interface{}) (interface{}, error) {
	enumSlice, ok := enum.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot cast enum to slice of interface{}: %v", enum)
	}
	n := len(enumSlice)
	i := randgen.Int31n(int32(n))
	return enumSlice[i], nil
}
