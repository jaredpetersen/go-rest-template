package task

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

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

func TestSave(t *testing.T) {
	ctx := context.Background()

	task := New()

	rdb := redisMock.Client{}
	rdb.On("Set", mock.Anything, "task."+task.Id, mock.MatchedBy(taskMatcher(*task)), mock.Anything).Return(nil)

	err := Save(ctx, &rdb, *task)
	if err != nil {
		t.Error("Encountered error", err)
	}

	rdb.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	id := "2b7e1292-a831-4df5-b00e-3105a51111bb"
	description := "buy socks"
	storedTask := fmt.Sprintf("{\"description\":\"%s\"}", description)

	rdb := redisMock.Client{}
	rdb.On("Get", mock.Anything, "task."+id).Return(&storedTask, nil)

	task, err := Get(ctx, &rdb, id)
	if err != nil {
		t.Error("Encountered error", err)
	}

	if *task != (Task{Description: description}) {
		t.Error("Invalid task")
	}

	rdb.AssertExpectations(t)
}
