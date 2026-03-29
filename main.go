package main

import (
	"os"

	"github.com/user/pkgwatch/cmd"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
