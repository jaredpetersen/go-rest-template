package task

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	redismock "github.com/jaredpetersen/go-rest-template/internal/redis/mocks"
)

func TestCacheRepoSave(t *testing.T) {
	ctx := context.Background()

	tsk := New()

	rdb := redismock.Client{}
	rdb.On("Set", mock.Anything, "task."+tsk.Id, mock.MatchedBy(taskMatcher(*tsk)), time.Duration(0)).Return(nil)

	tcr := CacheRepo{Redis: &rdb}

	err := tcr.Save(ctx, *tsk)
	assert.NoError(t, err, "Returned error")

	rdb.AssertExpectations(t)
}

func TestCacheRepoSaveReturnsRedisError(t *testing.T) {
	ctx := context.Background()

	task := New()

	expectedErr := errors.New("Failed")

	rdb := redismock.Client{}
	rdb.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	tcr := CacheRepo{Redis: &rdb}

	err := tcr.Save(ctx, *task)
	assert.EqualError(t, err, expectedErr.Error(), "Did not return error")
}

func TestCacheRepoGet(t *testing.T) {
	ctx := context.Background()

	id := "2b7e1292-a831-4df5-b00e-3105a51111bb"
	description := "buy socks"
	storedTask := fmt.Sprintf("{\"description\":\"%s\"}", description)

	rdb := redismock.Client{}
	rdb.On("Get", mock.Anything, "task."+id).Return(&storedTask, nil)

	tcr := CacheRepo{Redis: &rdb}

	tsk, err := tcr.Get(ctx, id)
	assert.NoError(t, err, "Returned error")
	assert.NotEqual(t, &tsk, Task{Description: description}, "Task is incorrect")

	rdb.AssertExpectations(t)
}

func TestCacheRepoGetNotExists(t *testing.T) {
	ctx := context.Background()

	id := "868e5655-660e-41f1-b271-b00172d7fa2d"

	rdb := redismock.Client{}
	rdb.On("Get", mock.Anything, "task."+id).Return(nil, nil)

	tcr := CacheRepo{Redis: &rdb}

	tsk, err := tcr.Get(ctx, id)
	assert.NoError(t, err, "Returned error")
	assert.Nil(t, tsk, "Task should be nil")

	rdb.AssertExpectations(t)
}

func TestCacheRepoGetReturnsRedisError(t *testing.T) {
	ctx := context.Background()

	id := "f48d0dbd-4cfe-4295-95cd-55cff8f5c149"

	expectedError := errors.New("Failed")

	rdb := redismock.Client{}
	rdb.On("Get", mock.Anything, mock.Anything).Return(nil, expectedError)

	tcr := CacheRepo{Redis: &rdb}

	tsk, err := tcr.Get(ctx, id)
	assert.Nil(t, tsk, "Task should be nil")
	assert.EqualError(t, err, expectedError.Error(), "Did not return error")
}
