package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AnnaTarantina/APIGateway/middleware"
)

const (
	commentServiceURL    = "http://localhost:8081"
	censorshipServiceURL = "http://localhost:8083"
)

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddComment(w, r)
	case http.MethodGet:
		GetCommentsByNewsID(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := ctx.Value(middleware.RequestIDKey).(string)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read body"})
		return
	}
	defer r.Body.Close()

	// 1. СИНХРОННЫЙ ЗАПРОС К СЕРВИСУ ЦЕНЗУРИРОВАНИЯ
	var commentReq struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(body, &commentReq); err == nil {
		checkBody, _ := json.Marshal(map[string]string{"text": commentReq.Text})
		checkURL := fmt.Sprintf("%s/check?request_id=%s", censorshipServiceURL, reqID)

		checkReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, checkURL, bytes.NewReader(checkBody))
		checkReq.Header.Set("Content-Type", "application/json")
		checkReq.Header.Set("X-Request-ID", reqID)

		checkResp, err := http.DefaultClient.Do(checkReq)
		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": "censorship service unavailable"})
			return
		}
		defer checkResp.Body.Close()

		if checkResp.StatusCode != http.StatusOK {
			// Пробрасываем ошибку от сервиса цензуры (400)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(checkResp.StatusCode)
			io.Copy(w, checkResp.Body)
			return
		}
	}

	// 2. ПРОКСИРОВАНИЕ В COMMENTS SERVICE (только если цензура пройдена)
	targetURL := fmt.Sprintf("%s/comment?request_id=%s", commentServiceURL, reqID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", reqID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "comment service unavailable"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func GetCommentsByNewsID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := ctx.Value(middleware.RequestIDKey).(string)
	newsID := r.URL.Query().Get("news_id")

	if newsID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "news_id is required"})
		return
	}

	targetURL := fmt.Sprintf("%s/comments?news_id=%s&request_id=%s", commentServiceURL, newsID, reqID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	req.Header.Set("X-Request-ID", reqID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "comment service unavailable"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
