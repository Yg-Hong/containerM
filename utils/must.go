package utils

import (
	"fmt"
	"os"
)

func Must(err error, context ...string) {
	if err != nil {
		if len(context) > 0 {
			fmt.Fprintf(os.Stderr, "Error (%s): %v\n", context[0], err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
