package dummy

import "time"

// Dummy is a main entity in the domain.
type Dummy struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ListRequest represents incoming params for Service.List method.
type ListRequest struct {
	FilterID   int64
	FilterName string
	Offset     int32
	Limit      int32
}
