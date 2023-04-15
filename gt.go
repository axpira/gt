package gt

import (
	"fmt"
	"reflect"
)

var (
	DefaultDiffFunc diffFunc = func(want, got any) string {
		if !reflect.DeepEqual(want, got) {
			return fmt.Sprintf("want: %#+v\ngot:  %#+v", want, got)
		}
		return ""
	}
	DefaultFailFunc = func(t testingT) {
		t.FailNow()
	}
	DefaultFailHookFunc = func() {}
)

type testingT interface {
	Fail()
	FailNow()
	Helper()
	Log(...any)
}

func ErrorIs(t testingT, prefix string, err, target error, opts ...option) bool {
	t.Helper()
	return newConfig(opts...).ErrorIs(t, prefix, err, target, opts...)
}

func Equal(t testingT, prefix string, want, got any, opts ...option) bool {
	t.Helper()
	return newConfig(opts...).Equal(t, prefix, want, got, opts...)
}

func NotEqual(t testingT, prefix string, want, got any, opts ...option) bool {
	t.Helper()
	return newConfig(opts...).NotEqual(t, prefix, want, got, opts...)
}
