package spy

import (
	"strings"
	"sync"
	"testing"
)

// helperFunction auxilia no teste de captura do nome da função chamadora.
func helperFunction(r *Recorder, args ...any) {
	r.Record(args...)
}

func TestRecorder_Record(t *testing.T) {
	t.Run("captures function name and arguments", func(t *testing.T) {
		r := &Recorder{}
		expectedArg := "test-value"

		// Chama através da helper para verificar se o nome "helperFunction" é capturado
		helperFunction(r, expectedArg)

		calls := r.History
		if len(calls) != 1 {
			t.Fatalf("got %d calls, want 1", len(calls))
		}

		call := calls[0]

		// Verifica o nome da função
		// O nome completo depende do caminho do pacote, então verificamos o sufixo.
		wantNameSuffix := "helperFunction"
		// call[0] é o nome da função
		gotName, ok := call[0].(string)
		if !ok {
			t.Fatalf("first element of call history should be string, got %T", call[0])
		}

		if !strings.HasSuffix(gotName, wantNameSuffix) {
			t.Errorf("got function name %q, want suffix %q", gotName, wantNameSuffix)
		}

		// Verifica os argumentos
		// call[1] é o primeiro argumento passado para Record
		if len(call) != 2 {
			t.Fatalf("got %d elements in call record, want 2 (name + 1 arg)", len(call))
		}
		if call[1] != expectedArg {
			t.Errorf("got arg %v, want %v", call[1], expectedArg)
		}
	})
}

func TestRecorder_Reset(t *testing.T) {
	r := &Recorder{}
	r.Record("first")
	r.Record("second")

	if got := len(r.History); got != 2 {
		t.Fatalf("pre-reset: got %d calls, want 2", got)
	}

	r.Reset()

	if got := len(r.History); got != 0 {
		t.Errorf("post-reset: got %d calls, want 0", got)
	}
}

func TestRecorder_Concurrency(t *testing.T) {
	r := &Recorder{}
	var wg sync.WaitGroup
	workers := 20
	iterations := 100

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				r.Record(id, j)
			}
		}(i)
	}

	wg.Wait()

	totalExpected := workers * iterations
	if got := len(r.History); got != totalExpected {
		t.Errorf("got %d calls, want %d", got, totalExpected)
	}
}
