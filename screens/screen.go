package screens

import (
	"context"
)

type Screen interface {
	Show(ctx context.Context)
}
