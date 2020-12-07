package main

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestExpander(t *testing.T) {
	var tests = []struct {
		name string
		inp  string
		exp  string
	}{
		{
			"string", `
type: string
`,
			`
type: string
`,
		},
		{
			"simple object", `
type: object
properties:
  thing:
    type: string
`,
			`
type: object
properties:
  thing:
    type: string
`,
		},
		{
			"array of numbers",
			`
type: array
items:
  type: number
`,
			`
type: array
items:
  type: number
`,
		},
		{
			"array of objects",
			`
type: array
items:
  type: object
  properties:
    id:
      type: integer
    name:
      type: string
`,
			`
type: array
items:
  type: object
  properties:
    id:
      type: integer
    name:
      type: string
`,
		},
	}
	for i, tc := range tests {
		desc := fmt.Sprintf("Test Case %d: %s", i, tc.name)
		inp := tc.inp
		exp := tc.exp
		t.Run(desc, func(t *testing.T) {
			var thing map[interface{}]interface{}
			if err := yaml.Unmarshal([]byte(inp), &thing); err != nil {
				t.Fatalf("cannot unmarshal '%s'. Error: %s", inp, err)
			}
			got, err := copyAndExpand(thing)
			if err != nil {
				t.Fatalf("cannot copyAndExpand(): %s", err)
			}
			gotYaml, err := yaml.Marshal(got)
			if err != nil {
				t.Fatalf("cannot yaml.Marshal(): %s", err)
			}
			if string(gotYaml) != exp {
				t.Errorf("no match:\nGot: %s\nExp: %s", gotYaml, exp)
			}
		})
	}
}
