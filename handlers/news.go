package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/AnnaTarantina/APIGateway/middleware"
	"github.com/AnnaTarantina/APIGateway/models"
)

const newsServiceURL = "http://localhost:8082"

func GetNewsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := ctx.Value(middleware.RequestIDKey).(string)

	targetURL := fmt.Sprintf("%s/news/100?request_id=%s", newsServiceURL, reqID)
	if s := r.URL.Query().Get("s"); s != "" {
		targetURL += "&s=" + s
	}
	if page := r.URL.Query().Get("page"); page != "" {
		targetURL += "&page=" + page
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	req.Header.Set("X-Request-ID", reqID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "news service unavailable"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func GetDetailedNews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	newsID := r.URL.Query().Get("id")
	if newsID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	reqID, _ := ctx.Value(middleware.RequestIDKey).(string)
	results := make(chan interface{}, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	// Горутина 1: Новость из GONEWS
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("%s/news/detail/%s?request_id=%s", newsServiceURL, newsID, reqID)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		req.Header.Set("X-Request-ID", reqID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			results <- fmt.Errorf("news service error: %v", err)
			return
		}
		defer resp.Body.Close()

		var news models.NewsFullDetailed
		json.NewDecoder(resp.Body).Decode(&news)
		results <- news
	}()

	// Горутина 2: Комментарии из CommentService
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("http://localhost:8081/comments?news_id=%s&request_id=%s", newsID, reqID)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		req.Header.Set("X-Request-ID", reqID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			results <- fmt.Errorf("comments service error: %v", err)
			return
		}
		defer resp.Body.Close()

		var comments []models.Comment
		json.NewDecoder(resp.Body).Decode(&comments)
		results <- comments
	}()

	wg.Wait()
	close(results)

	var finalNews models.NewsFullDetailed
	var finalComments []models.Comment
	var fetchErr error

	for res := range results {
		switch v := res.(type) {
		case error:
			fetchErr = v
		case models.NewsFullDetailed:
			finalNews = v
		case []models.Comment:
			finalComments = v
		}
	}

	if fetchErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": fetchErr.Error()})
		return
	}

	finalNews.Comments = finalComments
	writeJSON(w, http.StatusOK, finalNews)
}

func FilterNews(w http.ResponseWriter, r *http.Request) {
	GetNewsList(w, r)
}
