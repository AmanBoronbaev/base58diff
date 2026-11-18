package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	btc "github.com/btcsuite/btcd/btcutil/base58"
	tron "github.com/mr-tron/base58/base58"
)

func checksum(data []byte) [4]byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return [4]byte(second[:4])
}

// data includes the version prefix. btcutil takes its first byte separately.
func compareCheckEncode(data []byte) (string, error) {
	sum := checksum(data)
	frame := make([]byte, 0, len(data)+len(sum))
	frame = append(frame, data...)
	frame = append(frame, sum[:]...)
	a, b := btc.CheckEncode(data[1:], data[0]), tron.Encode(frame)
	if a != b {
		return "", fmt.Errorf("Base58Check encode: btcutil=%q mr-tron=%q", a, b)
	}
	decoded, err := compareCheckDecode(a)
	if err != nil {
		return "", err
	}
	if decoded != hex.EncodeToString(data) {
		return "", fmt.Errorf("Base58Check round trip via %q returned %s", a, decoded)
	}
	return a, nil
}

func compareCheckDecode(text string) (string, error) {
	payload, version, btcErr := btc.CheckDecode(text)
	frame, tronErr := tron.Decode(text)
	valid := false
	var data []byte
	if tronErr == nil && len(frame) >= 5 {
		data = frame[:len(frame)-4]
		sum := checksum(data)
		valid = bytes.Equal(sum[:], frame[len(data):])
	}
	if valid != (btcErr == nil) {
		return "", fmt.Errorf("Base58Check decode %q: btcutil error=%v, reference accepts=%t", text, btcErr, valid)
	}
	if !valid {
		return "rejected", nil
	}
	if version != data[0] || !bytes.Equal(payload, data[1:]) {
		return "", fmt.Errorf("Base58Check decode %q: btcutil=%02x%x mr-tron=%x", text, version, payload, data)
	}
	return hex.EncodeToString(data), nil
}
