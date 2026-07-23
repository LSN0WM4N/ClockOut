package dispatcher_test

import (
	"ClockOut/internal/dispatcher"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testEvent struct {
	Type    string
	Payload string
}

func (e testEvent) GetType() string {
	return e.Type
}

func (e testEvent) GetPayload() string {
	return e.Payload
}

func TestRegisterHandler(t *testing.T) {
	ch := make(chan testEvent)

	d := dispatcher.NewDispatcher(ch)

	called := false

	d.RegisterHandler("test", func(testEvent) {
		called = true
	})

	handler, ok := d.GetHandler("test")
	if !ok {
		t.Fatal("handler not registered")
	}

	handler(testEvent{})

	if !called {
		t.Fatal("handler was not executed")
	}
}

func TestDispatcherDispatchesEvent(t *testing.T) {
	ch := make(chan testEvent)

	d := dispatcher.NewDispatcher(ch)

	var wg sync.WaitGroup
	wg.Add(1)

	d.RegisterHandler("test", func(e testEvent) {
		if e.Payload != "123" {
			t.Errorf("unexpected ID %s", e.Payload)
		}
		wg.Done()
	})

	d.Start()

	ch <- testEvent{
		Type:    "test",
		Payload: "123",
	}

	wait(&wg, t)
}

func TestDispatcherIgnoresUnknownEvent(t *testing.T) {
	ch := make(chan testEvent)

	d := dispatcher.NewDispatcher(ch)

	d.Start()

	ch <- testEvent{
		Type: "unknown",
	}

	time.Sleep(20 * time.Millisecond)
}

func TestDispatcherRecoversFromPanic(t *testing.T) {
	ch := make(chan testEvent)

	d := dispatcher.NewDispatcher(ch)

	var wg sync.WaitGroup
	wg.Add(1)

	d.RegisterHandler("panic", func(testEvent) {
		defer wg.Done()
		panic("boom")
	})

	d.Start()

	ch <- testEvent{
		Type: "panic",
	}

	wait(&wg, t)
}

func TestDispatcherMultipleEvents(t *testing.T) {
	ch := make(chan testEvent)

	d := dispatcher.NewDispatcher(ch)
	const total = 100

	var count atomic.Int32
	var wg sync.WaitGroup
	wg.Add(total)

	d.RegisterHandler("test", func(testEvent) {
		count.Add(1)
		wg.Done()
	})

	d.Start()

	for i := 0; i < total; i++ {
		ch <- testEvent{
			Type: "test",
		}
	}

	wait(&wg, t)

	if count.Load() != total {
		t.Fatalf("expected %d events got %d", total, count.Load())
	}
}

func wait(wg *sync.WaitGroup, t *testing.T) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for handler")
	}
}
