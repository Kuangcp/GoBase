package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0")).
			Padding(0, 1)

	collectionStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	requestStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	responseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00D68F")).
			Padding(1)

	methodStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	sendButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#00D68F")).
			Padding(0, 2).
			Bold(true)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Underline(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A0A0A0"))

	// Modal styles
	modalOverlayStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#000000"))

	modalBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(50)

	modalTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	modalHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))

	// Input modes
	modeNormal        = 0
	modeInput         = 1
	modeConfirmDelete = 2
)

// RequestData stores all data for a single request
type RequestData struct {
	Name    string
	URL     string
	Method  string
	Headers string
	Body    string
	Params  string
	Auth    string

	Response     string
	ResponseTime int
	StatusCode   int
	IsLoading    bool
	StartTime    time.Time
}

// Model represents the application state
type Model struct {
	width  int
	height int

	// Tree structure
	collections []Collection
	treeItems   []TreeItem
	expanded    map[string]bool
	selected    int

	// Active request (pointer into collections)
	activeReq *RequestData

	// Focus: 0=tree, 1=url, 2=tab_content
	focus int

	// Active tab in request area
	activeTab int // 0:Body, 1:Headers, 2:Params, 3:Auth

	// Text inputs for editing
	urlInput    textinput.Model
	bodyInput   textinput.Model
	headerInput textinput.Model
	paramInput  textinput.Model
	authInput   textinput.Model

	// Response viewport
	viewport viewport.Model

	// Modal input state
	inputMode   int
	inputBuf    textinput.Model
	inputTitle  string
	pendingPath string // path of item being operated on
}

func NewModel() Model {
	urlInput := textinput.New()
	urlInput.Placeholder = "Enter URL..."
	urlInput.CharLimit = 256
	urlInput.Width = 60

	bodyInput := textinput.New()
	bodyInput.Placeholder = "Request body..."
	bodyInput.CharLimit = 4096
	bodyInput.Width = 80

	headerInput := textinput.New()
	headerInput.Placeholder = "Key: Value"
	headerInput.CharLimit = 1024
	headerInput.Width = 80

	paramInput := textinput.New()
	paramInput.Placeholder = "key=value"
	paramInput.CharLimit = 1024
	paramInput.Width = 80

	authInput := textinput.New()
	authInput.Placeholder = "Bearer token..."
	authInput.CharLimit = 1024
	authInput.Width = 80

	vp := viewport.New(0, 0)

	inputBuf := textinput.New()
	inputBuf.CharLimit = 64
	inputBuf.Width = 40

	expanded := map[string]bool{}

	// Load from disk
	collections := loadAll()

	// Expand all collections by default
	for i := range collections {
		expanded[collections[i].Name] = true
	}

	m := Model{
		collections: collections,
		expanded:    expanded,
		selected:    0,
		focus:       0,
		activeTab:   0,
		urlInput:    urlInput,
		bodyInput:   bodyInput,
		headerInput: headerInput,
		paramInput:  paramInput,
		authInput:   authInput,
		viewport:    vp,
		inputBuf:    inputBuf,
		inputMode:   modeNormal,
	}

	// Build tree items and load first request
	m.treeItems = flattenTree(m.collections, m.expanded)
	if len(m.treeItems) > 0 {
		// Find first request item
		for i, item := range m.treeItems {
			if item.Request != nil {
				m.selected = i
				m.activeReq = item.Request
				m.loadActiveRequest()
				break
			}
		}
	}

	return m
}

func (m *Model) loadActiveRequest() {
	if m.activeReq == nil {
		return
	}
	req := m.activeReq

	m.urlInput.SetValue(req.URL)
	m.bodyInput.SetValue(req.Body)
	m.headerInput.SetValue(req.Headers)
	m.paramInput.SetValue(req.Params)
	m.authInput.SetValue(req.Auth)

	m.viewport.SetContent(req.Response)
}

