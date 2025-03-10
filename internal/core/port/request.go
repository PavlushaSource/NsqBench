package port

import (
	"context"
)

type RequestService interface {
	Run(ctx context.Context, iterations int) error
	Close() error
}
