package queue

import "github.com/google/uuid"

// Entry is a scan queue or active-scan row for tabular CLI output.
type Entry struct {
	ProjectId uuid.UUID
	BranchId  uuid.UUID
	ScanId    uuid.UUID
	Stage     string
}
