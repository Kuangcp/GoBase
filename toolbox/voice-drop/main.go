package main

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

var hub = NewHub()

//go:embed frontend
var staticFiles embed.FS

var staticFS = func() fs.FS {
	s, _ := fs.Sub(staticFiles, "frontend")
	return s
}()

type sendRequest struct {
	Text string `json:"text"`
}

type sendResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func main() {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("POST /api/send", handleSend)
	httpMux.HandleFunc("POST /api/start-record", handleStartRecord)
	httpMux.HandleFunc("POST /api/stop-record", handleStopRecord)

	go func() {
		log.Printf("http api on :6600")
		log.Fatal(http.ListenAndServe(":6600", withCORS(httpMux)))
	}()

	certFile, keyFile, err := certFiles()
	if err != nil {
		log.Fatalf("cert: %v", err)
	}

	httpsMux := http.NewServeMux()
	httpsMux.Handle("GET /", http.FileServerFS(staticFS))
	httpsMux.HandleFunc("POST /api/send", handleSend)
	httpsMux.HandleFunc("POST /api/start-record", handleStartRecord)
	httpsMux.HandleFunc("POST /api/stop-record", handleStopRecord)
	httpsMux.HandleFunc("GET /ws", hub.HandleWS)

	log.Printf("https on :6601")
	log.Fatal(http.ListenAndServeTLS(":6601", certFile, keyFile, withCORS(httpsMux)))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, sendResponse{Ok: false, Error: "invalid JSON"})
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		writeJSON(w, http.StatusBadRequest, sendResponse{Ok: false, Error: "text is empty"})
		return
	}
	pasteText(req.Text)
	writeJSON(w, http.StatusOK, sendResponse{Ok: true})
}

func handleStartRecord(w http.ResponseWriter, r *http.Request) {
	if err := hub.Send("start_recording"); err != nil {
		msg := "phone not connected"
		if !errors.Is(err, errNoClient) {
			msg = "send failed"
		}
		writeJSON(w, http.StatusServiceUnavailable, sendResponse{Ok: false, Error: msg})
		return
	}
	writeJSON(w, http.StatusOK, sendResponse{Ok: true})
}

func handleStopRecord(w http.ResponseWriter, r *http.Request) {
	if err := hub.Send("stop_recording"); err != nil {
		msg := "phone not connected"
		if !errors.Is(err, errNoClient) {
			msg = "send failed"
		}
		writeJSON(w, http.StatusServiceUnavailable, sendResponse{Ok: false, Error: msg})
		return
	}
	writeJSON(w, http.StatusOK, sendResponse{Ok: true})
}

func xenv() []string {
	env := os.Environ()
	display := os.Getenv("DISPLAY")
	if display == "" {
		display = ":0"
	}
	out := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, "DISPLAY=") {
			out = append(out, e)
		}
	}
	return append(out, "DISPLAY="+display)
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
