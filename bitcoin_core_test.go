package main

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

func TestBitcoinCoreVectors(t *testing.T) {
	data, err := os.ReadFile("testdata/bitcoin-core-base58.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors [][]string
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	if len(vectors) == 0 {
		t.Fatal("empty Bitcoin Core corpus")
	}
	for i, pair := range vectors {
		if len(pair) != 2 {
			t.Fatalf("vector %d: expected a hex and a Base58 string", i)
		}
		raw, err := hex.DecodeString(pair[0])
		if err != nil {
			t.Fatalf("vector %d: %v", i, err)
		}
		encoded, err := compareEncode(raw)
		if err != nil || encoded != pair[1] {
			t.Errorf("vector %d: encode got %q, %v; want %q", i, encoded, err, pair[1])
		}
		if pair[1] == "" {
			continue
		} // The decoders differ on empty input.
		decoded, err := compareDecode(pair[1])
		if err != nil || decoded != hex.EncodeToString(raw) {
			t.Errorf("vector %d: decode got %q, %v; want %x", i, decoded, err, raw)
		}
	}
}
