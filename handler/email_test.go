package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nando-suka/email-service/model"
	"github.com/stretchr/testify/assert"
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

func TestSendEmail_Success(t *testing.T) {
	mockQ := &MockQueue{}
	h := NewEmailHandler(mockQ, 3)

	body := map[string]interface{}{
		"to":      []string{"user@example.com"},
		"subject": "Test",
		"body":    "<h1>Test</h1>",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/emails", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.SendEmail(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["id"])

	assert.Equal(t, 1, len(mockQ.Tasks))
	assert.Equal(t, "user@example.com", mockQ.Tasks[0].To[0])
}

func TestSendEmail_InvalidBody(t *testing.T) {
	mockQ := &MockQueue{}
	h := NewEmailHandler(mockQ, 3)

	req := httptest.NewRequest("POST", "/api/emails", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()
	h.SendEmail(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
