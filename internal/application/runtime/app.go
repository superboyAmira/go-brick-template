package runtime

import (
	"context"
	"fmt"
	"log/slog"
)

type Info struct {
	Name    string
	Version string
}

type Module interface {
	Init(ctx context.Context, info Info) error
	Run(ctx context.Context, c *Closer)
}

type Application struct {
	name    string
	version string
	modules []Module
}

func New(name string) *Application {
	return &Application{name: name, version: "0.1.0"}
}

func (a *Application) WithModule(m Module) *Application {
	a.modules = append(a.modules, m)
	return a
}

func (a *Application) Build(ctx context.Context) error {
	info := Info{Name: a.name, Version: a.version}
	for i, m := range a.modules {
		if err := m.Init(ctx, info); err != nil {
			return fmt.Errorf("init module %d: %w", i, err)
		}
	}
	return nil
}

func (a *Application) Run(ctx context.Context, c *Closer) {
	for _, m := range a.modules {
		m.Run(ctx, c)
	}
	slog.Info("application running", "name", a.name, "modules", len(a.modules))
	<-ctx.Done()
	slog.Info("shutting down", "name", a.name)
}