func (m *Model) saveActiveRequest() {
	if m.activeReq == nil {
		return
	}
	req := m.activeReq
	req.URL = m.urlInput.Value()
	req.Body = m.bodyInput.Value()
	req.Headers = m.headerInput.Value()
	req.Params = m.paramInput.Value()
	req.Auth = m.authInput.Value()
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle modal input mode
		if m.inputMode == modeInput {
			switch msg.String() {
			case "enter":
				name := strings.TrimSpace(m.inputBuf.Value())
				if name != "" {
					m.applyInput(name)
				}
				m.inputMode = modeNormal
				m.inputBuf.Blur()
				m.refreshTree()
				return m, nil
			case "esc":
				m.inputMode = modeNormal
				m.inputBuf.Blur()
				return m, nil
			default:
				var cmd tea.Cmd
				m.inputBuf, cmd = m.inputBuf.Update(msg)
				return m, cmd
			}
		}

		// Handle confirm delete mode
		if m.inputMode == modeConfirmDelete {
			switch msg.String() {
			case "y", "Y":
				m.deleteItem()
				m.inputMode = modeNormal
				m.refreshTree()
				return m, nil
			case "n", "N", "esc":
				m.inputMode = modeNormal
				return m, nil
			}
			return m, nil
		}

		// Normal mode
		switch msg.String() {
		case "ctrl+c":
			m.saveActiveRequest()
			saveAll(m.collections)
			return m, tea.Quit

		case "tab":
			m.focus = (m.focus + 1) % 3
			m.updateFocus()

		case "shift+tab":
			m.focus = (m.focus + 2) % 3
			m.updateFocus()

		case "ctrl+s":
			if m.focus >= 1 && m.activeReq != nil && !m.activeReq.IsLoading {
				m.saveActiveRequest()
				m.activeReq.IsLoading = true
				m.activeReq.StartTime = time.Now()
				return m, tea.Batch(m.sendCurrentRequest(), m.startTimer())
			}

		case "up":
			if m.focus == 0 {
				m.moveUp()
			}

		case "down":
			if m.focus == 0 {
				m.moveDown()
			}

		case "right":
			if m.focus == 0 {
				m.toggleExpand(true)
			}

		case "left":
			if m.focus == 0 {
				m.toggleExpand(false)
			}

		case "enter":
			if m.focus == 0 {
				m.toggleExpand(!m.isExpanded())
			}

		case "n":
			if m.focus == 0 {
				m.startCreateRequest()
			}

		case "N":
			if m.focus == 0 {
				m.startCreateFolder()
			}

		case "r":
			if m.focus == 0 {
				m.startRename()
			}

		case "x":
			if m.focus == 0 {
				m.startDelete()
			}

		case "alt+m":
			if m.activeReq != nil {
				methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
				current := m.activeReq.Method
				for i, method := range methods {
					if method == current {
						m.activeReq.Method = methods[(i+1)%len(methods)]
						break
					}
				}
			}
		}

	case tea.MouseMsg:
		if m.inputMode != modeNormal {
			return m, nil
		}
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			x, y := msg.X, msg.Y
			m.handleClick(x, y)
		}
		if msg.Action == tea.MouseActionRelease {
			if msg.Button == tea.MouseButtonWheelUp {
				m.viewport.YOffset--
			} else if msg.Button == tea.MouseButtonWheelDown {
				m.viewport.YOffset++
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width/2 - 4
		m.viewport.Height = msg.Height/2 - 8

	case httpResponseMsg:
		if m.activeReq != nil {
			m.activeReq.IsLoading = false
			if msg.err != nil {
				m.activeReq.Response = fmt.Sprintf("Error: %v", msg.err)
				m.activeReq.StatusCode = 0
			} else {
				m.activeReq.StatusCode = msg.StatusCode
				m.activeReq.Response = msg.Body
				m.activeReq.ResponseTime = int(msg.Duration.Milliseconds())
				m.viewport.SetContent(m.activeReq.Response)
			}
		}

	case loadingTickMsg:
		if m.activeReq != nil && m.activeReq.IsLoading {
			cmds = append(cmds, m.startTimer())
		}
	}

	// Update focused input (only when not in modal mode)
	if m.inputMode == modeNormal {
		var cmd tea.Cmd
		switch m.focus {
		case 1:
			m.urlInput, cmd = m.urlInput.Update(msg)
		case 2:
			switch m.activeTab {
			case 0:
				m.bodyInput, cmd = m.bodyInput.Update(msg)
			case 1:
				m.headerInput, cmd = m.headerInput.Update(msg)
			case 2:
				m.paramInput, cmd = m.paramInput.Update(msg)
			case 3:
				m.authInput, cmd = m.authInput.Update(msg)
			}
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateFocus() {
	m.urlInput.Blur()
	m.bodyInput.Blur()
	m.headerInput.Blur()
	m.paramInput.Blur()
	m.authInput.Blur()

	switch m.focus {
	case 1:
		m.urlInput.Focus()
	case 2:
		switch m.activeTab {
		case 0:
			m.bodyInput.Focus()
		case 1:
			m.headerInput.Focus()
		case 2:
			m.paramInput.Focus()
		case 3:
			m.authInput.Focus()
		}
	}
}

func (m *Model) switchTab(tab int) {
	m.saveActiveRequest()
	m.activeTab = tab
	m.focus = 2
	m.updateFocus()
}

// Tree navigation

func (m *Model) moveUp() {
	if m.selected > 0 {
		m.saveActiveRequest()
		m.selected--
		if m.treeItems[m.selected].Request != nil {
			m.activeReq = m.treeItems[m.selected].Request
			m.loadActiveRequest()
		}
	}
}

func (m *Model) moveDown() {
	if m.selected < len(m.treeItems)-1 {
		m.saveActiveRequest()
		m.selected++
		if m.treeItems[m.selected].Request != nil {
			m.activeReq = m.treeItems[m.selected].Request
			m.loadActiveRequest()
		}
	}
}

func (m *Model) toggleExpand(expand bool) {
	if m.selected >= len(m.treeItems) {
		return
	}
	item := &m.treeItems[m.selected]
	if item.Request != nil {
		return // Can't expand a request
	}
	m.expanded[item.Path] = expand
	m.treeItems = flattenTree(m.collections, m.expanded)
	// Keep selection in bounds
	if m.selected >= len(m.treeItems) {
		m.selected = len(m.treeItems) - 1
	}
}

func (m *Model) isExpanded() bool {
	if m.selected >= len(m.treeItems) {
		return false
	}
	item := &m.treeItems[m.selected]
	if item.Request != nil {
		return false
	}
	return m.expanded[item.Path]
}

// Modal input handlers

func (m *Model) startCreateRequest() {
	if m.selected >= len(m.treeItems) {
		return
	}
	item := &m.treeItems[m.selected]
	m.inputMode = modeInput
	m.inputTitle = "New Request"
	m.inputBuf.SetValue("New Request")
	m.inputBuf.Focus()
	m.pendingPath = item.Path
}

func (m *Model) startCreateFolder() {
	if m.selected >= len(m.treeItems) {
		return
	}
	item := &m.treeItems[m.selected]
	m.inputMode = modeInput
	m.inputTitle = "New Folder"
	m.inputBuf.SetValue("New Folder")
	m.inputBuf.Focus()
	m.pendingPath = item.Path
}

func (m *Model) startRename() {
	if m.selected >= len(m.treeItems) {
		return
	}
	item := &m.treeItems[m.selected]
	m.inputMode = modeInput
	m.inputTitle = "Rename"
	m.inputBuf.SetValue(item.Name)
	m.inputBuf.Focus()
	m.pendingPath = item.Path
}

func (m *Model) startDelete() {
	if m.selected >= len(m.treeItems) {
		return
	}
	item := &m.treeItems[m.selected]
	m.inputMode = modeConfirmDelete
	m.inputTitle = fmt.Sprintf("Delete \"%s\" ?", item.Name)
	m.pendingPath = item.Path
}

func (m *Model) applyInput(name string) {
	parts := strings.SplitN(m.pendingPath, "/", 2)
	if len(parts) == 0 {
		return
	}

	collName := parts[0]

	// Find collection index
	collIdx := -1
	for i := range m.collections {
		if m.collections[i].Name == collName {
			collIdx = i
			break
		}
	}
	if collIdx < 0 {
		return
	}

	switch m.inputTitle {
	case "New Request":
		req := RequestData{
			Name:     name,
			Method:   "GET",
			Response: "Click Send to load response...",
		}

		if len(parts) == 1 {
			// Add to collection root as a folder with one request
			folder := Folder{
				Name:     name,
				Requests: []RequestData{req},
			}
			m.collections[collIdx].Folders = append(m.collections[collIdx].Folders, folder)
			m.expanded[collName] = true
		} else {
			// Add to folder
			folderName := parts[1]
			for fi := range m.collections[collIdx].Folders {
				if m.collections[collIdx].Folders[fi].Name == folderName {
					m.collections[collIdx].Folders[fi].Requests = append(
						m.collections[collIdx].Folders[fi].Requests, req,
					)
					m.expanded[collName+"/"+folderName] = true
					break
				}
			}
		}

	case "New Folder":
		folder := Folder{Name: name}
		if len(parts) == 1 {
			// Add folder to collection
			m.collections[collIdx].Folders = append(m.collections[collIdx].Folders, folder)
			m.expanded[collName] = true
		}

	case "Rename":
		if len(parts) == 1 {
			// Rename collection
			oldPath := m.collections[collIdx].Name
			m.collections[collIdx].Name = name
			// Update expanded map
			delete(m.expanded, oldPath)
			m.expanded[name] = true
		} else {
			folderName := parts[1]
			for fi := range m.collections[collIdx].Folders {
				if m.collections[collIdx].Folders[fi].Name == folderName {
					oldPath := collName + "/" + folderName
					m.collections[collIdx].Folders[fi].Name = name
					delete(m.expanded, oldPath)
					m.expanded[collName+"/"+name] = true

					// Also rename request file if it's a single-request folder
					if len(m.collections[collIdx].Folders[fi].Requests) == 1 {
						m.collections[collIdx].Folders[fi].Requests[0].Name = name
					}
					break
				}
			}
		}
	}
}

func (m *Model) deleteItem() {
	parts := strings.SplitN(m.pendingPath, "/", 2)
	if len(parts) == 0 {
		return
	}

	collName := parts[0]

	// Find collection index
	collIdx := -1
	for i := range m.collections {
		if m.collections[i].Name == collName {
			collIdx = i
			break
		}
	}
	if collIdx < 0 {
		return
	}

	if len(parts) == 1 {
		// Delete entire collection
		dir := filepath.Join(dataDir(), collName)
		os.RemoveAll(dir)
		m.collections = append(m.collections[:collIdx], m.collections[collIdx+1:]...)
	} else {
		folderName := parts[1]
		for fi := range m.collections[collIdx].Folders {
			if m.collections[collIdx].Folders[fi].Name == folderName {
				// Delete folder from disk
				dir := filepath.Join(dataDir(), collName, folderName)
				os.RemoveAll(dir)
				// Delete from memory
				folders := m.collections[collIdx].Folders
				m.collections[collIdx].Folders = append(folders[:fi], folders[fi+1:]...)
				break
			}
		}
	}

	// Clear active request if it was deleted
	if m.activeReq != nil {
		found := false
		for _, item := range m.treeItems {
			if item.Request == m.activeReq {
				found = true
				break
			}
		}
		if !found {
			m.activeReq = nil
		}
	}
}

func (m *Model) refreshTree() {
	m.treeItems = flattenTree(m.collections, m.expanded)
	if m.selected >= len(m.treeItems) {
		m.selected = len(m.treeItems) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
	// Try to load request at current selection
	if m.selected < len(m.treeItems) && m.treeItems[m.selected].Request != nil {
		m.activeReq = m.treeItems[m.selected].Request
		m.loadActiveRequest()
	}
}

func (m *Model) handleClick(x, y int) {
	statusHeight := 1
	collectionWidth := m.width / 4

	// Click in tree area
	if x < collectionWidth && y > statusHeight {
		itemIndex := y - statusHeight - 2
		if itemIndex >= 0 && itemIndex < len(m.treeItems) {
			m.saveActiveRequest()
			m.selected = itemIndex
			item := m.treeItems[m.selected]
			if item.Request != nil {
				m.activeReq = item.Request
				m.loadActiveRequest()
			} else {
				// Toggle expand/collapse on click
				m.toggleExpand(!m.expanded[item.Path])
			}
			m.focus = 0
		}
		return
	}

	// Click in right panel
	if x >= collectionWidth {
		rightX := x - collectionWidth

		// Tab bar is at y offset 6 (border + padding + url bar + blank + tabs)
		tabY := statusHeight + 5
		if y == tabY {
			// Calculate tab click
			tabStart := rightX - 3
			if tabStart >= 0 {
				tabWidths := []int{8, 12, 12, 8}
				offset := 0
				for i, tw := range tabWidths {
					if tabStart >= offset && tabStart < offset+tw {
						m.switchTab(i)
						return
					}
					offset += tw + 2
				}
			}
		}

		// Click in URL bar area
		if y == statusHeight+2 {
			m.focus = 1
			m.updateFocus()
			return
		}

		// Click in content area (below tabs)
		if y > tabY {
			m.focus = 2
			m.updateFocus()
			return
		}

		m.focus = 1
		m.updateFocus()
	}
}

func (m Model) sendCurrentRequest() tea.Cmd {
	return func() tea.Msg {
		req := m.activeReq
		if req == nil {
			return httpResponseMsg{err: fmt.Errorf("no request selected")}
		}
		start := time.Now()

		var body io.Reader
		if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
			body = strings.NewReader(req.Body)
		}

		httpReq, err := http.NewRequest(req.Method, req.URL, body)
		if err != nil {
			return httpResponseMsg{err: err}
		}

		// Parse headers
		for _, line := range strings.Split(req.Headers, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				httpReq.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}

		// Append query params
		if req.Params != "" {
			sep := "?"
			if strings.Contains(req.URL, "?") {
				sep = "&"
			}
			httpReq.URL.RawQuery += sep + req.Params
		}

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			return httpResponseMsg{err: err}
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return httpResponseMsg{err: err}
		}

		return httpResponseMsg{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
			Duration:   time.Since(start),
		}
	}
}

type httpResponseMsg struct {
	StatusCode int
	Body       string
	Duration   time.Duration
	err        error
}

type loadingTickMsg struct{}

func (m Model) startTimer() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return loadingTickMsg{}
	})
}

