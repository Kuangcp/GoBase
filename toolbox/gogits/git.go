package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or git not available: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func getTotalFiles(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func getTotalLines(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	files := strings.Fields(string(out))
	var total int
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(repoPath, f))
		if err != nil {
			continue
		}
		total += bytes.Count(data, []byte{'\n'})
	}
	return total
}

func getExtensionStats(repoPath string) []ExtensionStat {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	files := strings.Fields(string(out))

	extFiles := make(map[string]int)
	extLines := make(map[string]int)
	for _, f := range files {
		ext := ""
		if idx := strings.LastIndex(f, "."); idx >= 0 {
			ext = f[idx+1:]
		}
		extFiles[ext]++
		data, err := os.ReadFile(filepath.Join(repoPath, f))
		if err != nil {
			continue
		}
		extLines[ext] += bytes.Count(data, []byte{'\n'})
	}

	var stats []ExtensionStat
	for ext, fc := range extFiles {
		stats = append(stats, ExtensionStat{
			Extension: ext,
			FileCount: fc,
			LineCount: extLines[ext],
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Extension == "" {
			return true
		}
		if stats[j].Extension == "" {
			return false
		}
		if stats[i].FileCount != stats[j].FileCount {
			return stats[i].FileCount > stats[j].FileCount
		}
		return stats[i].Extension < stats[j].Extension
	})

	return stats
}

func dayRange(start, end time.Time) ([]time.Time, []string) {
	if start.IsZero() || end.IsZero() || start.After(end) {
		return nil, nil
	}
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	var times []time.Time
	var labels []string
	for !s.After(e) {
		times = append(times, s)
		labels = append(labels, s.Format("2006-01-02"))
		s = s.AddDate(0, 0, 1)
	}
	return times, labels
}

func runGitLog(repoPath string) ([]CommitInfo, error) {
	cmd := exec.Command("git", "-C", repoPath, "log",
		"--format=COMMIT%n%H|%an|%ae|%ai",
		"--numstat", "HEAD")

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	defer cmd.Wait()

	scanner := bufio.NewScanner(out)
	var commits []CommitInfo
	state := 0
	var current *CommitInfo

	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case 0:
			if line == "COMMIT" {
				state = 1
			}
		case 1:
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "|", 4)
			if len(parts) >= 4 {
				t, err := time.Parse(gitTimeFormat, parts[3])
				if err != nil {
					t = time.Now()
				}
				current = &CommitInfo{
					Hash:   parts[0],
					Author: parts[1],
					Email:  parts[2],
					Date:   t,
				}
				state = 2
			}
		case 2:
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "COMMIT") {
				if current != nil {
					commits = append(commits, *current)
					current = nil
				}
				state = 1
			} else {
				parts := strings.SplitN(line, "\t", 3)
				if len(parts) >= 3 {
					added, err1 := strconv.Atoi(parts[0])
					deleted, err2 := strconv.Atoi(parts[1])
					if err1 == nil && err2 == nil {
						current.Added += added
						current.Deleted += deleted
					}
				}
			}
		}
	}

	if current != nil {
		commits = append(commits, *current)
	}

	return commits, scanner.Err()
}

func getDailyFileCounts(repoPath string) ([]string, []int) {
	cmd := exec.Command("git", "-C", repoPath, "log",
		"--reverse",
		"--format=COMMIT%n%H|%ai",
		"--diff-filter=AD",
		"--name-status",
		"HEAD")

	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	files := make(map[string]bool)
	type point struct {
		date  time.Time
		count int
	}
	var points []point
	var currentDate time.Time

	scanner := bufio.NewScanner(bytes.NewReader(out))
	state := 0

	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case 0:
			if line == "COMMIT" {
				state = 1
			}
		case 1:
			parts := strings.SplitN(line, "|", 2)
			if len(parts) == 2 {
				t, err := time.Parse(gitTimeFormat, parts[1])
				if err == nil {
					currentDate = t
				}
			}
			state = 2
		case 2:
			if strings.HasPrefix(line, "COMMIT") {
				points = append(points, point{currentDate, len(files)})
				state = 1
				continue
			}
			if len(line) > 1 && line[1] == '\t' {
				filename := line[2:]
				switch line[0] {
				case 'A':
					files[filename] = true
				case 'D':
					delete(files, filename)
				}
			}
		}
	}

	if !currentDate.IsZero() {
		points = append(points, point{currentDate, len(files)})
	}

	if len(points) == 0 {
		return nil, nil
	}

	for i := range points {
		points[i].date = time.Date(points[i].date.Year(), points[i].date.Month(), points[i].date.Day(), 0, 0, 0, 0, points[i].date.Location())
	}

	start := points[0].date
	end := points[len(points)-1].date
	allDays, allDayLabels := dayRange(start, end)

	data := make([]int, len(allDays))
	pi := 0
	for i, day := range allDays {
		for pi+1 < len(points) && !points[pi+1].date.After(day) {
			pi++
		}
		if pi < len(points) {
			data[i] = points[pi].count
		}
	}

	return allDayLabels, data
}

