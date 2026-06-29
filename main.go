package main

import (
	"log"
	"net/http"

	"github.com/AnnaTarantina/APIGateway/handlers"
)

func main() {
	// Эндпоинты для новостей (заглушки)
	http.HandleFunc("/news", handlers.GetNewsList)
	http.HandleFunc("/news/filter", handlers.FilterNews)
	http.HandleFunc("/news/detail", handlers.GetDetailedNews)

	// Эндпоинты для комментариев (прокси в CommentService)
	http.HandleFunc("/comments", handlers.CommentsHandler)
	http.HandleFunc("/comments/by-news-id", handlers.GetCommentsByNewsID)

	log.Println("API Gateway is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
