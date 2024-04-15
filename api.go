package alor

import (
	"context"
	"time"
)

// какой api реализован
type IAlorClient interface {
	// текущее время сервера
	GetTime(ctx context.Context) (time.Time, error)

}
