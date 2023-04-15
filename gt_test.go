package gt

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type mockTestingT struct {
	called  map[string]int
	logFunc func(...any)
}

func (m *mockTestingT) Helper() {
}

func (m *mockTestingT) FailNow() {
	m.called["FailNow"]++
}

func (m *mockTestingT) Fail() {
	m.called["Fail"]++
}

func (m *mockTestingT) Log(args ...any) {
	m.called["Log"]++
	m.logFunc(args...)
}

func TestEqual(t *testing.T) {
	errUnknown := errors.New("unknown")
	tests := map[string]struct {
		wantLog    []any
		wantCalled map[string]int
		wantResult bool
		execute    func(testingT) bool
	}{
		"not equal": {
			wantCalled: map[string]int{},
			wantLog:    []any{""},
			wantResult: true,
			execute: func(t testingT) bool {
				return NotEqual(t, "prefix", "a", "b")
			},
		},
		"not equal with equal values": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog:    []any{`prefix is equals`},
			wantResult: false,
			execute: func(t testingT) bool {
				return NotEqual(t, "prefix", "a", "a")
			},
		},
		"equal": {
			wantCalled: map[string]int{},
			wantLog:    []any{""},
			wantResult: true,
			execute: func(t testingT) bool {
				return Equal(t, "prefix", "a", "a")
			},
		},
		"equal with custom diff func": {
			wantCalled: map[string]int{},
			wantLog:    []any{""},
			wantResult: true,
			execute: func(t testingT) bool {
				return Equal(t, "prefix2", "a", "b", WithDiffFunc(func(want, got any) string {
					return ""
				}))
			},
		},
		"equal: not equal struct": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog: []any{`struct
want: struct { name string; age int }{name:"no name", age:39}
got:  struct { name string; age int }{name:"anyone", age:39}`},
			wantResult: false,
			execute: func(t testingT) bool {
				return Equal(t,
					"struct",
					struct {
						name string
						age  int
					}{
						name: "no name",
						age:  39,
					},
					struct {
						name string
						age  int
					}{
						name: "anyone",
						age:  39,
					},
				)
			},
		},
		"equal: not equal": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog: []any{`not equals
want: "a"
got:  "b"`},
			wantResult: false,
			execute: func(t testingT) bool {
				return Equal(t, "not equals", "a", "b")
			},
		},
		"iserror fail": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog:    []any{"iserror fail: want error \"unknown\" got \"test\""},
			wantResult: false,
			execute: func(t testingT) bool {
				return ErrorIs(t, "iserror fail", fmt.Errorf("test"), errUnknown)
			},
		},
		"iserror success": {
			wantCalled: map[string]int{},
			wantLog:    []any{},
			wantResult: true,
			execute: func(t testingT) bool {
				return ErrorIs(t, "iserror fail", fmt.Errorf("test: %w", errUnknown), errUnknown)
			},
		},
		"iserror nil": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog:    []any{`iserror fail: want no error got "test: unknown"`},
			wantResult: false,
			execute: func(t testingT) bool {
				return ErrorIs(t, "iserror fail", fmt.Errorf("test: %w", errUnknown), nil)
			},
		},
		"with json": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog: []any{`with json
want: map[string]interface {}{"id":1, "name":"test"}
got:  map[string]interface {}{"id":2, "name":"test"}`},
			wantResult: false,
			execute: func(t testingT) bool {
				return Equal(t, "with json", []byte(`{"id":1,"name": "test"}`), []byte(`{"id":2,"name": "test"}`), WithJSONDiff)
			},
		},
		"with json invalid type": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog: []any{`with json invalid type
error on parse {"id":1,"name": "test"} want []byte and got string`},
			wantResult: false,
			execute: func(t testingT) bool {
				return Equal(t, "with json invalid type", `{"id":1,"name": "test"}`, []byte(`{"id":2,"name": "test"}`), WithJSONDiff)
			},
		},
		"with json invalid type 2": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog: []any{`with json invalid type 2
error on parse {"id":2,"name": "test"} want []byte and got string`},
			wantResult: false,
			execute: func(t testingT) bool {
				return Equal(t, "with json invalid type 2", []byte(`{"id":1,"name": "test"}`), `{"id":2,"name": "test"}`, WithJSONDiff)
			},
		},
		"with json invalid json": {
			wantCalled: map[string]int{"FailNow": 1, "Log": 1},
			wantLog: []any{`with json invalid json
invalid character '}' looking for beginning of object key string`},
			wantResult: false,
			execute: func(t testingT) bool {
				return Equal(t, "with json invalid json", []byte(`{"id":1,"name": "test",}`), []byte(`{"id":2,"name": "test"}`), WithJSONDiff)
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			mockedTestingT := &mockTestingT{
				called: make(map[string]int, 5),
				logFunc: func(args ...any) {
					if !reflect.DeepEqual(tc.wantLog, args) {
						t.Errorf("want log '%v' got '%v'", tc.wantLog, args)
					}
				},
			}

			result := tc.execute(mockedTestingT)
			if result != tc.wantResult {
				t.Errorf("\n++want %v\n--got %v", tc.wantResult, result)
			}
			if !reflect.DeepEqual(tc.wantCalled, mockedTestingT.called) {
				t.Errorf("want %v got %v", tc.wantCalled, mockedTestingT.called)
			}
		})
	}
}
