package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Delta struct {
	Content  string `json:"content,omitempty"`
	Thinking string `json:"thinking,omitempty"`
}

type Choice struct {
	Index int   `json:"index"`
	Delta Delta `json:"delta"`
}

type SSEEvent struct {
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Message struct {
	Role     string `json:"role"`
	Content  string `json:"content"`
	Thinking string `json:"thinking,omitempty"`
}

type ResponseChoice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []ResponseChoice `json:"choices"`
	Usage   *Usage           `json:"usage,omitempty"`
}

func isTerminal() bool {
	fi, _ := os.Stdout.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func parseArgs() (method string, headers []string, data string, url string) {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-X" || arg == "--request":
			if i+1 < len(args) {
				method = args[i+1]
				i++
			}
		case arg == "-H" || arg == "--header":
			if i+1 < len(args) {
				headers = append(headers, args[i+1])
				i++
			}
		case arg == "-d" || arg == "--data" || arg == "--data-binary":
			if i+1 < len(args) {
				data = args[i+1]
				i++
			}
		case !strings.HasPrefix(arg, "-"):
			url = arg
		}
	}
	return
}

func buildBody(data string) io.Reader {
	if data == "" {
		return nil
	}
	if strings.HasPrefix(data, "@") {
		name := data[1:]
		if name == "-" {
			return os.Stdin
		}
		b, err := os.ReadFile(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file '%s': %v\n", name, err)
			os.Exit(1)
		}
		return bytes.NewReader(b)
	}
	return strings.NewReader(data)
}

func renderSSE(r io.Reader) {
	useANSI := isTerminal()
	var dim, reset string
	if useANSI {
		dim = "\033[2m"
		reset = "\033[0m"
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	inThinking := false
	hasOutput := false
	out := os.Stdout

	startTime := time.Now()
	firstTokenTime := time.Time{}
	tokenCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimSpace(line[6:])
		if data == "[DONE]" {
			continue
		}

		var event SSEEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		if len(event.Choices) == 0 {
			if event.Usage != nil {
				fmt.Fprintf(os.Stderr, "\nusage: %d prompt + %d completion = %d total tokens\n",
					event.Usage.PromptTokens, event.Usage.CompletionTokens, event.Usage.TotalTokens)
			}
			continue
		}

		delta := event.Choices[0].Delta

		if delta.Thinking != "" {
			if firstTokenTime.IsZero() {
				firstTokenTime = time.Now()
			}
			tokenCount++

			if !inThinking && hasOutput {
				out.WriteString(reset)
				out.WriteString("\n\n")
			}
			if !inThinking {
				out.WriteString(dim)
				inThinking = true
			}
			out.WriteString(delta.Thinking)
			out.Sync()
			hasOutput = true
		}

		if delta.Content != "" {
			if firstTokenTime.IsZero() {
				firstTokenTime = time.Now()
			}
			tokenCount++

			if inThinking {
				out.WriteString(reset)
				out.WriteString("\n\n")
				inThinking = false
			}
			out.WriteString(delta.Content)
			out.Sync()
			hasOutput = true
		}
	}

	if hasOutput {
		out.WriteString(reset)
		out.WriteString("\n")
	}

	if !firstTokenTime.IsZero() {
		elapsed := time.Since(startTime).Seconds()
		avgTokenLatency := elapsed / float64(tokenCount)
		fmt.Fprintf(os.Stderr, "\r首token延迟: %.2fs | 平均token延迟: %.2fs\n", time.Since(firstTokenTime).Seconds(), avgTokenLatency)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading stream: %v\n", err)
	}
}

func renderJSON(r io.Reader) {
	var resp ChatResponse
	body, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading response: %v\n", err)
		return
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		os.Stdout.Write(body)
		return
	}
	if len(resp.Choices) == 0 {
		os.Stdout.Write(body)
		return
	}

	useANSI := isTerminal()
	msg := resp.Choices[0].Message

	if msg.Thinking != "" {
		if useANSI {
			fmt.Printf("\033[2m%s\033[0m\n\n", msg.Thinking)
		} else {
			fmt.Printf("%s\n\n", msg.Thinking)
		}
	}
	fmt.Print(msg.Content)
	fmt.Println()

	if resp.Usage != nil {
		fmt.Fprintf(os.Stderr, "usage: %d prompt + %d completion = %d total tokens\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}
}

func main() {
	method, headers, data, url := parseArgs()
	if url == "" {
		fmt.Fprintln(os.Stderr, "error: no URL specified")
		fmt.Fprintln(os.Stderr, "usage: goai [-X METHOD] [-H HEADER]... [-d DATA|@FILE] URL")
		os.Exit(1)
	}

	if method == "" {
		if data != "" {
			method = "POST"
		} else {
			method = "GET"
		}
	}

	body := buildBody(data)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating request: %v\n", err)
		os.Exit(1)
	}

	for _, h := range headers {
		if k, v, ok := strings.Cut(h, ":"); ok {
			req.Header.Set(strings.TrimSpace(k), strings.TrimSpace(v))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		os.Exit(1)
	}

	ctype := resp.Header.Get("Content-Type")
	isSSE := strings.Contains(ctype, "text/event-stream")
	isChat := strings.Contains(url, "/chat/completions")

	if isChat && isSSE {
		renderSSE(resp.Body)
	} else if isChat {
		renderJSON(resp.Body)
	} else {
		io.Copy(os.Stdout, resp.Body)
	}
}
