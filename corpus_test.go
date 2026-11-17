package main

import (
	"strings"
	"testing"
)

func TestCheckCorpus(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		count, code int
		diagnostic  string
	}{
		{"shared empty encoding and mixed hex case", `[["",""],["3A","21"]]`, 2, 0, ""},
		{"wrong expectation", `[["00","1"],["0039","z"]]`, 1, 1, "vector 2"},
		{"bad hex", `[["GG","1"]]`, 0, 2, "vector 1"},
		{"null string", `[[null,""]]`, 0, 2, "vector 1"},
		{"extra column", `[["00","1","x"]]`, 0, 2, "vector 1"},
		{"truncated", `[["00","1"]`, 1, 2, ""},
		{"extra document", `[["00","1"]] true`, 1, 2, "after corpus"},
		{"empty", `[]`, 0, 2, "empty corpus"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			count, code, err := checkCorpus(strings.NewReader(tc.input))
			if count != tc.count || code != tc.code || (err == nil) != (code == 0) || (err != nil && !strings.Contains(err.Error(), tc.diagnostic)) {
				t.Fatalf("got count=%d code=%d err=%v", count, code, err)
			}
		})
	}
}
