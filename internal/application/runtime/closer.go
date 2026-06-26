package runtime

import (
	"context"
	"log/slog"
	"sync"
)

type CloseFunc func() error

type Closer struct {
	mu    sync.Mutex
	funcs []namedFunc
}

type namedFunc struct {
	name string
	fn   CloseFunc
}

func NewCloser() *Closer {
	return &Closer{}
}

func (c *Closer) Add(name string, fn CloseFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, namedFunc{name: name, fn: fn})
}

func (c *Closer) Close() error {
	c.CloseContext(context.Background())
	return nil
}

func (c *Closer) CloseContext(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	funcs := make([]namedFunc, len(c.funcs))
	copy(funcs, c.funcs)
	c.mu.Unlock()

	for i := len(funcs) - 1; i >= 0; i-- {
		nf := funcs[i]
		if err := nf.fn(); err != nil {
			slog.Warn("close failed", "component", nf.name, "error", err)
		}
	}
}
