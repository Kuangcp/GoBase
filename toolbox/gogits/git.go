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