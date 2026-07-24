package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "6601"
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(staticFS))
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

	env := xenv()

	clipCmd := exec.Command("xclip", "-selection", "clipboard")
	clipCmd.Env = env
	clipCmd.Stdin = strings.NewReader(req.Text)
	if err := clipCmd.Start(); err != nil {
		log.Printf("xclip start error: %v", err)
		writeJSON(w, http.StatusInternalServerError, sendResponse{Ok: false, Error: "failed to start xclip"})
		return
	}

	time.Sleep(200 * time.Millisecond)

	pasteCmd := exec.Command("xdotool", "key", "ctrl+v")
	pasteCmd.Env = env
	if out, err := pasteCmd.CombinedOutput(); err != nil {
		clipCmd.Process.Kill()
		log.Printf("xdotool error: %v, output: %s", err, string(out))
		writeJSON(w, http.StatusInternalServerError, sendResponse{Ok: false, Error: "failed to paste"})
		return
	}

	clipDone := make(chan error, 1)
	go func() { clipDone <- clipCmd.Wait() }()
	select {
	case <-clipDone:
	case <-time.After(time.Second):
		clipCmd.Process.Kill()
	}

	log.Printf("pasted %d characters", len(req.Text))
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
