package branch

import "github.com/google/uuid"

type Branch struct {
	Id          uuid.UUID
	Name        string
	Description string
	IsWorking   bool
}

func NewBranch(id uuid.UUID, name string, description string, isWorking bool) Branch {
	return Branch{
		Id:          id,
		Name:        name,
		Description: description,
		IsWorking:   isWorking,
	}
}
