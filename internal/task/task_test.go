package task

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func taskMatcher(expectedTask Task) func(value []byte) bool {
	return func(value []byte) bool {
		var unmarshaledTask Task
		err := json.Unmarshal(value, &unmarshaledTask)
		if err != nil {
			return false
		}

		return cmp.Equal(expectedTask, unmarshaledTask)
	}
}

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return (err == nil)
}

func TestNew(t *testing.T) {
	tsk := New()

	assert.True(t, isValidUUID(tsk.ID), "Did not generate valid ID")

	assert.Nil(t, tsk.DateDue, "Initialized DateDue")

	assert.False(t, tsk.DateCreated.IsZero(), "Did not initialize DateCreated")

	assert.False(t, tsk.DateUpdated.IsZero(), "Did not initialize DateUpdated")

	assert.Equal(t, tsk.DateCreated, tsk.DateUpdated, "DateCreated and DateUpdated are not equal")

	expectedTask := Task{ID: tsk.ID, DateCreated: tsk.DateCreated, DateUpdated: tsk.DateUpdated}
	assert.Equal(t, expectedTask, *tsk, "Task is setting more defaults than expected")
}
