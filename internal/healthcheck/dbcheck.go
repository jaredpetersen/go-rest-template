package healthcheck

import (
	"context"
	"database/sql"
	"github.com/jaredpetersen/go-health/health"
)

type DBDetails struct {
	ConnectionsInUse int
	ConnectionsIdle  int
}

func BuildDBHealthCheckFunc(db *sql.DB) health.CheckFunc {
	return func(ctx context.Context) health.Status {
		err := db.PingContext(ctx)
		if err != nil {
			return health.Status{State: health.StateDown}
		}

		dbStats := db.Stats()
		dbDetails := DBDetails{ConnectionsInUse: dbStats.InUse, ConnectionsIdle: dbStats.Idle}

		return health.Status{State: health.StateUp, Details: dbDetails}
	}
}
