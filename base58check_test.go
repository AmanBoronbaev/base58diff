package main

import (
	"bytes"
	"testing"
)

func FuzzBase58Check(f *testing.F) {
	f.Add([]byte{0})
	f.Add([]byte{0xff, 0, 0})
	f.Add(append([]byte{0x04, 0x88, 0xb2, 0x1e}, make([]byte, 74)...))
	f.Add(bytes.Repeat([]byte{0}, 21))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 || len(data) > 256 {
			return
		}
		if _, err := compareCheckEncode(data); err != nil {
			t.Fatal(err)
		}
	})
}
