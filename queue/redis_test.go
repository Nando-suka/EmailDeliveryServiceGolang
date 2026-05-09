package queue

import (
	"testing"
	"time"

	"github.com/Nando-suka/email-service/model"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisQueue_EnqueueDequeue(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	q := NewRedisQueue(mr.Addr(), "", "email:queue")

	task := model.EmailTask{
		ID:      "123",
		To:      []string{"test@test.com"},
		Subject: "Hello",
		Body:    "Body",
	}
	err := q.Enqueue(task)
	assert.NoError(t, err)

	dequeued, err := q.Dequeue(1 * time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "123", dequeued.ID)
}
