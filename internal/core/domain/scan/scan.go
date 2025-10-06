package scan

import "github.com/google/uuid"

type Scan struct {
	Id         uuid.UUID
	SettingsId uuid.UUID
}

func NewScan(id, settingsId uuid.UUID) *Scan {
	return &Scan{Id: id, SettingsId: settingsId}
}
