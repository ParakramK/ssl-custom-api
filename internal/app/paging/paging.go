package paging

import (
	"github.com/google/uuid"
)

// Sort is one ORDER BY term. Field is a logical name each repository
// maps to a real column; Desc selects descending order.
type Sort struct {
	Field string
	Desc  bool
}

// Cursor is a decoded keyset position. Repositories compare the leading
// sort column against it, using ID as the tiebreak.
type Cursor struct {
	ID  uuid.UUID
	Key string
}

// Query is the neutral pagination/filter request handlers build from
// fiber's paginate.PageInfo and pass down to services and repositories.
// It keeps fiber types out of the domain layers.
type Query struct {
	Limit  int
	Sorts  []Sort
	Cursor *Cursor // nil = first page
}

// LeadingSort returns the primary sort, defaulting to id ascending
// (mirrors the paginate middleware defaults).
func (q Query) LeadingSort() Sort {
	if len(q.Sorts) > 0 {
		return q.Sorts[0]
	}
	return Sort{Field: "id"}
}
