package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var errNoClient = errors.New("no ws client connected")

type wsMsg struct {
	Cmd   string `json:"cmd,omitempty"`
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
}

type Hub struct {
	mu   sync.Mutex
	conn *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade: %v", err)
		return
	}

	h.mu.Lock()
	if h.conn != nil {
		h.conn.Close()
	}
	h.conn = conn
	h.mu.Unlock()

	log.Println("ws client connected")

	defer func() {
		h.mu.Lock()
		h.conn = nil
		h.mu.Unlock()
		conn.Close()
		log.Println("ws client disconnected")
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg wsMsg
		if json.Unmarshal(data, &msg) != nil {
			continue
		}
		if msg.Type == "transcript" {
			text := strings.TrimSpace(msg.Text)
			if text == "" {
				continue
			}
			log.Printf("ws transcript: %s", text)
			pasteText(text)
		}
	}
}

func (h *Hub) Send(cmd string) error {
	h.mu.Lock()
	conn := h.conn
	h.mu.Unlock()
	if conn == nil {
		return errNoClient
	}
	return conn.WriteJSON(wsMsg{Cmd: cmd})
}

func (h *Hub) Connected() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.conn != nil
}

func execCmd(env []string, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	return cmd
}

func pasteText(text string) {
	env := xenv()

	clipCmd := execCmd(env, "xclip", "-selection", "clipboard")
	clipCmd.Stdin = strings.NewReader(text)
	if err := clipCmd.Start(); err != nil {
		log.Printf("xclip start: %v", err)
		return
	}

	time.Sleep(200 * time.Millisecond)

	pasteCmd := execCmd(env, "xdotool", "key", "ctrl+v")
	if out, err := pasteCmd.CombinedOutput(); err != nil {
		clipCmd.Process.Kill()
		log.Printf("xdotool error: %v, output: %s", err, string(out))
		return
	}

	done := make(chan error, 1)
	go func() { done <- clipCmd.Wait() }()
	select {
	case <-done:
	case <-time.After(time.Second):
		clipCmd.Process.Kill()
	}

	log.Printf("pasted %d characters", len(text))
}
