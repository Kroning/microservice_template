package dummy

import "context"

//go:generate go run go.uber.org/mock/mockgen -source=interfaces.go -destination=dummymocks/dummymocks.go -package=dummymocks

type Repository interface {
	Create(ctx context.Context, name string) (*Dummy, error)
	List(ctx context.Context, request ListRequest) ([]Dummy, error)
}

type Service interface {
	Create(ctx context.Context, name string) (*Dummy, error)
	List(ctx context.Context, request ListRequest) ([]Dummy, error)
}
