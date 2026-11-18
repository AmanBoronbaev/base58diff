package main

import (
	"os"
	"path/filepath"
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
			count, code, err := checkCorpus(strings.NewReader(tc.input), false)
			if count != tc.count || code != tc.code || (err == nil) != (code == 0) || (err != nil && !strings.Contains(err.Error(), tc.diagnostic)) {
				t.Fatalf("got count=%d code=%d err=%v", count, code, err)
			}
		})
	}
}

func TestCheckBase58CheckCLI(t *testing.T) {
	for _, tc := range []struct {
		file string
		data string
		code int
		want string
	}{
		{"testdata/bip32-xpub.json", "", 0, "17 vectors passed"},
		{"", `[["00","1Wh4bh"],["ff","wrong"]]`, 1, "vector 2"},
		{"", `[["",""]]`, 2, "vector 1"},
	} {
		path := tc.file
		if path == "" {
			path = filepath.Join(t.TempDir(), "vectors.json")
			if err := os.WriteFile(path, []byte(tc.data), 0600); err != nil {
				t.Fatal(err)
			}
		}
		var out, diag strings.Builder
		code := run([]string{"check", "--base58check", path}, &out, &diag)
		message := diag.String()
		if tc.code == 0 {
			message = out.String()
		}
		if code != tc.code || !strings.Contains(message, tc.want) {
			t.Errorf("check %q: exit %d, stdout %q, stderr %q", path, code, out.String(), diag.String())
		}
	}
}