// --- View ---

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	statusHeight := 1
	collectionWidth := m.width / 4
	requestWidth := m.width - collectionWidth

	statusBar := m.renderStatusBar()
	treePanel := m.renderTree(collectionWidth, m.height-statusHeight)
	rightPanel := m.renderRight(requestWidth, m.height-statusHeight)

	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		lipgloss.JoinHorizontal(lipgloss.Top, treePanel, rightPanel),
	)

	// Overlay modal if in input mode
	if m.inputMode == modeInput {
		return m.renderInputModal(mainView)
	}
	if m.inputMode == modeConfirmDelete {
		return m.renderConfirmModal(mainView)
	}

	return mainView
}

func (m Model) renderStatusBar() string {
	title := titleStyle.Render("API-Y - HTTP Client")

	focusLabel := ""
	switch m.focus {
	case 0:
		focusLabel = statusBarStyle.Render(" [Tree]")
	case 1:
		focusLabel = statusBarStyle.Render(" [URL]")
	case 2:
		focusLabel = statusBarStyle.Render(" [Request]")
	}

	hints := statusBarStyle.Render("  Tab:switch  Ctrl+S:send  n:new  r:rename  x:delete")
	used := lipgloss.Width(title) + lipgloss.Width(focusLabel)
	padding := m.width - used - lipgloss.Width(hints)
	if padding > 0 {
		hints = statusBarStyle.Render(strings.Repeat(" ", padding)) + hints
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, title, focusLabel, hints)
}

