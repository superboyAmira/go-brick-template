package runtime

import (
	"context"
	"testing"
)

type stubModule struct {
	inited bool
}

func (s *stubModule) Init(ctx context.Context, _ Info) error {
	s.inited = true
	return nil
}

func (s *stubModule) Run(ctx context.Context, _ *Closer) {}

func TestApplicationBuild(t *testing.T) {
	ctx := context.Background()
	mod := &stubModule{}
	a := New("test").WithModule(mod)
	if err := a.Build(ctx); err != nil {
		t.Fatal(err)
	}
	if !mod.inited {
		t.Fatal("expected module init")
	}
}
