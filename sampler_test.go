package main

import (
	"fmt"
	"testing"

	"github.com/breise/rstack"
	"gopkg.in/yaml.v2"
)

func TestSam(t *testing.T) {
	var tests = []struct {
		name string
		inp  string
		exp  string
	}{
		{
			"string", `
type: string
`,
			"string0001",
		},
		{
			"simple object", `
type: object
properties:
  thing:
    type: string
`,
			"map[thing:string0001]",
		},
		{
			"array of numbers",
			`
type: array
items:
  type: number
`,
			"[1.20]",
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
			"[map[id:1 name:string0002]]",
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
			got, err := NewSampleBuilder().renderSample(rstack.New().Push("Top"), thing)
			if err != nil {
				t.Fatal("cannot mkSample(Top)")
			}
			gotString := fmt.Sprintf("%+v", got)
			if gotString != exp {
				t.Errorf("no match:\nGot: %s\nExp: %s", gotString, exp)
			}
		})
	}
}
