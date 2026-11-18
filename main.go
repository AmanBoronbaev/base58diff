package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const usage = "usage: base58diff encode HEX | decode BASE58 | encode-check HEX | decode-check BASE58 | check FILE|-"

func run(args []string, out, diagnostics io.Writer) int {
	if len(args) != 2 {
		fmt.Fprintln(diagnostics, usage)
		return 2
	}
	var value string
	var err error
	switch args[0] {
	case "check":
		var input io.Reader = os.Stdin
		if args[1] != "-" {
			file, err := os.Open(args[1])
			if err != nil {
				fmt.Fprintln(diagnostics, err)
				return 2
			}
			defer file.Close()
			input = file
		}
		count, code, err := checkCorpus(input)
		if err != nil {
			fmt.Fprintln(diagnostics, err)
			return code
		}
		fmt.Fprintf(out, "%d vectors passed\n", count)
		return 0
	case "encode", "encode-check":
		data, parseErr := hex.DecodeString(args[1])
		if parseErr != nil {
			fmt.Fprintln(diagnostics, parseErr)
			return 2
		}
		if args[0] == "encode-check" {
			if len(data) == 0 {
				fmt.Fprintln(diagnostics, "Base58Check requires a version byte")
				return 2
			}
			value, err = compareCheckEncode(data)
		} else {
			value, err = compareEncode(data)
		}
	case "decode-check":
		value, err = compareCheckDecode(args[1])
	case "decode":
		if args[1] == "" {
			fmt.Fprintln(diagnostics, "empty text is outside the shared decode contract")
			return 2
		}
		value, err = compareDecode(args[1])
	default:
		fmt.Fprintln(diagnostics, usage)
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
