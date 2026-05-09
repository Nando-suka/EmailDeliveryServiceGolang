package worker

import (
	"errors"
	"testing"
	"time"

	"github.com/Nando-suka/email-service/config"
	"github.com/Nando-suka/email-service/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
)

type MockDialer struct {
    Err error
}

func (m *MockDialer) DialAndSend(msg *gomail.Message) error {
	return m.Err
}

type MockQueue struct {
    Tasks []model.EmailTask
    Err   error
}

func (m *MockQueue) Enqueue(task model.EmailTask) error {
    m.Tasks = append(m.Tasks, task)
    return m.Err
}

func (m *MockQueue) Dequeue(timeout time.Duration) (*model.EmailTask, error) {
    if len(m.Tasks) == 0 {
        // blocking? untuk test kita return nil, time.Sleep sebentar
        time.Sleep(10*time.Millisecond)
        return nil, errors.New("queue empty")
    }
    task := m.Tasks[0]
    m.Tasks = m.Tasks[1:]
    return &task, nil
}

func TestSender_ProcessTask_Success(t *testing.T) {
    cfg := &config.Config{FromEmail: "test@example.com", FromName: "Test"}
    mockQ := &MockQueue{
        Tasks: []model.EmailTask{
            {ID: "1", To: []string{"to@example.com"}, Subject: "Sub", Body: "Body", MaxRetries: 3},
        },
    }
    dialer := &MockDialer{} // sukses, error nil
    s := &Sender{cfg: cfg, queue: mockQ, dialer: dialer, quit: make(chan bool)}
    // jalankan processTask langsung
    s.processTask(1, &mockQ.Tasks[0]) // task akan dikirim
	// processTask langsung tidak mengeluarkan task dari queue; yang dicek adalah tidak ada retry enqueue tambahan.
	assert.Len(t, mockQ.Tasks, 1)
}

func TestSender_ProcessTask_FailThenRetry(t *testing.T) {
	cfg := &config.Config{FromEmail: "test@example.com", FromName: "Test"}
	q := &MockQueue{}
	dialer := &MockDialer{Err: errors.New("smtp error")}
	s := &Sender{cfg: cfg, queue: q, dialer: dialer, quit: make(chan bool)}
	task := model.EmailTask{ID: "2", To: []string{"to"}, Subject: "S", Body: "B", MaxRetries: 3}
	s.processTask(1, &task)
	// harusnya retry: task di-enqueue lagi dengan retries=1
	assert.Len(t, q.Tasks, 1)
	assert.Equal(t, 1, q.Tasks[0].Retries)
}