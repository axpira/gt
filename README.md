# Package gt (Go Test)

[![Go Report Card](https://goreportcard.com/badge/github.com/axpira/gt)](https://goreportcard.com/report/github.com/axpira/gt)
[![GoDoc](https://pkg.go.dev/badge/github.com/axpira/gt.svg)](https://pkg.go.dev/github.com/axpira/gt)
[![Coverage Status](https://coveralls.io/repos/github/axpira/gt/badge.svg)](https://coveralls.io/github/axpira/gt)

*gt* a basic but powerful test library used along side native go testing.

Highly customizable, simple and Zero Dependency

## Installation

Use go get.

```sh
go get github.com/axpira/gt
```

Then import the assert package into your own code.

```sh
import "github.com/axpira/gt"
```

## Usage and documentation

Please see [Godoc](https://pkg.go.dev/github.com/axpira/gt) for detailed usage
docs.

### Example

```go
package whatever

import (
    "errors"
    "testing"
    "github.com/axpira/gt"
)

func TestEqual(t *testing.T) {

    gotRes, gotErr := awesomefunc()
    gt.ErrorIs(t, "awesome response err", gotErr, targetErr)
    gt.Equal(t, "awesome response", wantRes, gotRes)
    // by default when not match will FailNow and stop execution of next steps

    // but you can change this behavior
    gt.ErrorIs(t, "awesome response err", gotErr, targetErr, gt.WithFailLazy())

    // or you can receive a hook when fail
    gt.ErrorIs(t, "awesome response err", gotErr, targetErr, gt.WithFailHook(func(){
        t.Errof("error is %s not %s", wantErr, gotErr)
    }))

    // the Equal func the diff is using reflect.DeepEqual but you can use whatever you prefer
    // here is changing to the great cmp, need to import "github.com/google/go-cmp/cmp"
    diffFunc = func(want, got any) string {
        if diff := cmp.Diff(want, got); diff != "" {
            return fmt.Sprintf("mismatch (-want +got):\n%s", diff)
        }
        return ""
    }
    gt.Equal(t, "awesome response", wantRes, gotRes, gt.WithDiffFunc(diffFunc))

    // and if you want to compare two json, you can use like this
    gt.Equal(t, "awesome json diff", []byte(`{"a":"b"}`), ]byte(`{"a":"b"}`), gt.WithJSONDiff)


    // and if you want to create gt with the config (then don't need to send always)
    gt := gt.With(gt.WithDiffFunc(diffFunc)) // or any other config
    // and use normal
    gt.Equal(t, "awesome response", wantRes, gotRes) // now will use your custom diff
```

## How to Contribute

Make a PR.

## License

Distributed under MIT License, please see license file in code for more details.
