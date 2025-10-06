package project

import "github.com/google/uuid"

type Project struct {
	Id   uuid.UUID
	Name string
}

func NewProject(id uuid.UUID, name string) Project {
	return Project{
		Id:   id,
		Name: name,
	}
}
