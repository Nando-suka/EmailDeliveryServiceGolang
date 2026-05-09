package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nando-suka/email-service/model"
	"github.com/Nando-suka/email-service/queue"
	"github.com/google/uuid"
)

type EmailHandler struct {
	queue      *queue.RedisQueue
	maxRetries int
}

func NewEmailHandler(q *queue.RedisQueue, maxRetries int) *EmailHandler {
	return &EmailHandler{queue: q, maxRetries: maxRetries}
}

type SendRequest struct {
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	ContentType string   `json:"content_type"`
}

func (h *EmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if len(req.To) == 0 || req.Subject == "" || req.Body == "" {
		http.Error(w, "to, subject, and body are required", http.StatusBadRequest)
		return
	}

	task := model.EmailTask{
		ID:          uuid.New().String(),
		To:          req.To,
		Subject:     req.Subject,
		Body:        req.Body,
		ContentType: req.ContentType,
		CreatedAt:   time.Now(),
		MaxRetries:  h.maxRetries,
	}

	if err := h.queue.Enqueue(task); err != nil {
		http.Error(w, "failed to queue email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "email queued",
		"id":      task.ID,
	})
}
