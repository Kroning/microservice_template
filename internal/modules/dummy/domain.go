package dummy

// Dummy is a main entity in the domain.
type Dummy struct {
	ID   int64
	Name string
}

// ListRequest represents incoming params for Service.List method.
type ListRequest struct {
	FilterID   int64
	FilterName string
	Offset     int32
	Limit      int32
}
