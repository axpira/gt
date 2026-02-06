package spy

import (
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name      string
		want      any
		got       any
		wantDiff  string
		wantPanic bool
	}{
		{
			name: "exact match",
			want: [][]any{{"func", 1}},
			got:  [][]any{{"func", 1}},
		},
		{
			name:     "size mismatch: got > want",
			want:     [][]any{},
			got:      [][]any{{"func", 1}},
			wantDiff: "call #0:\n    WANT: (should not be called)\n    GOT:  func(1)",
		},
		{
			name:     "size mismatch: want > got",
			want:     [][]any{{"func", 1}},
			got:      [][]any{},
			wantDiff: "call #0:\n    WANT: func(1)\n    GOT:  (no call recorded)",
		},
		{
			name:     "value mismatch",
			want:     [][]any{{"func", 1}},
			got:      [][]any{{"func", 2}},
			wantDiff: "call #0:\n    WANT: func(1)\n    GOT:  func(2)",
		},
		{
			name:      "panic: nil want",
			want:      nil,
			got:       [][]any{},
			wantPanic: true,
		},
		{
			name:      "panic: nil got",
			want:      [][]any{},
			got:       nil,
			wantPanic: true,
		},
		{
			name:      "panic: invalid type want",
			want:      "invalid",
			got:       [][]any{},
			wantPanic: true,
		},
		{
			name:      "panic: invalid type got",
			want:      [][]any{},
			got:       "invalid",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.wantPanic && r == nil {
					t.Error("want panic, got none")
				}
				if !tt.wantPanic && r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()

			gotDiff := Diff(tt.want, tt.got)

			if tt.wantPanic {
				return
			}

			if tt.wantDiff == "" {
				if gotDiff != "" {
					t.Errorf("want empty diff, got:\n%s", gotDiff)
				}
				return
			}

			if !strings.Contains(gotDiff, tt.wantDiff) {
				t.Errorf("diff mismatch.\nwant substring:\n%s\ngot:\n%s", tt.wantDiff, gotDiff)
			}
		})
	}
}
