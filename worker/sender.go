package worker

import (
	"log"
	"time"

	"github.com/Nando-suka/email-service/config"
	"github.com/Nando-suka/email-service/model"
	"github.com/Nando-suka/email-service/queue"
	"gopkg.in/gomail.v2"
)

type Sender struct {
	cfg   *config.Config
	queue *queue.RedisQueue
	quit  chan bool
}

func NewSender(cfg *config.Config, q *queue.RedisQueue) *Sender {
	return &Sender{
		cfg:   cfg,
		queue: q,
		quit:  make(chan bool),
	}
}

func (s *Sender) Start(workerID int) {
	log.Printf("[Worker %d] dimulai", workerID)
	for {
		select {
		case <-s.quit:
			return
		default:
			task, err := s.queue.Dequeue(5 * time.Second)
			if err != nil {
				// Jika timeout/no task, lanjutkan loop
				continue
			}
			s.processTask(workerID, task)
		}
	}
}

func (s *Sender) processTask(workerID int, task *model.EmailTask) {
	log.Printf("[Worker %d] mengirim email %s ke %v", workerID, task.ID, task.To)
	err := s.sendEmail(task)
	if err != nil {
		log.Printf("[Worker %d] gagal kirim email %s: %v", workerID, task.ID, err)
		task.Retries++
		if task.Retries < task.MaxRetries {
			// Requefue untuk retry (bisa dengan delay)
			time.Sleep(2 * time.Second) // delay sederhana
			if err := s.queue.Enqueue(*task); err != nil {
				log.Printf("[Worker %d] gagal me-reenqueue email %s: %v", workerID, task.ID, err)
			}
		} else {
			log.Printf("[Worker %d] email %s gagal setelah %d kali percobaan", workerID, task.ID, task.MaxRetries)
			// Bisa simpan ke dead letter queue atau logging khusus
		}
	} else {
		log.Printf("[Worker %d] email %s berhasil dikirim", workerID, task.ID)
	}
}

func (s *Sender) sendEmail(task *model.EmailTask) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.cfg.FromEmail, s.cfg.FromName))
	m.SetHeader("To", task.To...)
	m.SetHeader("Subject", task.Subject)
	contentType := "text/html"
	if task.ContentType != "" {
		contentType = task.ContentType
	}
	m.SetBody(contentType, task.Body)

	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword)
	return d.DialAndSend(m)
}

func (s *Sender) Stop() {
	close(s.quit)
}
