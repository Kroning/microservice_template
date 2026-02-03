package repository

import (
	"time"

	"{{index .App "git"}}/internal/modules/dummy"
)

// dummyModel is the database representation of the Dummy entity.
type dummyModel struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// toDomain converts the database model to the domain entity.
func (m *dummyModel) toDomain() dummy.Dummy {
	return dummy.Dummy{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
