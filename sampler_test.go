package swagsampler_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/breise/rstack"
	"github.com/breise/swagsampler"

	"gopkg.in/yaml.v2"
)

func TestSam(t *testing.T) {
	var tests = []struct {
		name string
		inp  string
		exp  func(interface{}) error
	}{
		{
			"string", `
type: string
`,
			func(x interface{}) error {
				var errs []string
				y := x.(string)
				if y == "" {
					errs = append(errs, "string is empty")
				}
				if len(y) > 255 {
					errs = append(errs, fmt.Sprintf("'%s' is > 255 characters", y))
				}
				if len(errs) > 0 {
					return errors.New(strings.Join(errs, "; "))
				}
				return nil
			},
		},
		{
			"string with min and max lengths", `
type: string
minLength: 10
maxLength: 20
`,
			func(x interface{}) error {
				var errs []string
				y := x.(string)
				if len(y) < 10 {
					errs = append(errs, fmt.Sprintf("'%s' is < 10 characters", y))
				}
				if len(y) > 20 {
					errs = append(errs, fmt.Sprintf("'%s' is > 20 characters", y))
				}
				if len(errs) > 0 {
					return errors.New(strings.Join(errs, "; "))
				}
				return nil
			},
		},
		// 		{
		// 			"simple object", `
		// type: object
		// properties:
		//   thing:
		//     type: string
		// `,
		// 			"map[thing:string0001]",
		// 		},
		// 		{
		// 			"array of numbers",
		// 			`
		// type: array
		// items:
		//   type: number
		// `,
		// 			"[1.20]",
		// 		},
		// 		{
		// 			"array of objects",
		// 			`
		// type: array
		// items:
		//   type: object
		//   properties:
		//     id:
		//       type: integer
		//     name:
		//       type: string
		// `,
		// 			"[map[id:1 name:string0002]]",
		// 		},
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
			got, err := swagsampler.New().RenderSample(rstack.New().Push("Top"), thing)
			if err != nil {
				t.Fatalf("cannot mkSample(Top): %s", err)
			}
			if err := exp(got); err != nil {
				t.Errorf("Test Case %d: %s: %s", i, desc, err)
			}
		})
	}
}