func (m Model) renderTree(width, height int) string {
	content := ""
	for i, item := range m.treeItems {
		prefix := "  "
		style := lipgloss.NewStyle()

		if i == m.selected {
			prefix = "> "
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
		}

		indent := strings.Repeat("  ", item.Level)

		switch {
		case item.Request != nil:
			// Request node
			methodColor := lipgloss.Color("#00D68F")
			switch item.Request.Method {
			case "POST":
				methodColor = lipgloss.Color("#FFD93D")
			case "PUT":
				methodColor = lipgloss.Color("#FFA500")
			case "DELETE":
				methodColor = lipgloss.Color("#FF6B6B")
			}
			methodTag := lipgloss.NewStyle().Foreground(methodColor).Bold(true).Render(item.Request.Method)
			content += style.Render(prefix+indent+methodTag+" "+item.Name) + "\n"

		default:
			// Collection or folder node
			icon := "  "
			if item.Expanded {
				icon = "▼ "
			} else {
				icon = "▶ "
			}
			nameStyle := lipgloss.NewStyle().Bold(true)
			if item.Level == 0 {
				nameStyle = nameStyle.Foreground(lipgloss.Color("#7D56F4"))
			}
			content += style.Render(prefix+indent+icon+nameStyle.Render(item.Name)) + "\n"
		}
	}

	style := collectionStyle
	if m.focus == 0 {
		style = focusedStyle
	}

	return style.Width(width - 2).Height(height - 2).Render(content)
}

