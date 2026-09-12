package scafeeds

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeScaFeeds Type = "sca_feeds"
)

type Status string

const (
	StatusActive     Status = "active"
	StatusCurrent    Status = "current"
	StatusArchived   Status = "archived"
	StatusRolledBack Status = "rolled_back"
)

func (s Status) Valid() bool {
	switch s {
	case StatusActive, StatusCurrent, StatusArchived, StatusRolledBack:
		return true
	default:
		return false
	}
}

type Trigger string

const (
	TriggerManual    Trigger = "manual"
	TriggerScheduled Trigger = "scheduled"
)

type PackageUser struct {
	Id        uuid.UUID
	UserName  string
	Email     string
	TokenName string
}

type Package struct {
	Id             int32
	PackageType    Type
	Version        string
	Status         Status
	FileName       string
	FileSize       int64
	FileHash       string
	TriggeredBy    Trigger
	UploadedAt     time.Time
	UploadedBy     *PackageUser
	LastModifiedAt time.Time
	LastModifiedBy *PackageUser
}
