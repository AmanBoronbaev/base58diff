# base58diff

Compare two Go implementations of Bitcoin's Base58 alphabet on the same input. Use it to replay a test vector, check a codec upgrade, or turn an encoding disagreement into a small Go regression test.

```sh
go run . encode 000039  # 11z
go run . decode 11z     # 000039
go run . decode 1O      # rejected
```

Each command runs both [btcutil](https://github.com/btcsuite/btcd/tree/master/btcutil) and [mr-tron](https://github.com/mr-tron/base58). Encoding checks the resulting text and decodes it back to the original bytes. Decoding compares acceptance, decoded bytes and canonical re-encoding. Leading zero bytes matter: removing one leading `1` changes the payload.

The program exits with 0 when both libraries agree, including when both reject a string. Exit 1 reports a disagreement or failed round trip on stderr. Exit 2 indicates a command or input error. Empty bytes encode to empty text. Empty text is excluded from decode comparisons because the libraries expose different empty-input behavior. Inputs use raw Base58 without a version byte or Base58Check checksum.

## Install

With Go 1.22 or later:

```sh
go install github.com/AmanBoronbaev/base58diff@v0.1.0
base58diff decode 11z
```

The command is installed in `GOBIN`, or in `$(go env GOPATH)/bin` by default.

## Replay a vector file

```sh
go run . check testdata/bitcoin-core-base58.json
# 21 vectors passed
printf '%s\n' '[["000039","11z"]]' | go run . check -
```

Supply a JSON array of hex/Base58 string pairs. The checker streams one pair at a time and reports the first failing vector by its one-based index. A compiled binary exits with 1 for a mismatch and 2 for malformed or empty files.

## Run the corpus

```sh
go test ./...
go test -fuzz=FuzzEncode -fuzztime=30s -parallel=2
```

The regular suite includes all 21 [Bitcoin Core v29.0 Base58 vectors](https://github.com/bitcoin/bitcoin/blob/v29.0/src/test/data/base58_encode_decode.json), stored under the [upstream MIT license](testdata/COPYING.bitcoin-core). It also covers fixed input/output pairs, a deterministic payload corpus with varying leading-zero prefixes, and invalid characters at decoder chunk boundaries. The fuzz target varies binary payloads up to 256 bytes and checks both directions. Go saves a minimized failure in `testdata/fuzz/FuzzEncode/` and prints the exact replay command. Commit a confirmed failure there to include it in subsequent ordinary test runs.

To compare a dependency upgrade, change its version in `go.mod`, run `go mod tidy` and rerun the corpus. A discrepancy identifies a case to investigate. Two implementations can agree on the same mistake, so the suite also checks fixed expected values.
