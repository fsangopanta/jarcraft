package main

import (
	"fmt"
	"os"

	"jarcraft/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, "Error:", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
