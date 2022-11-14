package modules

import (
	"context"
	"io"
)

// Module -
type Module interface {
	io.Closer

	Start(ctx context.Context)
}
