package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nando-suka/email-service/config"
	"github.com/Nando-suka/email-service/handler"
	"github.com/Nando-suka/email-service/queue"
	"github.com/Nando-suka/email-service/worker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	// Inisialisasi antrean
	q := queue.NewRedisQueue(cfg.RedisAddr, cfg.RedisPass, "email:queue")

	// Router HTTP
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	emailHandler := handler.NewEmailHandler(q, cfg.MaxRetries)
	r.Post("/api/emails", emailHandler.SendEmail)

	// Jalankan worker pool
	sender := worker.NewSender(cfg, q)
	for i := 1; i <= cfg.WorkerCount; i++ {
		go sender.Start(i)
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Mematikan worker...")
		sender.Stop()
		os.Exit(0)
	}()

	log.Println("Server berjalan di  :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
