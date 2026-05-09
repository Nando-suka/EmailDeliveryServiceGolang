package queue

import (
	"time"

	"github.com/Nando-suka/email-service/model"
)

type MockQueue struct {
	Tasks []model.EmailTask
	Err   error
}

func (m *MockQueue) Enqueue(task model.EmailTask) error {
	if m.Err != nil {
		return m.Err
	}
	m.Tasks = append(m.Tasks, task)
	return nil
}

func (m *MockQueue) Dequeue(timeout time.Duration) (*model.EmailTask, error) {
	// tidak diperlukan untuk handler
	return nil, nil
}