func (m Model) renderRight(width, height int) string {
	requestHeight := height / 2
	responseHeight := height - requestHeight

	reqPanel := m.renderRequest(width, requestHeight)
	respPanel := m.renderResponse(width, responseHeight)

	return lipgloss.JoinVertical(lipgloss.Left, reqPanel, respPanel)
}

func (m Model) renderRequest(width, height int) string {
	req := m.activeReq
	if req == nil {
		return requestStyle.Width(width - 2).Height(height - 2).Render("No request selected")
	}

	methodColor := lipgloss.Color("#7D56F4")
	switch req.Method {
	case "POST":
		methodColor = lipgloss.Color("#FFD93D")
	case "PUT":
		methodColor = lipgloss.Color("#FFA500")
	case "DELETE":
		methodColor = lipgloss.Color("#FF6B6B")
	}

	methodTag := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(methodColor).
		Padding(0, 1).
		Render(req.Method + " ")

	urlLine := lipgloss.JoinHorizontal(lipgloss.Top, methodTag, " ", m.urlInput.View())

	sendBtn := sendButtonStyle.Render("[ Send ]")
	if req.IsLoading {
		sendBtn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#FFA500")).
			Padding(0, 2).
			Render(" Loading...")
	}

	// Tabs
	tabs := []string{"Body", "Headers", "Params", "Auth"}
	tabBar := " "
	for i, tab := range tabs {
		if i == m.activeTab {
			tabBar += activeTabStyle.Render(fmt.Sprintf("[%s]", tab)) + " "
		} else {
			tabBar += inactiveTabStyle.Render(fmt.Sprintf(" %s ", tab)) + " "
		}
	}

	// Active tab content
	var content string
	switch m.activeTab {
	case 0:
		content = m.bodyInput.View()
	case 1:
		content = m.headerInput.View()
	case 2:
		content = m.paramInput.View()
	case 3:
		content = m.authInput.View()
	}

	style := requestStyle
	if m.focus >= 1 {
		style = focusedStyle
	}

	return style.Width(width - 2).Height(height - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, urlLine, " ", sendBtn),
			"",
			tabBar,
			"",
			content,
		),
	)
}

