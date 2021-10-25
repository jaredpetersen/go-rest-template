package healthcheck

import (
	"context"
	"github.com/jaredpetersen/go-health/health"
	"github.com/jaredpetersen/go-rest-template/internal/redis"
)

func BuildRedisHealthCheckFunc(rdb redis.Client) health.CheckFunc {
	return func(ctx context.Context) health.Status {
		err := rdb.Ping(ctx)
		if err != nil {
			// In our case, Redis is just a caching layer to improve performance and a failed connection does not mean
			// that the application is down. Return a warning status instead so that the outage is visible but does not
			// trigger application restarts.
			return health.Status{State: health.StateWarn}
		}

		return health.Status{State: health.StateUp}
	}
}
