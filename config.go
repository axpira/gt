package gt

import (
	"encoding/json"
	"errors"
	"fmt"
)

func newConfig(opts ...option) config {
	return config{
		diff:     DefaultDiffFunc,
		fail:     DefaultFailFunc,
		failHook: DefaultFailHookFunc,
	}.With(opts...)
}

type diffFunc func(want, got any) string

type config struct {
	diff     diffFunc
	fail     func(testingT)
	failHook func()
}

func With(opts ...option) config {
	return newConfig(opts...)
}

func (c config) With(opts ...option) config {
	cfg := c
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func (c config) ErrorIs(t testingT, prefix string, err, target error, opts ...option) bool {
	if !errors.Is(err, target) {
		t.Helper()
		config := c.With(opts...)
		if target == nil {
			t.Log(fmt.Sprintf("%s: want no error got %q", prefix, err))
		} else {
			t.Log(fmt.Sprintf("%s: want error %q got %q", prefix, target, err))
		}
		config.failHook()
		config.fail(t)
		return false
	}
	return true
}

func (c config) Equal(t testingT, prefix string, want, got any, opts ...option) bool {
	config := c.With(opts...)
	if diff := config.diff(want, got); diff != "" {
		t.Helper()
		t.Log(fmt.Sprintf("%s\n%s", prefix, diff))
		config.failHook()
		config.fail(t)
		return false
	}
	return true
}

func (c config) NotEqual(t testingT, prefix string, want, got any, opts ...option) bool {
	config := c.With(opts...)
	if diff := config.diff(want, got); diff == "" {
		t.Helper()
		t.Log(fmt.Sprintf("%s is equals", prefix))
		config.failHook()
		config.fail(t)
		return false
	}
	return true
}

type option func(*config)

func WithFailLazy() option {
	return func(c *config) {
		c.fail = func(t testingT) {
			t.Fail()
		}
	}
}

func WithFailHook(fn func()) option {
	return func(c *config) {
		c.failHook = fn
	}
}

func WithDiffFunc(diffFn diffFunc) option {
	return func(c *config) {
		c.diff = diffFn
	}
}

func WithMarshal(marshal func(input any) (any, error)) option {
	return func(cfg *config) {
		cfg.diff = func(want, got any) string {
			w, err := marshal(want)
			if err != nil {
				return err.Error()
			}
			g, err := marshal(got)
			if err != nil {
				return err.Error()
			}
			return DefaultDiffFunc(w, g)
		}
	}
}

func WithJSONDiff(cfg *config) {
	WithMarshal(func(value any) (m any, err error) {
		m = make(map[string]interface{}, 0)
		defer func() {
			if r := recover(); r != nil {
				m = nil
				err = fmt.Errorf("error on parse %v want []byte and got %T", value, value)
			}
		}()
		err = json.Unmarshal(value.([]byte), &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	})(cfg)
}
