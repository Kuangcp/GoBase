package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

//go:embed frontend
var frontendFS embed.FS

type sendRequest struct {
	Text string `json:"text"`
}

type sendResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6601"
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(frontendFS))
	mux.HandleFunc("POST /api/send", handleSend)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, withCORS(mux)))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, sendResponse{
			Ok:    false,
			Error: "invalid JSON",
		})
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		writeJSON(w, http.StatusBadRequest, sendResponse{
			Ok:    false,
			Error: "text is empty",
		})
		return
	}

	clipCmd := exec.Command("xclip", "-selection", "clipboard")
	clipCmd.Stdin = strings.NewReader(req.Text)
	if out, err := clipCmd.CombinedOutput(); err != nil {
		log.Printf("xclip error: %v, output: %s", err, string(out))
		writeJSON(w, http.StatusInternalServerError, sendResponse{
			Ok:    false,
			Error: "failed to copy to clipboard",
		})
		return
	}

	pasteCmd := exec.Command("xdotool", "key", "ctrl+v")
	if out, err := pasteCmd.CombinedOutput(); err != nil {
		log.Printf("xdotool key error: %v, output: %s", err, string(out))
		writeJSON(w, http.StatusInternalServerError, sendResponse{
			Ok:    false,
			Error: "failed to paste",
		})
		return
	}

	log.Printf("pasted %d characters", len(req.Text))
	writeJSON(w, http.StatusOK, sendResponse{Ok: true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
