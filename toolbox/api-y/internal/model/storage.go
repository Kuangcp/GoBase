package model

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Collection represents a top-level group of folders and requests
type Collection struct {
	Name    string
	Folders []Folder
}

// Folder represents a group of requests within a collection
type Folder struct {
	Name     string
	Requests []RequestData
}

// TreeItem represents a flattened node in the UI tree
type TreeItem struct {
	Level    int // 0=collection, 1=folder, 2=request
	Name     string
	Expanded bool
	Request  *RequestData // nil for collections/folders
	Path     string       // filesystem path
}

func dataDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "api-y")
}

// loadAll loads all collections from the data directory
func loadAll() []Collection {
	dir := dataDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return defaultCollections()
	}

	var collections []Collection
	entries, err := os.ReadDir(dir)
	if err != nil {
		return defaultCollections()
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		collection := loadCollection(filepath.Join(dir, entry.Name()), entry.Name())
		collections = append(collections, collection)
	}

	if len(collections) == 0 {
		return defaultCollections()
	}

	return collections
}

func loadCollection(path, name string) Collection {
	collection := Collection{Name: name}

	entries, err := os.ReadDir(path)
	if err != nil {
		return collection
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folder := loadFolder(filepath.Join(path, entry.Name()), entry.Name())
			collection.Folders = append(collection.Folders, folder)
		} else if strings.HasSuffix(entry.Name(), ".http") || strings.HasSuffix(entry.Name(), ".rest") {
			// Single .http file at collection root (no folder)
			requests := parseHTTPFile(filepath.Join(path, entry.Name()))
			if len(requests) > 0 {
				folder := Folder{
					Name:     strings.TrimSuffix(strings.TrimSuffix(entry.Name(), ".http"), ".rest"),
					Requests: requests,
				}
				collection.Folders = append(collection.Folders, folder)
			}
		}
	}

	return collection
}

func loadFolder(path, name string) Folder {
	folder := Folder{Name: name}

	entries, err := os.ReadDir(path)
	if err != nil {
		return folder
	}

	// Sort: folders first, then files
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			// Subfolder: load its .http files as requests in this folder
			subPath := filepath.Join(path, entry.Name())
			subEntries, err := os.ReadDir(subPath)
			if err != nil {
				continue
			}
			for _, subEntry := range subEntries {
				if !subEntry.IsDir() && (strings.HasSuffix(subEntry.Name(), ".http") || strings.HasSuffix(subEntry.Name(), ".rest")) {
					requests := parseHTTPFile(filepath.Join(subPath, subEntry.Name()))
					folder.Requests = append(folder.Requests, requests...)
				}
			}
		} else if strings.HasSuffix(entry.Name(), ".http") || strings.HasSuffix(entry.Name(), ".rest") {
			requests := parseHTTPFile(filepath.Join(path, entry.Name()))
			folder.Requests = append(folder.Requests, requests...)
		}
	}

	return folder
}

// parseHTTPFile parses an IntelliJ .http file into RequestData slices
func parseHTTPFile(path string) []RequestData {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var requests []RequestData
	var current *RequestData
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Separator: ### optionally followed by a name
		if strings.HasPrefix(line, "###") {
			name := strings.TrimSpace(strings.TrimPrefix(line, "###"))
			req := RequestData{
				Name:     name,
				Response: "Click Send to load response...",
			}
			requests = append(requests, req)
			current = &requests[len(requests)-1]
			continue
		}

		if current == nil {
			continue
		}

		// Parse method line: METHOD URL
		if current.Method == "" {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			parts := strings.SplitN(trimmed, " ", 2)
			if len(parts) == 2 {
				current.Method = strings.ToUpper(parts[0])
				rawURL := parts[1]

				// Extract query params from URL
				if idx := strings.Index(rawURL, "?"); idx >= 0 {
					current.URL = rawURL[:idx]
					current.Params = rawURL[idx+1:]
				} else {
					current.URL = rawURL
				}

				if current.Name == "" {
					current.Name = fmt.Sprintf("%s %s", current.Method, current.URL)
				}
			}
			continue
		}

		// Empty line separates headers from body
		if strings.TrimSpace(line) == "" && len(current.Headers) > 0 {
			// Read remaining lines as body
			var bodyLines []string
			for scanner.Scan() {
				bodyLines = append(bodyLines, scanner.Text())
			}
			current.Body = strings.Join(bodyLines, "\n")
			current.Body = strings.TrimSpace(current.Body)
			continue
		}

		// Parse header line
		if strings.TrimSpace(line) != "" {
			// Check if this is the Authorization header → store as Auth
			lower := strings.ToLower(strings.TrimSpace(line))
			if strings.HasPrefix(lower, "authorization:") {
				authValue := strings.TrimSpace(strings.TrimPrefix(line, "Authorization:"))
				authValue = strings.TrimSpace(strings.TrimPrefix(line, "authorization:"))
				current.Auth = authValue
			} else {
				if current.Headers != "" {
					current.Headers += "\n"
				}
				current.Headers += strings.TrimSpace(line)
			}
		}
	}

	return requests
}

