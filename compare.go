package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	btc "github.com/btcsuite/btcd/btcutil/base58"
	tron "github.com/mr-tron/base58/base58"
)

// compareDecode checks the two APIs on their common domain: nonempty text.
// btcutil reports invalid input as zero bytes; every valid nonempty Base58
// string decodes to at least one byte, including strings made only of '1'.
func compareDecode(text string) (string, error) {
	a := btc.Decode(text)
	b, err := tron.Decode(text)
	if (len(a) == 0) != (err != nil) || (err == nil && !bytes.Equal(a, b)) {
		return "", fmt.Errorf("decode %q: btcutil=%x mr-tron=%x err=%v", text, a, b, err)
	}
	if err != nil {
		return "rejected", nil
	}
	if btc.Encode(a) != text || tron.Encode(b) != text {
		return "", fmt.Errorf("decode %q: canonical re-encoding differs", text)
	}
	return hex.EncodeToString(a), nil
}

func compareEncode(data []byte) (string, error) {
	a, b := btc.Encode(data), tron.Encode(data)
	if a != b {
		return "", fmt.Errorf("encode %x: btcutil=%q mr-tron=%q", data, a, b)
	}
	if len(data) > 0 {
		got, err := compareDecode(a)
		if err != nil {
			return "", err
		}
		if got != hex.EncodeToString(data) {
			return "", fmt.Errorf("round trip %x via %q returned %s", data, a, got)
		}
	}
	return a, nil
}
