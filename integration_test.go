package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nando-suka/email-service/config"
	"github.com/Nando-suka/email-service/handler"
	"github.com/Nando-suka/email-service/queue"
	"github.com/Nando-suka/email-service/worker"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_EmailFlow(t *testing.T) {
	// Setup config dengan Mailpit
	cfg := &config.Config{
		SMTPHost:     "localhost",
		SMTPPort:     1025,
		SMTPUser:     "",
		SMTPPassword: "",
		FromEmail:    "test@service.com",
		FromName:     "Test Service",
		WorkerCount:  1,
		MaxRetries:   2,
	}
	// Gunakan queue in-memory
	q := queue.NewInMemQueue(100)
	dialer := &worker.RealDialer{
		Host: cfg.SMTPHost,
		Port: cfg.SMTPPort,
		User: cfg.SMTPUser,
		Pass: cfg.SMTPPassword,
	}
	sender := worker.NewSender(cfg, q, dialer)
	go sender.Start(1)
	defer sender.Stop()

	// Handler
	h := handler.NewEmailHandler(q, cfg.MaxRetries)

	// Kirim HTTP request
	body := map[string]interface{}{
		"to":      []string{"penerima@example.com"},
		"subject": "Integration Test",
		"body":    "<p>Halo dari integration test</p>",
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

	// Tunggu worker mengirim
	time.Sleep(2 * time.Second)

	// Verifikasi email masuk di Mailpit API
	client := http.Client{}
	req, _ = http.NewRequest("GET", "http://localhost:8025/api/v1/messages", nil)
	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()
	var mailpitMessages struct {
		Messages []struct {
			Subject string   `json:"subject"`
			To      []string `json:"to"`
		} `json:"messages"`
	}

	json.NewDecoder(res.Body).Decode(&mailpitMessages)
	assert.Greater(t, len(mailpitMessages.Messages), 0)
	// Cari email kita (opsional)
	found := false
	for _, m := range mailpitMessages.Messages {
		if m.Subject == "Integration Test" {
			found = true
			break
		}
	}
	assert.True(t, found, "Email harus ditemukan di Mailpit")

}
