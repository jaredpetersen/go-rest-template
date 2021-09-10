package taskservice

import (
	"context"
	"errors"
	"testing"

	"github.com/jaredpetersen/go-rest-template/internal/task"
	taskmock "github.com/jaredpetersen/go-rest-template/internal/task/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetReturnsCachedTask(t *testing.T) {
	ctx := context.Background()

	storedTask := task.Task{Id: "someid"}

	tcr := taskmock.CacheClient{}
	tcr.On("Get", mock.Anything, storedTask.Id).Return(&storedTask, nil)

	tdbr := taskmock.DBClient{}

	retrievedTask, err := Get(ctx, &tcr, &tdbr, storedTask.Id)
	assert.NoError(t, err, "Returned error")
	assert.Equal(t, &storedTask, retrievedTask, "Returned incorrect task")

	tdbr.AssertExpectations(t)
	tcr.AssertExpectations(t)
}

func TestGetReturnsStoredTaskOnCacheMiss(t *testing.T) {
	ctx := context.Background()

	storedTask := task.Task{Id: "someid"}

	tcr := taskmock.CacheClient{}
	tcr.On("Get", mock.Anything, storedTask.Id).Return(nil, nil)

	tdbr := taskmock.DBClient{}
	tdbr.On("Get", mock.Anything, storedTask.Id).Return(&storedTask, nil)

	retrievedTask, err := Get(ctx, &tcr, &tdbr, storedTask.Id)
	assert.NoError(t, err, "Returned error")
	assert.Equal(t, &storedTask, retrievedTask, "Returned incorrect task")

	tcr.AssertExpectations(t)
	tdbr.AssertExpectations(t)
}

func TestGetReturnsStoredTaskOnCacheError(t *testing.T) {
	ctx := context.Background()

	storedTask := task.Task{Id: "someid"}

	tcr := taskmock.CacheClient{}
	tcr.On("Get", mock.Anything, storedTask.Id).Return(nil, errors.New("Failed"))

	tdbr := taskmock.DBClient{}
	tdbr.On("Get", mock.Anything, storedTask.Id).Return(&storedTask, nil)

	retrievedTask, err := Get(ctx, &tcr, &tdbr, storedTask.Id)
	assert.NoError(t, err, "Returned error")
	assert.Equal(t, &storedTask, retrievedTask, "Returned incorrect task")

	tcr.AssertExpectations(t)
	tdbr.AssertExpectations(t)
}

func TestGetReturnsErrorOnDBError(t *testing.T) {
	ctx := context.Background()

	storedTask := task.Task{Id: "someid"}
	dbErr := errors.New("Fail")

	tcr := taskmock.CacheClient{}
	tcr.On("Get", mock.Anything, storedTask.Id).Return(nil, nil)

	tdbr := taskmock.DBClient{}
	tdbr.On("Get", mock.Anything, storedTask.Id).Return(nil, dbErr)

	retrievedTask, err := Get(ctx, &tcr, &tdbr, storedTask.Id)
	assert.ErrorIs(t, dbErr, err, "Incorrect error")
	assert.Nil(t, retrievedTask, "Task must be nil")

	tcr.AssertExpectations(t)
	tdbr.AssertExpectations(t)
}

func TestSave(t *testing.T) {
	ctx := context.Background()

	tsk := task.Task{Id: "someid"}

	tcr := taskmock.CacheClient{}
	tcr.On("Save", mock.Anything, tsk).Return(nil)

	tdb := taskmock.DBClient{}
	tdb.On("Save", mock.Anything, tsk).Return(nil)

	err := Save(ctx, &tcr, &tdb, tsk)
	assert.NoError(t, err, "Returned error")

	tdb.AssertExpectations(t)
	tcr.AssertExpectations(t)
}

func TestSaveOnCacheError(t *testing.T) {
	ctx := context.Background()

	tsk := task.Task{Id: "someid"}

	tcr := taskmock.CacheClient{}
	tcr.On("Save", mock.Anything, tsk).Return(errors.New("Failed"))

	tdb := taskmock.DBClient{}
	tdb.On("Save", mock.Anything, tsk).Return(nil)

	err := Save(ctx, &tcr, &tdb, tsk)
	assert.NoError(t, err, "Returned error")

	tdb.AssertExpectations(t)
	tcr.AssertExpectations(t)
}

func TestSaveReturnsErrorOnDBError(t *testing.T) {
	ctx := context.Background()

	tsk := task.Task{Id: "someid"}
	dbErr := errors.New("Failed")

	tcr := taskmock.CacheClient{}
	tcr.On("Save", mock.Anything, tsk).Return(nil)

	tdb := taskmock.DBClient{}
	tdb.On("Save", mock.Anything, tsk).Return(dbErr)

	err := Save(ctx, &tcr, &tdb, tsk)
	assert.ErrorIs(t, dbErr, err, "Incorrect error")

	tdb.AssertExpectations(t)
	tcr.AssertExpectations(t)
}
