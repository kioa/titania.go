package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

type ShowUnit struct {
	Name   string      `json:"unit"`
	Fruits []*ShowCode `json:"fruits"`
}

type ShowCode struct {
	Name     string      `json:"code"`
	Language string      `json:"language"`
	Details  []*ShowCase `json:"details"`
}

type ShowCase struct {
	Name   string `json:"case"`
	Result string `json:"result"`
	Time   string `json:"time"`
	OutPut string `json:"output"`
	Error  string `json:"error"`
}

// 流石に雑すぎる ちゃんと要約して
func WrapUp(outcomes []*ShowUnit) {
	for _, outcome := range outcomes {

		fmt.Fprintf(os.Stderr, "\n%s\n", pretty.Bold(pretty.Cyan(outcome.Name)))

		for _, fruit := range outcome.Fruits {

			fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Bold(fruit.Language), pretty.Bold(pretty.Blue(fruit.Name)))

			for _, detail := range fruit.Details {
				switch detail.Result {
				case "PASS":
					fmt.Fprintf(os.Stderr, "%s: %s %ss\n", pretty.Green(detail.Name), pretty.Green(detail.Result), detail.Time)
				case "FAIL":
					fmt.Fprintf(os.Stderr, "%s: %s %ss\n", pretty.Yellow(detail.Name), pretty.Yellow(detail.Result), detail.Time)
				case "CLIENT ERROR":
					fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Magenta(detail.Name), pretty.Magenta(detail.Result))
				case "SERVER ERROR":
					fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Blue(detail.Name), pretty.Blue(detail.Result))
				case "TESTER ERROR":
					fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Bold(pretty.Red(detail.Name)), pretty.Bold(pretty.Red(detail.Result)))
				default:
					fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Red(detail.Name), pretty.Red(detail.Result))
				}
			}
		}
	}
}
