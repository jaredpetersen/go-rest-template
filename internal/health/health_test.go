package health

import "testing"

func TestCheck(t *testing.T) {
	expectedStatus := Status{Up: true}

	status := Check()

	if *status != expectedStatus {
		t.Error("health returns incorrect status")
	}
}
