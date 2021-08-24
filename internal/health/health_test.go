package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	status := Check()

	assert.Equal(t, Status{Up: true}, *status, "Status is incorrect")
}
