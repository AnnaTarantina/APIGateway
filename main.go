package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnnaTarantina/APIGateway/handlers"
	"github.com/AnnaTarantina/APIGateway/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Новости
	mux.HandleFunc("/news", handlers.GetNewsList)
	mux.HandleFunc("/news/filter", handlers.FilterNews)
	mux.HandleFunc("/news/detail", handlers.GetDetailedNews)

	// Комментарии
	mux.HandleFunc("/comments", handlers.CommentsHandler)
	mux.HandleFunc("/comments/by-news-id", handlers.GetCommentsByNewsID)

	handler := middleware.RequestIDMiddleware(middleware.LoggingMiddleware(mux))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Запуск сервера в горутине
	go func() {
		log.Println("[*] API Gateway HTTP server is started on localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("[*] API Gateway HTTP server has been stopped. Reason: got %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
