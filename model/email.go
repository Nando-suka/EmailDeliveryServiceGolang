package model

import "time"

type EmailTask struct {
	ID          string    `json:"id"`
	To          []string  `json:"to"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`         // HTML body
	ContentType string    `json:"content_type"` // text/html atau text/plain
	CreatedAt   time.Time `json:"created_at"`
	Retries     int       `json:"retries"`
	MaxRetries  int       `json:"max_retries"`
}
