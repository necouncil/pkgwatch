package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type Result struct {
	Package string `json:"package"`
	Current string `json:"current"`
	Latest  string `json:"latest"`
}

func Print(results []Result) {
	for _, r := range results {
		fmt.Printf("%s %s -> %s\n", r.Package, r.Current, r.Latest)
	}
}

func PrintJSON(results []Result) {
	if results == nil {
		results = []Result{}
	}
	out := struct {
		Outdated []Result `json:"outdated"`
	}{Outdated: results}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}
