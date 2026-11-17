package main

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"strings"
	"testing"
)

func TestVectors(t *testing.T) {
	for _, v := range []struct{ hex, text string }{
		{"", ""}, {"00", "1"}, {"0000", "11"}, {"0039", "1z"}, {"3a", "21"},
		{"61", "2g"}, {"626262", "a3gV"}, {"636363", "aPEr"},
	} {
		raw, err := hex.DecodeString(v.hex)
		if err != nil {
			t.Fatal(err)
		}
		text, err := compareEncode(raw)
		if err != nil || text != v.text {
			t.Errorf("%q: got %q, %v; want %q", v.hex, text, err, v.text)
		}
	}
}

func TestCLI(t *testing.T) {
	for _, v := range []struct {
		args []string
		code int
		want string
	}{
		{[]string{"encode", "000039"}, 0, "11z\n"},
		{[]string{"decode", "11z"}, 0, "000039\n"},
		{[]string{"decode", "0"}, 0, "rejected\n"},
		{[]string{"decode", ""}, 2, ""}, {[]string{"encode", "f"}, 2, ""},
	} {
		var out, diag bytes.Buffer
		code := run(v.args, &out, &diag)
		if code != v.code || out.String() != v.want || (code != 0) != (diag.Len() > 0) {
			t.Errorf("%q: code %d, stdout %q, stderr %q", v.args, code, out.String(), diag.String())
		}
	}
}

func TestCorpus(t *testing.T) {
	rng := rand.New(rand.NewSource(58))
	for _, n := range []int{1, 2, 7, 8, 9, 10, 11, 20, 25, 32, 33, 64, 78, 82, 128, 256} {
		for zeroes := 0; zeroes <= n; zeroes++ {
			raw := make([]byte, n)
			rng.Read(raw[zeroes:])
			if _, err := compareEncode(raw); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, bad := range []string{"0", "O", "I", "l", " ", "\x00", "é"} {
		for _, n := range []int{0, 1, 9, 10, 11, 20} {
			text := strings.Repeat("z", n) + bad + strings.Repeat("1", 20-n)
			got, err := compareDecode(text)
			if err != nil || got != "rejected" {
				t.Fatalf("%q: %q, %v", text, got, err)
			}
		}
	}
}

func FuzzEncode(f *testing.F) {
	for _, data := range [][]byte{{0}, {0, 0, 57}, {58}, bytes.Repeat([]byte{255}, 32)} {
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 256 {
			return
		}
		if _, err := compareEncode(data); err != nil {
			t.Fatal(err)
		}
	})
}
