package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func run(args []string, out, diagnostics io.Writer) int {
	if len(args) != 2 {
		fmt.Fprintln(diagnostics, "usage: base58diff encode HEX | decode BASE58")
		return 2
	}
	var value string
	var err error
	switch args[0] {
	case "encode":
		data, parseErr := hex.DecodeString(args[1])
		if parseErr != nil {
			fmt.Fprintln(diagnostics, parseErr)
			return 2
		}
		value, err = compareEncode(data)
	case "decode":
		if args[1] == "" {
			fmt.Fprintln(diagnostics, "empty text is outside the shared decode contract")
			return 2
		}
		value, err = compareDecode(args[1])
	default:
		fmt.Fprintln(diagnostics, "expected encode or decode")
		return 2
	}
	if err != nil {
		fmt.Fprintln(diagnostics, err)
		return 1
	}
	fmt.Fprintln(out, value)
	return 0
}

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
