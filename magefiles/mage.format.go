// nolint
package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Lint Runs golangci-lint checks over the code.
func Lint() error {
	out, err := sh.Output(
		"golangci-lint",
		"run", "./...",
		"--allow-parallel-runners",
		"--skip-dirs", `(node_modules|magefiles|\.serverless|mod|bin|vendor|\.github|\.git)`,
	)

	fmt.Println(out)
	return err
}

// Format Runs gofmt over the code.
func Format() error {
	out, err := sh.Output("goimports-reviser", "-format", "./...")
	fmt.Println(out)

	return err
}