// saveAll saves all collections to the data directory
func saveAll(collections []Collection) error {
	dir := dataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data dir: %w", err)
	}

	// Remove existing directories (simple approach: delete and recreate)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read data dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			os.RemoveAll(filepath.Join(dir, entry.Name()))
		}
	}

	// Write each collection
	for _, collection := range collections {
		collectionDir := filepath.Join(dir, collection.Name)
		if err := os.MkdirAll(collectionDir, 0755); err != nil {
			return fmt.Errorf("failed to create collection dir: %w", err)
		}

		for _, folder := range collection.Folders {
			folderDir := filepath.Join(collectionDir, folder.Name)
			if err := os.MkdirAll(folderDir, 0755); err != nil {
				return fmt.Errorf("failed to create folder dir: %w", err)
			}

			for i, req := range folder.Requests {
				filename := fmt.Sprintf("%s.http", sanitizeFilename(req.Name))
				if i > 0 {
					// Handle duplicate names
					filename = fmt.Sprintf("%s_%d.http", sanitizeFilename(req.Name), i)
				}
				reqPath := filepath.Join(folderDir, filename)
				if err := writeHTTPFile(reqPath, req); err != nil {
					return fmt.Errorf("failed to write request: %w", err)
				}
			}
		}
	}

	return nil
}

// writeHTTPFile writes a RequestData to an IntelliJ .http file
func writeHTTPFile(path string, req RequestData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	// Request name
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("%s %s", req.Method, req.URL)
	}
	fmt.Fprintf(w, "### %s\n", name)

	// Method + URL (append params)
	url := req.URL
	if req.Params != "" {
		if strings.Contains(url, "?") {
			url += "&" + req.Params
		} else {
			url += "?" + req.Params
		}
	}
	fmt.Fprintf(w, "%s %s\n", req.Method, url)

	// Headers
	if req.Headers != "" {
		for _, line := range strings.Split(req.Headers, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				fmt.Fprintf(w, "%s\n", line)
			}
		}
	}

	// Auth as Authorization header
	if req.Auth != "" {
		fmt.Fprintf(w, "Authorization: %s\n", req.Auth)
	}

	// Blank line before body
	fmt.Fprintf(w, "\n")

	// Body
	if req.Body != "" {
		fmt.Fprintf(w, "%s\n", req.Body)
	}

	w.Flush()
	return nil
}

func sanitizeFilename(name string) string {
	// Replace spaces and special chars
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "*", "")
	name = strings.ReplaceAll(name, "?", "")
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "<", "")
	name = strings.ReplaceAll(name, ">", "")
	name = strings.ReplaceAll(name, "|", "")
	if len(name) > 64 {
		name = name[:64]
	}
	return name
}

// flattenTree converts collections into a flat list of tree items for UI rendering
func flattenTree(collections []Collection, expanded map[string]bool) []TreeItem {
	var items []TreeItem

	for ci := range collections {
		c := &collections[ci]
		collPath := c.Name
		collExpanded := expanded[collPath]

		items = append(items, TreeItem{
			Level:    0,
			Name:     c.Name,
			Expanded: collExpanded,
			Path:     collPath,
		})

		if !collExpanded {
			continue
		}

		for fi := range c.Folders {
			f := &c.Folders[fi]
			folderPath := collPath + "/" + f.Name
			folderExpanded := expanded[folderPath]

			items = append(items, TreeItem{
				Level:    1,
				Name:     f.Name,
				Expanded: folderExpanded,
				Path:     folderPath,
			})

			if !folderExpanded {
				continue
			}

			for ri := range f.Requests {
				req := &f.Requests[ri]
				items = append(items, TreeItem{
					Level:   2,
					Name:    req.Name,
					Request: req,
					Path:    folderPath + "/" + req.Name,
				})
			}
		}
	}

	return items
}

// findRequest finds a RequestData by its tree path
func findRequest(collections []Collection, path string) *RequestData {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return nil
	}

	collName := parts[0]
	folderName := parts[1]
	reqName := strings.Join(parts[2:], "/")

	for ci := range collections {
		if collections[ci].Name != collName {
			continue
		}
		for fi := range collections[ci].Folders {
			if collections[ci].Folders[fi].Name != folderName {
				continue
			}
			for ri := range collections[ci].Folders[fi].Requests {
				if collections[ci].Folders[fi].Requests[ri].Name == reqName {
					return &collections[ci].Folders[fi].Requests[ri]
				}
			}
		}
	}
	return nil
}

func defaultCollections() []Collection {
	return []Collection{
		{
			Name: "Examples",
			Folders: []Folder{
				{
					Name: "HTTPBin",
					Requests: []RequestData{
						{
							Name:     "GET /users",
							URL:      "https://httpbin.org/get",
							Method:   "GET",
							Headers:  "Content-Type: application/json",
							Body:     "{}",
							Params:   "page=1&limit=10",
							Response: "Click Send to load response...",
						},
						{
							Name:     "POST /login",
							URL:      "https://httpbin.org/post",
							Method:   "POST",
							Headers:  "Content-Type: application/json",
							Body:     `{"username":"admin","password":"123456"}`,
							Response: "Click Send to load response...",
						},
						{
							Name:     "GET /health",
							URL:      "https://httpbin.org/status/200",
							Method:   "GET",
							Response: "Click Send to load response...",
						},
						{
							Name:     "PUT /update",
							URL:      "https://httpbin.org/put",
							Method:   "PUT",
							Headers:  "Content-Type: application/json",
							Body:     `{"id":1,"name":"updated"}`,
							Auth:     "Bearer abc123",
							Response: "Click Send to load response...",
						},
						{
							Name:     "DELETE /item",
							URL:      "https://httpbin.org/delete",
							Method:   "DELETE",
							Params:   "id=42",
							Response: "Click Send to load response...",
						},
					},
				},
			},
		},
	}
}