func (m Model) renderResponse(width, height int) string {
	req := m.activeReq

	statusColor := lipgloss.Color("#A0A0A0")
	statusText := "Ready"

	if req != nil && req.IsLoading {
		elapsed := time.Since(req.StartTime).Milliseconds()
		statusColor = lipgloss.Color("#FFA500")
		statusText = fmt.Sprintf("Loading... %dms", elapsed)
		m.viewport.SetContent("Request in progress...")
	} else if req != nil && req.StatusCode > 0 {
		if req.StatusCode >= 200 && req.StatusCode < 300 {
			statusColor = lipgloss.Color("#00D68F")
		} else if req.StatusCode >= 300 && req.StatusCode < 400 {
			statusColor = lipgloss.Color("#FFD93D")
		} else if req.StatusCode >= 400 {
			statusColor = lipgloss.Color("#FF6B6B")
		}
		statusText = fmt.Sprintf("Status: %d | Time: %dms", req.StatusCode, req.ResponseTime)
	}

	statusBar := lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render(statusText)

	return responseStyle.Width(width - 2).Height(height - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			statusBar,
			"",
			m.viewport.View(),
		),
	)
}

func (m Model) renderInputModal(baseView string) string {
	modalContent := lipgloss.JoinVertical(lipgloss.Left,
		modalTitleStyle.Render(m.inputTitle),
		"",
		m.inputBuf.View(),
		"",
		modalHintStyle.Render("Enter: confirm    Esc: cancel"),
	)

	modal := modalBoxStyle.Render(modalContent)
	return centeredOverlay(baseView, modal, m.width, m.height)
}

func (m Model) renderConfirmModal(baseView string) string {
	modalContent := lipgloss.JoinVertical(lipgloss.Left,
		modalTitleStyle.Render(m.inputTitle),
		"",
		modalHintStyle.Render("y: confirm    n/Esc: cancel"),
	)

	modal := modalBoxStyle.Render(modalContent)
	return centeredOverlay(baseView, modal, m.width, m.height)
}

func centeredOverlay(base, overlay string, width, height int) string {
	overlayWidth := lipgloss.Width(overlay)
	overlayHeight := lipgloss.Height(overlay)

	x := (width - overlayWidth) / 2
	y := (height - overlayHeight) / 2

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	// Build overlay by placing content at the right position
	lines := strings.Split(base, "\n")
	for i := y; i < height && (i-y) < overlayHeight; i++ {
		if i >= len(lines) {
			break
		}
		overlayLine := ""
		olLines := strings.Split(overlay, "\n")
		if (i-y) < len(olLines) {
			overlayLine = olLines[i-y]
		}

		lineRunes := []rune(lines[i])
		overlayRunes := []rune(overlayLine)

		start := x
		for j, r := range overlayRunes {
			if start+j < len(lineRunes) {
				lineRunes[start+j] = r
			} else {
				lineRunes = append(lineRunes, r)
			}
		}
		lines[i] = string(lineRunes)
	}

	return strings.Join(lines, "\n")
}
