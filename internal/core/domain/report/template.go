package report

import "github.com/google/uuid"

// Template is a report template listing row (id + name).
type Template struct {
	Id   uuid.UUID
	Name string
}