func getLargeFileCount(repoPath string, threshold int) int {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	files := strings.Fields(string(out))
	count := 0
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(repoPath, f))
		if err != nil {
			continue
		}
		if bytes.Count(data, []byte{'\n'}) > threshold {
			count++
		}
	}
	return count
}

func getTopLinesFiles(repoPath string, topN int) []FileHotspot {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	files := strings.Fields(string(out))
	type fileLine struct {
		path string
		line int
	}
	var list []fileLine
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(repoPath, f))
		if err != nil {
			continue
		}
		list = append(list, fileLine{path: f, line: bytes.Count(data, []byte{'\n'})})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].line > list[j].line
	})
	if len(list) > topN {
		list = list[:topN]
	}
	var result []FileHotspot
	for _, fl := range list {
		result = append(result, FileHotspot{Path: fl.path, ModifyCount: fl.line})
	}
	return result
}

func getTodoCount(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	files := strings.Fields(string(out))
	count := 0
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(repoPath, f))
		if err != nil {
			continue
		}
		content := string(data)
		for _, kw := range []string{"TODO", "FIXME", "HACK", "XXX"} {
			idx := 0
			for {
				pos := strings.Index(content[idx:], kw)
				if pos < 0 {
					break
				}
				count++
				idx += pos + len(kw)
			}
		}
	}
	return count
}

func getHotspots(repoPath string, topN int) []FileHotspot {
	cmd := exec.Command("git", "-C", repoPath, "log", "--since=120 days ago",
		"--pretty=format:", "--name-only", "--no-merges", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	freq := make(map[string]int)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		freq[line]++
	}
	type kv struct {
		path string
		cnt  int
	}
	var sorted []kv
	for k, v := range freq {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].cnt > sorted[j].cnt
	})
	if len(sorted) > topN {
		sorted = sorted[:topN]
	}
	result := make([]FileHotspot, len(sorted))
	for i, s := range sorted {
		result[i] = FileHotspot{Path: s.path, ModifyCount: s.cnt}
	}
	return result
}

func getAbandonedData(repoPath string, totalLOC int) (float64, int) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--since=1 year ago",
		"--pretty=format:", "--name-only", "--diff-filter=AM", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	recent := make(map[string]bool)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		recent[line] = true
	}
	cmd2 := exec.Command("git", "-C", repoPath, "ls-files")
	out2, err := cmd2.Output()
	if err != nil {
		return 0, 0
	}
	abandonedLOC := 0
	for _, f := range strings.Fields(string(out2)) {
		if recent[f] {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repoPath, f))
		if err != nil {
			continue
		}
		abandonedLOC += bytes.Count(data, []byte{'\n'})
	}
	if totalLOC <= 0 {
		return 0, abandonedLOC
	}
	return float64(abandonedLOC) / float64(totalLOC), abandonedLOC
}

