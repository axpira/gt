# gt (Go Test)

[![Go Report Card](https://goreportcard.com/badge/github.com/axpira/gt)](https://goreportcard.com/report/github.com/axpira/gt)
[![GoDoc](https://pkg.go.dev/badge/github.com/axpira/gt.svg)](https://pkg.go.dev/github.com/axpira/gt)
[![Coverage Status](https://coveralls.io/repos/github/axpira/gt/badge.svg?branch=main)](https://coveralls.io/github/axpira/gt?branch=main)

**gt** is a minimalist, idiomatic testing library for Go. It provides powerful assertions and a lightweight spying mechanism without the bloat of heavy frameworks.

Designed for **Effective Go**:
- **Zero Dependencies**: Pure Go standard library.
- **No Magic**: Explicit, readable, and type-safe.
- **Extensible**: Plug in your favorite diff tools (like `go-cmp`) or JSON comparators.

## Installation

```sh
go get github.com/axpira/gt
```

## Quick Start

`gt` works alongside the standard `testing` package.

```go
package main_test

import (
	"errors"
	"testing"

	"github.com/axpira/gt"
)

func TestExample(t *testing.T) {
	got := 2 + 2
	gt.Equal(t, "math check", 4, got)

	err := errors.New("something went wrong")
	target := errors.New("something went wrong")
	
	// Checks if err matches target (using errors.Is)
	gt.ErrorIs(t, "error check", err, target)
}
```

## Advanced Usage

### Custom Diff Function
By default, `gt` uses `reflect.DeepEqual`. You can easily plug in `google/go-cmp` for better diffs:

```go
import "github.com/google/go-cmp/cmp"

// ... inside test
gt.Equal(t, "complex struct", want, got, gt.WithDiffFunc(func(want, got any) string {
    if diff := cmp.Diff(want, got); diff != "" {
        return diff
    }
    return ""
}))
```

### JSON Comparison
Compare JSON strings or bytes ignoring whitespace:

```go
gt.Equal(t, "json check", 
    []byte(`{"id": 1, "name": "foo"}`), 
    []byte(`{
        "name": "foo",
        "id": 1
    }`), 
    gt.WithJSONDiff,
)
```

### Custom Failure Behavior
Change how `gt` handles failures (default is `FailNow`):

```go
// Don't stop execution on failure
gt.Equal(t, "soft check", 1, 2, gt.WithFailLazy())

// Custom hook (e.g., for logging)
gt.Equal(t, "hook check", 1, 2, gt.WithFailHook(func() {
    t.Log("Assertion failed!")
}))
```

## Mocking & Spying

`gt` includes a `spy` package to facilitate **Manual Mocking**. Instead of generating complex mock code, we encourage defining simple mock structs that record their interactions.

### The `spy.Recorder`

The `spy.Recorder` is a thread-safe helper that records function calls and arguments.

### Step-by-Step Example

**1. Define the Interface**
```go
type Repository interface {
    Save(ctx context.Context, data string) error
}
```

**2. Create the Mock**
Embed `spy.Recorder` in your struct.

```go
import "github.com/axpira/gt/spy"

type MockRepo struct {
    spy.Recorder
    
    // Configurable return values
    SaveErr error
}

func (m *MockRepo) Save(ctx context.Context, data string) error {
    // Record the call. 
    // It automatically captures the method name "Save" and arguments.
    m.Record(ctx, data) 
    return m.SaveErr
}
```

**3. Use in Test**

```go
func TestService_Save(t *testing.T) {
    // Setup Mock
    mock := &MockRepo{SaveErr: nil}
    
    // Execute logic
    err := mock.Save(context.Background(), "important data")
    
    // Assertions
    gt.ErrorIs(t, "save error", err, nil)
    
    // Verify Interactions
    // History is [][]any: [[MethodName, Arg1, Arg2...], [MethodName, ...]]
    
    gt.Equal(t, "call count", 1, len(mock.History))
    
    lastCall := mock.History[0]
    gt.Equal(t, "method name", "Save", lastCall[0])
    gt.Equal(t, "argument", "important data", lastCall[2]) // [0]=Name, [1]=ctx, [2]=data
}
```

### Why Manual Mocks?
- **Type Safety**: Compiler checks your mocks.
- **Refactoring**: Rename methods and your mocks update automatically (mostly).
- **Simplicity**: No code generation steps or complex DSLs.

## License

MIT License. See LICENSE file for details.
