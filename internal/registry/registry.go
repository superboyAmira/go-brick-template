package registry

import (
	"github.com/go-brick-template/go-brick-template/internal/brick/item"
)

// Registry holds wired brick contracts for HTTP handlers and cross-brick DI.
type Registry struct {
	Item item.Service
}
