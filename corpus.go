package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// checkCorpus consumes a JSON array of [hex, encoded text] pairs.
// It keeps one vector in memory and stops at the first failing row.
func checkCorpus(r io.Reader, base58Check bool) (count, code int, err error) {
	encode, decode := compareEncode, compareDecode
	if base58Check {
		encode, decode = compareCheckEncode, compareCheckDecode
	}
	d := json.NewDecoder(r)
	token, err := d.Token()
	if err != nil {
		return 0, 2, err
	}
	if token != json.Delim('[') {
		return 0, 2, fmt.Errorf("expected a JSON array")
	}
	for d.More() {
		var pair []*string
		if err := d.Decode(&pair); err != nil {
			return count, 2, fmt.Errorf("vector %d: %w", count+1, err)
		}
		if len(pair) != 2 || pair[0] == nil || pair[1] == nil {
			return count, 2, fmt.Errorf("vector %d: expected two strings", count+1)
		}
		raw, err := hex.DecodeString(*pair[0])
		if err != nil {
			return count, 2, fmt.Errorf("vector %d: %w", count+1, err)
		}
		if base58Check && len(raw) == 0 {
			return count, 2, fmt.Errorf("vector %d: Base58Check requires a version byte", count+1)
		}
		encoded, err := encode(raw)
		if err != nil {
			return count, 1, fmt.Errorf("vector %d: %w", count+1, err)
		}
		if encoded != *pair[1] {
			return count, 1, fmt.Errorf("vector %d: encode %s: got %q, want %q", count+1, *pair[0], encoded, *pair[1])
		}
		if encoded != "" {
			decoded, err := decode(*pair[1])
			if err != nil {
				return count, 1, fmt.Errorf("vector %d: %w", count+1, err)
			}
			if !strings.EqualFold(decoded, *pair[0]) {
				return count, 1, fmt.Errorf("vector %d: decode %q: got %s, want %s", count+1, encoded, decoded, *pair[0])
			}
		}
		count++
	}
	if _, err := d.Token(); err != nil {
		return count, 2, err
	}
	if _, err := d.Token(); err != io.EOF {
		return count, 2, fmt.Errorf("unexpected data after corpus")
	}
	if count == 0 {
		return 0, 2, fmt.Errorf("empty corpus")
	}
	return count, 0, nil
}
