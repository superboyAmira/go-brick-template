package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const sep = "|"

// Encode builds an opaque cursor from sort key (updated_at) and row id.
func Encode(updatedAt time.Time, id uuid.UUID) string {
	raw := updatedAt.UTC().Format(time.RFC3339Nano) + sep + id.String()
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// Decode parses a cursor into updated_at and id.
func Decode(cursor string) (time.Time, uuid.UUID, error) {
	if cursor == "" {
		return time.Time{}, uuid.Nil, fmt.Errorf("empty cursor")
	}
	b, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor: %w", err)
	}
	parts := strings.SplitN(string(b), sep, 2)
	if len(parts) != 2 {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor format")
	}
	ts, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor time: %w", err)
	}
	id, err := uuid.Parse(parts[1])
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor id: %w", err)
	}
	return ts, id, nil
}

// ClampLimit normalizes list page size.
func ClampLimit(limit int) int {
	switch {
	case limit <= 0:
		return 20
	case limit > 100:
		return 100
	default:
		return limit
	}
}
