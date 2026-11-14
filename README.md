# go.os
go network utility

## Lock33 ID package

The `lock33id` package provides the ISUD identity witness, derived modulus,
shift, and seed helpers, plus the Lock33 pipe encoder and decoder.

The encoder writes each Unicode code point, shifted by `(DeriveSeed(x) % 33) +
1`, as hexadecimal and separates fields with `|`.

Run the package tests with:

```sh
go test ./...
```