func getCodeAgeDays(repoPath string) float64 {
	cmd := exec.Command("git", "-C", repoPath, "log", "--diff-filter=A",
		"--name-only", "--pretty=format:%ai|", "--no-merges", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	created := make(map[string]time.Time)
	var curDate time.Time
	for _, line := range strings.Split(string(out), "\n") {
		if idx := strings.IndexByte(line, '|'); idx >= 0 {
			t, err := time.Parse("2006-01-02 15:04:05 -0700", line[:idx])
			if err == nil {
				curDate = t
			}
		} else {
			line = strings.TrimSpace(line)
			if line != "" {
				if _, ok := created[line]; !ok {
					created[line] = curDate
				}
			}
		}
	}
	cmd2 := exec.Command("git", "-C", repoPath, "ls-files")
	out2, err := cmd2.Output()
	if err != nil {
		return 0
	}
	now := time.Now()
	totalDays := 0.0
	count := 0
	for _, f := range strings.Fields(string(out2)) {
		if cd, ok := created[f]; ok {
			totalDays += now.Sub(cd).Hours() / 24
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return totalDays / float64(count)
}

func getReleaseCount(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "tag", "--list")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func getOldCodeTouchRate(repoPath string) float64 {
	cmd := exec.Command("git", "-C", repoPath, "log", "--diff-filter=A",
		"--name-only", "--pretty=format:%ai|", "--no-merges", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	created := make(map[string]time.Time)
	var curDate time.Time
	for _, line := range strings.Split(string(out), "\n") {
		if idx := strings.IndexByte(line, '|'); idx >= 0 {
			t, err := time.Parse("2006-01-02 15:04:05 -0700", line[:idx])
			if err == nil {
				curDate = t
			}
		} else {
			line = strings.TrimSpace(line)
			if line != "" {
				if _, ok := created[line]; !ok {
					created[line] = curDate
				}
			}
		}
	}
	cmd2 := exec.Command("git", "-C", repoPath, "log", "--since=30 days ago",
		"--pretty=format:", "--name-only", "--no-merges", "HEAD")
	out2, err := cmd2.Output()
	if err != nil {
		return 0
	}
	cutoff := time.Now().AddDate(0, -3, 0)
	recentFiles := make(map[string]bool)
	for _, line := range strings.Split(string(out2), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		recentFiles[line] = true
	}
	oldCount := 0
	for f := range recentFiles {
		if cd, ok := created[f]; ok && cd.Before(cutoff) {
			oldCount++
		}
	}
	if len(recentFiles) == 0 {
		return 0
	}
	return float64(oldCount) / float64(len(recentFiles))
}

func isTestFile(path string) bool {
	lower := strings.ToLower(path)
	base := filepath.Base(lower)
	if strings.Contains(base, "_test") || strings.Contains(base, ".test.") ||
		strings.Contains(base, "_spec") || strings.Contains(base, ".spec.") ||
		strings.Contains(base, "_unittest") || strings.HasPrefix(base, "test_") {
		return true
	}
	dir := filepath.Dir(lower)
	for _, d := range []string{"/test/", "/tests/", "/__tests__/", "/spec/", "/specs/"} {
		if strings.Contains(dir, d) || dir == strings.Trim(d, "/") {
			return true
		}
	}
	return false
}

func getTestFileCount(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	files := strings.Fields(string(out))
	count := 0
	for _, f := range files {
		if isTestFile(f) {
			count++
		}
	}
	return count
}

func getAvgFilesPerCommit(repoPath string) float64 {
	cmd := exec.Command("git", "-C", repoPath, "log", "--no-merges",
		"--format=COMMIT%n", "--name-only", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(string(out), "\n")
	totalFiles := 0
	commitCount := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "COMMIT" {
			commitCount++
		} else if line != "" {
			totalFiles++
		}
	}
	if commitCount == 0 {
		return 0
	}
	return float64(totalFiles) / float64(commitCount)
}

func getRecentFileCount(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "log", "--since=120 days ago",
		"--pretty=format:", "--name-only", "--no-merges", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	files := make(map[string]bool)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		files[line] = true
	}
	return len(files)
}

func getDailyLocCounts(repoPath string, totalLoc int) ([]string, []int) {
	rootCmd := exec.Command("git", "-C", repoPath, "rev-parse", "--show-toplevel")
	rootOut, err := rootCmd.Output()
	if err != nil {
		return nil, nil
	}
	repoRoot := strings.TrimSpace(string(rootOut))
	prefix, err := filepath.Rel(repoRoot, repoPath)
	if err != nil {
		prefix = "."
	}
	args := []string{"-C", repoPath, "log",
		"--reverse",
		"--format=COMMIT%n%H|%ai",
		"--numstat",
		"HEAD"}
	if prefix != "." {
		args = append(args, "--", prefix+"/")
	}
	cmd := exec.Command("git", args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	type point struct {
		date time.Time
		loc  int
	}
	var points []point
	var currentDate time.Time
	fileLOC := make(map[string]int)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	state := 0

	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case 0:
			if line == "COMMIT" {
				state = 1
			}
		case 1:
			parts := strings.SplitN(line, "|", 2)
			if len(parts) == 2 {
				t, err := time.Parse(gitTimeFormat, parts[1])
				if err == nil {
					currentDate = t
				}
			}
			state = 2
		case 2:
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "COMMIT") {
				var total int
				for _, v := range fileLOC {
					total += v
				}
				points = append(points, point{currentDate, total})
				state = 1
				continue
			}
			parts := strings.SplitN(line, "\t", 3)
			if len(parts) >= 3 {
				added, err1 := strconv.Atoi(parts[0])
				deleted, err2 := strconv.Atoi(parts[1])
				if err1 == nil && err2 == nil {
					name := parts[2]
					if strings.Contains(name, " => ") {
						sub := strings.SplitN(name, " => ", 2)
						oldSize := fileLOC[sub[0]]
						delete(fileLOC, sub[0])
						name = sub[1]
						fileLOC[name] = oldSize
					}
					fileLOC[name] += added - deleted
					if fileLOC[name] <= 0 {
						delete(fileLOC, name)
					}
				}
			}
		}
	}

	if !currentDate.IsZero() {
		var total int
		for _, v := range fileLOC {
			total += v
		}
		points = append(points, point{currentDate, total})
	}

	if len(points) == 0 {
		return nil, nil
	}

	for i := range points {
		points[i].date = time.Date(points[i].date.Year(), points[i].date.Month(), points[i].date.Day(), 0, 0, 0, 0, points[i].date.Location())
	}

	start := points[0].date
	end := points[len(points)-1].date
	allDays, dayLabels := dayRange(start, end)

	data := make([]int, len(allDays))
	pi := 0
	for i, day := range allDays {
		for pi+1 < len(points) && !points[pi+1].date.After(day) {
			pi++
		}
		if pi < len(points) {
			data[i] = points[pi].loc
		}
	}

	if len(data) > 0 {
		adjust := totalLoc - data[len(data)-1]
		for i := range data {
			data[i] += adjust
		}
	}

	return dayLabels, data
}

func extractAuthorMonthlyStats(repoPath string) (map[string]map[string][]string, map[string]map[string]int, error) {
	cmd := exec.Command("git", "-C", repoPath, "log",
		"--all", "--no-merges",
		"--format=%H|%an|%ai",
		"--name-only", ".")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(string(out), "\n")

	authorMonthFiles := make(map[string]map[string][]string)
	authorMonthCommits := make(map[string]map[string]int)

	var currentAuthor, currentDate string
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		if strings.Contains(line, "|") && len(line) > 41 && line[40] == '|' {
			// Commit format line: hash|author|date
			parts := strings.SplitN(line, "|", 3)
			if len(parts) >= 3 {
				currentAuthor = parts[1]
				currentDate = parts[2]
				// Count this commit once for the author's month
				month := ""
				if len(currentDate) >= 7 {
					month = currentDate[:7]
				}
				if month != "" && currentAuthor != "" {
					if authorMonthCommits[currentAuthor] == nil {
						authorMonthCommits[currentAuthor] = make(map[string]int)
					}
					authorMonthCommits[currentAuthor][month]++
				}
			}
			continue
		}
		// File path line
		if currentAuthor == "" || currentDate == "" {
			continue
		}
		month := ""
		if len(currentDate) >= 7 {
			month = currentDate[:7]
		}
		if month == "" {
			continue
		}
		if authorMonthFiles[currentAuthor] == nil {
			authorMonthFiles[currentAuthor] = make(map[string][]string)
		}
		authorMonthFiles[currentAuthor][month] = append(authorMonthFiles[currentAuthor][month], line)
	}
	return authorMonthFiles, authorMonthCommits, nil
}

func getTotalFilesByMonth(repoPath string, months []string) map[string]int {
	sorted := make([]string, len(months))
	copy(sorted, months)
	sort.Strings(sorted)

	result := make(map[string]int)
	for _, month := range sorted {
		endOfMonth := month + "-31 23:59:59"
		cmd := exec.Command("git", "-C", repoPath, "log",
			"--all", "--before="+endOfMonth,
			"--format=%H", "-1", ".")
		hashBytes, err := cmd.Output()
		if err != nil {
			continue
		}
		hash := strings.TrimSpace(string(hashBytes))
		if hash == "" {
			continue
		}
		cmd2 := exec.Command("git", "-C", repoPath, "ls-tree", "-r", "--name-only", hash)
		filesBytes, err := cmd2.Output()
		if err != nil {
			continue
		}
		count := 0
		if trimmed := strings.TrimSpace(string(filesBytes)); trimmed != "" {
			count = len(strings.Split(trimmed, "\n"))
		}
		result[month] = count
	}
	return result
}
