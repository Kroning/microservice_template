package postgresql

import (
	"context"
	"time"
)

type metricsDB struct {
	enable bool
}

func (w *metricsDB) write(ctx context.Context, started time.Time, query string, err error) {
	// TODO: implement me with DB metrics
	if !w.enable {
		return
	}
}
