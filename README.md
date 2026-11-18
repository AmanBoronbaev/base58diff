# base58diff

Compare two Go implementations of Bitcoin's Base58 alphabet on the same input. Use it to replay a test vector, check a codec upgrade, or turn an encoding disagreement into a small Go regression test.

```sh
go run . encode 000039  # 11z
go run . decode 11z     # 000039
go run . decode 1O      # rejected
```

Each command runs both [btcutil](https://github.com/btcsuite/btcd/tree/master/btcutil) and [mr-tron](https://github.com/mr-tron/base58). Encoding checks the resulting text and decodes it back to the original bytes. Decoding compares acceptance, decoded bytes and canonical re-encoding. Leading zero bytes matter: removing one leading `1` changes the payload.

The program exits with 0 when both libraries agree, including when both reject a string. Exit 1 reports a disagreement or failed round trip on stderr. Exit 2 indicates a command or input error. Empty bytes encode to empty text. Empty text is excluded from decode comparisons because the libraries expose different empty-input behavior. The `encode` and `decode` commands operate on raw Base58.

## Compare Base58Check encodings

Pass the complete version-prefixed bytes in hex. For example, this Bitcoin address has version byte `00` and a 20-byte hash:

```sh
base58diff encode-check 00010966776006953d5567439e5e39f86a0d273bee
# 16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM
base58diff decode-check 16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM
# 00010966776006953d5567439e5e39f86a0d273bee
```

The checker compares btcutil's Base58Check API with mr-tron plus a four-byte double-SHA-256 checksum. Decoding validates the checksum and preserves every version byte. It prints `rejected` when both paths reject the input, including a changed checksum. The tool checks the encoding layer only. It accepts arbitrary version prefixes and payload lengths.

## Install

With Go 1.22 or later:

```sh
go install github.com/AmanBoronbaev/base58diff@v0.3.0
base58diff decode 11z
```

The command is installed in `GOBIN`, or in `$(go env GOPATH)/bin` by default.

## Replay a vector file

```sh
go run . check testdata/bitcoin-core-base58.json
# 21 vectors passed
printf '%s\n' '[["000039","11z"]]' | go run . check -
```

For Base58Check vectors, pass `--base58check`. The first column contains version-prefixed hex bytes with the checksum omitted. The second column contains the complete checksummed string.

```sh
go run . check --base58check testdata/bip32-xpub.json
# 17 vectors passed
```

The bundled corpus covers all 17 extended public keys in [BIP 32 test vectors 1–4](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#test-vectors). The expected 78-byte payloads were recovered from the published strings using integer base conversion and verified against their double-SHA-256 checksums. The vectors are distributed under the [BIP's license](testdata/COPYING.bip32).

Supply a JSON array of hex/Base58 string pairs. The checker streams one pair at a time and reports the first failing vector by its one-based index. A compiled binary exits with 1 for a mismatch and 2 for malformed or empty files.

## Run the corpus

```sh
go test ./...
go test -fuzz=FuzzEncode -fuzztime=30s -parallel=2
go test -fuzz=FuzzBase58Check -fuzztime=30s -parallel=2
```

The regular suite includes all 21 [Bitcoin Core v29.0 Base58 vectors](https://github.com/bitcoin/bitcoin/blob/v29.0/src/test/data/base58_encode_decode.json), stored under the [upstream MIT license](testdata/COPYING.bitcoin-core). It also covers fixed input/output pairs, a deterministic payload corpus with varying leading-zero prefixes, and invalid characters at decoder chunk boundaries. The fuzz targets vary binary payloads up to 256 bytes and check both directions, with and without a checksum. Go saves a minimized failure under `testdata/fuzz/` and prints the exact replay command. Commit a confirmed failure there to include it in subsequent ordinary test runs.

To compare a dependency upgrade, change its version in `go.mod`, run `go mod tidy` and rerun the corpus. A discrepancy identifies a case to investigate. Two implementations can agree on the same mistake, so the suite also checks fixed expected values.
