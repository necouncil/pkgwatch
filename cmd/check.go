package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/pkgwatch/internal/output"
	"github.com/user/pkgwatch/internal/pkgbuild"
	"github.com/user/pkgwatch/internal/upstream"
)

func Run(args []string) error {
	jsonMode := false
	var dirs []string

	for _, a := range args {
		if a == "--json" {
			jsonMode = true
		} else {
			dirs = append(dirs, a)
		}
	}

	if len(dirs) == 0 {
		entries, err := os.ReadDir(".")
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() {
				dirs = append(dirs, e.Name())
			}
		}
		if _, err := os.Stat("PKGBUILD"); err == nil {
			dirs = []string{"."}
		}
	}

	var results []output.Result

	for _, dir := range dirs {
		path := filepath.Join(dir, "PKGBUILD")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = dir
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}
		}

		pkg, err := pkgbuild.Parse(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, err)
			continue
		}

		latest, err := upstream.LatestVersion(pkg.Source)
		if err != nil || latest == "" {
			continue
		}

		if upstream.IsNewer(latest, pkg.Version) {
			results = append(results, output.Result{
				Package: pkg.Name,
				Current: pkg.Version,
				Latest:  latest,
			})
		}
	}

	if jsonMode {
		output.PrintJSON(results)
	} else {
		output.Print(results)
	}

	return nil
}
