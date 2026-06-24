package main

import (
	"log"
	"net/http"

	"./handlers"
)

func main() {
	http.HandleFunc("/news", handlers.GetNewsList)
	http.HandleFunc("/news/filter", handlers.FilterNews)
	http.HandleFunc("/news/detail", handlers.GetDetailedNews)
	http.HandleFunc("/comments", handlers.AddComment)
	http.HandleFunc("/comments/by-news-id", handlers.GetCommentsByNewsID)

	log.Println("API Gateway is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
