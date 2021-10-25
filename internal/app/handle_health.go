package app

import (
	"github.com/jaredpetersen/go-health/health"
	"github.com/jaredpetersen/go-rest-template/api"
	"net/http"
)

// handleLiveness creates a HTTP handler that indicates when the application is alive or dead
func (a *app) handleLiveness() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	}
}

// handleReadiness creates a HTTP handler that indicates when the application is ready to serve traffic
func (a *app) handleReadiness() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		monitorStatus := a.HealthMonitor.Check()

		redisStatus := monitorStatus.CheckStatuses["redis"]
		dbStatus := monitorStatus.CheckStatuses["database"]

		res := api.Health{
			State: transformState(monitorStatus.State),
			Components: api.HealthComponents{
				Redis: api.HealthComponent{
					State: transformState(redisStatus.Status.State),
					Timestamp: redisStatus.Timestamp,
				},
				CockroachDb: api.HealthComponent{
					State: transformState(dbStatus.Status.State),
					Timestamp: redisStatus.Timestamp,
				},
			},
		}

		var statusCode int
		if res.State == api.HealthStateDOWN {
			statusCode = http.StatusServiceUnavailable
		} else {
			statusCode = http.StatusOK
		}

		respond(w, res, statusCode)
	}
}

func transformState(monitorState health.State) api.HealthState {
	var state api.HealthState
	if monitorState == health.StateUp {
		state = api.HealthStateUP
	} else if monitorState == health.StateWarn {
		state = api.HealthStateWARN
	} else {
		state = api.HealthStateDOWN
	}

	return state
}
