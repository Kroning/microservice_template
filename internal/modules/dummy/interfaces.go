package dummy

import "context"

//go:generate mockgen -source=./domain.go -destination=usermocks/usermocks.go -package=usermocks

type Service interface {
	Create(ctx context.Context, name string) (*Dummy, error)
	List(ctx context.Context, request ListRequest) ([]Dummy, error)
}
