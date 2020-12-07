package main

import "fmt"

func copyAndExpand(node interface{}) (interface{}, error) {
	var rv interface{}
	if m, isMap := node.(map[interface{}]interface{}); isMap {
		tmp := map[string]interface{}{}
		for name, val := range m {
			nameString, ok := name.(string)
			if !ok {
				return nil, fmt.Errorf("cannot cast '%v' to a string", name)
			}
			var err error
			tmp[nameString], err = copyAndExpand(val)
			if err != nil {
				return nil, fmt.Errorf("cannot copyAndExpand(): %s", err)
			}
		}
		rv = tmp
	} else if a, isArray := node.([]interface{}); isArray {
		tmp := make([]interface{}, len(a))
		for i, val := range a {
			var err error
			tmp[i], err = copyAndExpand(val)
			if err != nil {
				return nil, fmt.Errorf("cannot copyAndExpand(): %s", err)
			}
		}
		rv = tmp
	} else {
		rv = node
	}
	return rv, nil
}
