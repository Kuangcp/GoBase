# Voice Drop — Design Spec

## Overview
Voice Drop is a PWA that captures speech on a mobile phone and types it into the current cursor position on a PC via X11's `xdotool`. The user presses a button to record, releases to see transcription, confirms/edits the text, then sends it to the Go backend over HTTP.

## Architecture

```
┌──────────────────────┐       HTTP POST /api/send       ┌──────────────────┐
│  Mobile Browser       │  ──────────────────────────►  │  Go Backend (PC) │
│  (PWA, Web Speech)    │  ◄──────────────────────────  │  (go:embed,      │
│                       │       200 OK / error          │   xdotool)       │
└──────────────────────┘                                └──────────────────┘
```

- **Single binary**: Go backend embeds all frontend files via `go:embed`.
- **No WebSocket**: HTTP POST is sufficient for one-shot text submission.
- **Keyboard simulation**: Backend shells out to `xdotool type --delay 0 <text>`.

## Component Breakdown

### Frontend (`frontend/`)
- `index.html` — single-page PWA with:
  - A press-and-hold record button (mic icon)
  - Text area showing recognition result (editable)
  - Send button → POST JSON to `/api/send`
  - Connection status indicator
- `manifest.json` — PWA manifest (name, icon, display: standalone)
- `sw.js` — Service Worker for offline cache (optional, v1 can be minimal)

### Backend (`main.go`)
- Embedded static file server for frontend assets
- `POST /api/send` endpoint: accepts `{"text":"..."}`, validates, executes `xdotool type`, responds JSON
- Configurable port via `-port` flag or `PORT` env var (default 8080)
- CORS headers for local development (mobile hits PC's IP)

## API

### `POST /api/send`
**Request:**
```json
{"text": "你好世界"}
```
**Response 200:**
```json
{"ok": true}
```
**Response 400/500:**
```json
{"ok": false, "error": "message"}
```

## Data Flow

1. User opens `http://<pc-ip>:8080` on phone
2. Holds record button → `SpeechRecognition.start()` captures audio
3. Releases button → `SpeechRecognition.stop()`, transcript appears in textarea
4. User edits transcript if needed, taps "Send"
5. `fetch POST /api/send {text}` → Go backend
6. Backend validates, runs `xdotool type --delay 0 <escaped text>`
7. Returns 200 OK; page shows success feedback

## Constraints & Assumptions

- PC runs Linux with X11 and `xdotool` installed
- Phone and PC are on the same LAN
- Browser supports Web Speech API (Chrome Android works well)
- No authentication (LAN-only, single-user)
- No HTTPS (dev-mode HTTP over LAN is acceptable; v2 can add TLS)

## Out of Scope (v1)

- No HTTPS/TLS
- No authentication
- No connection keepalive / heartbeat (HTTP request-response is stateless)
- No multiple-language support (default to system language)
- No history or saved recordings
