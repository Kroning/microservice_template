package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"{{index .App "git"}}/internal/modules/dummy"
	"{{index .App "git"}}/pkg/storage"
)

// PostgresRepository implements dummy.Repository interface using PostgreSQL.
type PostgresRepository struct {
	db storage.AbstractDB
}

// NewPostgresRepository creates a new PostgreSQL repository instance.
func NewPostgresRepository(db storage.AbstractDB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create inserts a new dummy record into the database.
func (r *PostgresRepository) Create(ctx context.Context, name string) (*dummy.Dummy, error) {
	now := time.Now()

	var id int64
	err := r.db.GetContext(ctx, &id, `
		INSERT INTO dummy (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`, name, now, now)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create dummy")
	}

	return &dummy.Dummy{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// List retrieves dummy records from the database based on the request filters.
func (r *PostgresRepository) List(ctx context.Context, request dummy.ListRequest) ([]dummy.Dummy, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM dummy
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	if request.FilterID > 0 {
		query += " AND id = $" + formatArgPos(argPos)
		args = append(args, request.FilterID)
		argPos++
	}

	if request.FilterName != "" {
		query += " AND name LIKE $" + formatArgPos(argPos)
		args = append(args, "%"+request.FilterName+"%")
		argPos++
	}

	query += " ORDER BY id DESC"

	// Set default limit if not specified
	limit := request.Limit
	if limit <= 0 {
		limit = 10
	}
	query += " LIMIT $" + formatArgPos(argPos)
	args = append(args, limit)
	argPos++

	if request.Offset > 0 {
		query += " OFFSET $" + formatArgPos(argPos)
		args = append(args, request.Offset)
	}

	var dummyModels []dummyModel
	err := r.db.SelectContext(ctx, &dummyModels, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list dummies")
	}

	dummies := make([]dummy.Dummy, len(dummyModels))
	for i, e := range dummyModels {
		dummies[i] = e.toDomain()
	}

	return dummies, nil
}

// formatArgPos formats argument position for PostgreSQL query.
func formatArgPos(pos int) string {
	return strconv.Itoa(pos)
}
