package queue

import (
	"time"

	"github.com/Nando-suka/email-service/model"
)

type Queue interface {
	Enqueue(task model.EmailTask) error
	Dequeue(timeout time.Duration) (*model.EmailTask, error)
}
