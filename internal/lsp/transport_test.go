package lsp

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

type nopWriteCloser struct {
	*bytes.Buffer
}

func (n nopWriteCloser) Close() error { return nil }

func TestCallReturnsWhenContextIsCanceled(t *testing.T) {
	client := &Client{
		stdin:    nopWriteCloser{Buffer: &bytes.Buffer{}},
		handlers: make(map[string]chan *Message),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := client.Call(ctx, "test/request", map[string]string{"key": "value"}, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Call() error = %v, want %v", err, context.DeadlineExceeded)
	}

	client.handlersMu.RLock()
	defer client.handlersMu.RUnlock()
	if len(client.handlers) != 0 {
		t.Fatalf("Call() left %d response handlers registered after cancellation", len(client.handlers))
	}
}
