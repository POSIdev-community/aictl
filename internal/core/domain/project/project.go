package project

import "github.com/google/uuid"

// Type is the user-facing project kind (never API *Based enum names).
type Type string

const (
	TypeSource Type = "source"
	TypeSbom   Type = "sbom"
)

func (t Type) IsSbom() bool {
	return t == TypeSbom
}

func (t Type) IsSource() bool {
	return t == TypeSource || t == ""
}

type Project struct {
	Id   uuid.UUID
	Name string
	Type Type
}

func NewProject(id uuid.UUID, name string) Project {
	return Project{
		Id:   id,
		Name: name,
		Type: TypeSource,
	}
}

func NewProjectWithType(id uuid.UUID, name string, typ Type) Project {
	if typ == "" {
		typ = TypeSource
	}

	return Project{
		Id:   id,
		Name: name,
		Type: typ,
	}
}
