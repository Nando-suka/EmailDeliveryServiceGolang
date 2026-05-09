package queue

import (
	"errors"
	"time"

	"github.com/Nando-suka/email-service/model"
)

type InMemQueue struct {
	ch chan model.EmailTask
}

func NewInMemQueue(buffer int) *InMemQueue {
	if buffer <= 0 {
		buffer = 1
	}
	return &InMemQueue{ch: make(chan model.EmailTask, buffer)}
}

func (q *InMemQueue) Enqueue(task model.EmailTask) error {
	q.ch <- task
	return nil
}

func (q *InMemQueue) Dequeue(timeout time.Duration) (*model.EmailTask, error) {
	select {
	case task := <-q.ch:
		return &task, nil
	case <-time.After(timeout):
		return nil, errors.New("queue timeout")
	}
}
