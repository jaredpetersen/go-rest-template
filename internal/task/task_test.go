package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	redisMock "github.com/jaredpetersen/go-rest-template/internal/redis/mocks"
)

func taskMatcher(expectedTask Task) func(value []byte) bool {
	return func(value []byte) bool {
		var unmarshaledTask Task
		err := json.Unmarshal(value, &unmarshaledTask)
		if err != nil {
			return false
		}

		return expectedTask == unmarshaledTask
	}
}

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return (err == nil)
}

func isValidDatetime(datetime string) bool {
	_, err := time.Parse(time.RFC3339, datetime)
	return (err == nil)
}

func TestNew(t *testing.T) {
	task := New()

	assert.True(t, isValidUUID(task.Id), "Did not generate valid ID")

	assert.True(t, isValidDatetime(task.DateCreated), "Did not generate valid datetime string for DateCreated")

	assert.True(t, isValidDatetime(task.DateUpdated), "Did not generate valid datetime string for DateUpdated")

	assert.Equal(t, task.DateCreated, task.DateUpdated, "DateCreated and DateUpdated are not equal")

	expectedTask := Task{Id: task.Id, DateCreated: task.DateCreated, DateUpdated: task.DateUpdated}
	assert.Equal(t, expectedTask, *task, "Task is setting more defaults than expected")
}

func TestSave(t *testing.T) {
	ctx := context.Background()

	task := New()

	rdb := redisMock.Client{}
	rdb.On("Set", mock.Anything, "task."+task.Id, mock.MatchedBy(taskMatcher(*task)), time.Duration(0)).Return(nil)

	err := Save(ctx, &rdb, *task)
	assert.NoError(t, err, "Returned error")

	rdb.AssertExpectations(t)
}

func TestSaveReturnsRedisError(t *testing.T) {
	ctx := context.Background()

	task := New()

	expectedErr := errors.New("Failed")

	rdb := redisMock.Client{}
	rdb.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	err := Save(ctx, &rdb, *task)
	assert.EqualError(t, err, expectedErr.Error(), "Did not return error")
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	id := "2b7e1292-a831-4df5-b00e-3105a51111bb"
	description := "buy socks"
	storedTask := fmt.Sprintf("{\"description\":\"%s\"}", description)

	rdb := redisMock.Client{}
	rdb.On("Get", mock.Anything, "task."+id).Return(&storedTask, nil)

	task, err := Get(ctx, &rdb, id)
	assert.NoError(t, err, "Returned error")
	assert.NotEqual(t, &task, Task{Description: description}, "Task is incorrect")

	rdb.AssertExpectations(t)
}

func TestGetNotExists(t *testing.T) {
	ctx := context.Background()

	id := "868e5655-660e-41f1-b271-b00172d7fa2d"

	rdb := redisMock.Client{}
	rdb.On("Get", mock.Anything, "task."+id).Return(nil, nil)

	task, err := Get(ctx, &rdb, id)
	assert.NoError(t, err, "Returned error")
	assert.Nil(t, task, "Task should be nil")

	rdb.AssertExpectations(t)
}

func TestGetReturnsRedisError(t *testing.T) {
	ctx := context.Background()

	id := "f48d0dbd-4cfe-4295-95cd-55cff8f5c149"

	expectedError := errors.New("Failed")

	rdb := redisMock.Client{}
	rdb.On("Get", mock.Anything, mock.Anything).Return(nil, expectedError)

	task, err := Get(ctx, &rdb, id)
	assert.Nil(t, task, "Task should be nil")
	assert.EqualError(t, err, expectedError.Error(), "Did not return error")
}
