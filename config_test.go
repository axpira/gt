package gt

import (
	"reflect"
	"testing"
)

type mockDiffFunc struct {
	called   int
	diffFunc func(want, got any) string
}

func (m *mockDiffFunc) diff(want, got any) string {
	m.called++
	return m.diffFunc(want, got)
}

func TestNewConfig(t *testing.T) {
	type diffFuncWant struct {
		want, got any
		result    string
	}
	tests := map[string]struct {
		diffFuncWant   diffFuncWant
		diffFuncCalled int
		wantCalled     map[string]int
		failHookCalled int
		execute        func(t testingT, c config)
	}{
		"t.Equal with not equal values and faillazy": {
			diffFuncWant:   diffFuncWant{"a", "b", "not equals"},
			diffFuncCalled: 1,
			wantCalled: map[string]int{
				"Fail": 1,
				"Log":  1,
			},
			failHookCalled: 1,
			execute: func(t testingT, c config) {
				c.Equal(t, "prefix", "a", "b", WithFailLazy())
			},
		},
		"t.Equal with not equal values": {
			diffFuncWant:   diffFuncWant{"a", "b", "not equals"},
			diffFuncCalled: 1,
			wantCalled: map[string]int{
				"FailNow": 1,
				"Log":     1,
			},
			failHookCalled: 1,
			execute: func(t testingT, c config) {
				c.Equal(t, "prefix", "a", "b")
			},
		},
		"t.NotEqual with equal values": {
			diffFuncWant:   diffFuncWant{"a", "b", ""},
			diffFuncCalled: 1,
			wantCalled: map[string]int{
				"FailNow": 1,
				"Log":     1,
			},
			failHookCalled: 1,
			execute: func(t testingT, c config) {
				c.NotEqual(t, "prefix", "a", "b")
			},
		},
		"t.Equal with equal values": {
			diffFuncWant:   diffFuncWant{"a", "b", ""},
			diffFuncCalled: 1,
			wantCalled:     map[string]int{},
			execute: func(t testingT, c config) {
				c.Equal(t, "prefix", "a", "b")
			},
		},
		"t.NotEqual with not equal values": {
			diffFuncWant:   diffFuncWant{"a", "b", "not equal"},
			diffFuncCalled: 1,
			wantCalled:     map[string]int{},
			execute: func(t testingT, c config) {
				c.NotEqual(t, "prefix", "a", "b")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockedTestingT := &mockTestingT{
				called:  make(map[string]int),
				logFunc: func(args ...any) {},
			}
			mockedDiffFunc := &mockDiffFunc{
				diffFunc: func(want any, got any) string {
					if want != tc.diffFuncWant.want {
						t.Errorf("diff func want param want %v got %v", tc.diffFuncWant.want, want)
					}
					if got != tc.diffFuncWant.got {
						t.Errorf("diff func got param want %v got %v", tc.diffFuncWant.got, got)
					}
					return tc.diffFuncWant.result
				},
			}

			failHookCalled := 0
			c := newConfig(WithDiffFunc(mockedDiffFunc.diff), WithFailHook(
				func() {
					failHookCalled++
				},
			))
			tc.execute(mockedTestingT, c)
			if mockedDiffFunc.called != tc.diffFuncCalled {
				t.Errorf("mocked diff func must be called %v and got %v", tc.diffFuncCalled, mockedDiffFunc.called)
			}
			if !reflect.DeepEqual(mockedTestingT.called, tc.wantCalled) {
				t.Errorf("called: want %#+v got %#+v", tc.wantCalled, mockedTestingT.called)
			}
			if failHookCalled != tc.failHookCalled {
				t.Errorf("failHookCalled: want %#+v got %#+v", tc.failHookCalled, failHookCalled)
			}
		})
	}
}

func TestWith(t *testing.T) {
	called := map[string]int{}
	c := With(
		WithDiffFunc(func(want, got any) string {
			called["diffFunc"]++
			return ""
		}),
	)

	c.Equal(t, "", nil, nil)
	want := map[string]int{
		"diffFunc": 1,
	}
	if !reflect.DeepEqual(want, called) {
		t.Errorf("diffFunc not called")
	}
}
