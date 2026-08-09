package knowledge

import "context"

type Repository interface {
	Store(ctx context.Context, source string, content []byte) error
	Search(ctx context.Context, query string) ([]string, error)
}
