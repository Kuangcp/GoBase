package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"os"
	"sort"
	"strings"
)

var (
	cfgFile = "/.ggi.json"
)

type Config struct {
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Alias   string `json:"alias"`
	Path    string `json:"path"`
	Comment string `json:"comment"`
}

func init() {
	home, err := ctool.Home()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgFile = home + cfgFile
}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}

	if addRepo != "" {
		addRepos()
		return
	} else if delRepo != "" {
		deleteRepos()
		return
	} else if jumpRepo != "" {
		cfg := Read()
		for _, r := range cfg.Repos {
			if r.Alias == jumpRepo {
				fmt.Print(r.Path)
				return
			}
		}
		return
	} else if listRepo {
		cfg := Read()
		for _, r := range cfg.Repos {
			fmt.Println(r.Alias, r.Path)
		}
		return
	}

	if pushCur {
		pushCurDir()
		return
	}

	if push {
		if allRepo {
			pushAllRepo()
			return
		} else {
			pushConfigedRepos()
			return
		}
	}
	if pull {
		pullRepo()
		return
	}

	if lod || lods {
		runLod(lods)
		return
	}

	checkRepoChange()
}

func deleteRepos() {
	cfg := Read()
	var nlist []Repo
	find := false
	for _, r := range cfg.Repos {
		if r.Alias != delRepo {
			nlist = append(nlist, r)
		} else {
			find = true
			logger.Info("Delete repo %s", delRepo)
		}
	}
	if find {
		cfg.Repos = nlist
		Write(cfg)
	} else {
		logger.Error("Repo %s not found", delRepo)
	}
}

func addRepos() {
	dir := FindRootDir()
	if dir == "" {
		logger.Error("None git repo")
		return
	}
	cfg := Read()
	for _, r := range cfg.Repos {
		if r.Path == dir {
			logger.Error("Repo %s already exists", dir)
			return
		}
	}
	repo := Repo{Path: dir, Alias: addRepo}
	cfg.Repos = append(cfg.Repos, repo)
	Write(cfg)
	return
}

func pushConfigedRepos() {
	cfg := Read()
	for _, repo := range cfg.Repos {
		out, err := runGit(repo.Path, "status", "--short", "--branch")
		if err != nil {
			logger.Error("Repo %s check error: %v", repo.Alias, err)
			continue
		}
		if !strings.Contains(out, "[ahead") {
			continue
		}

		fmt.Printf("\033[35m%-20s\033[0m %s\n", repo.Alias, repo.Path)
		msg, err := runGitErr(repo.Path, "push")
		if err != nil {
			logger.Error("push error: %v\n%s", err, msg)
		} else if msg != "" {
			fmt.Println(msg)
		}
	}
}

func pushCurDir() {
	rootDir := FindRootDir()
	if rootDir == "" {
		return
	}
	out, err := runGit(rootDir, "remote")
	if err != nil || out == "" {
		logger.Error("Repo %s has no remote", rootDir)
		return
	}
	for _, remote := range strings.Split(out, "\n") {
		remote = strings.TrimSpace(remote)
		if remote == "" {
			continue
		}
		msg, err := runGitErr(rootDir, "push", remote)
		if err != nil {
			logger.Error("Repo %s push %s error: %v\n%s", rootDir, remote, err, msg)
		} else if msg != "" {
			logger.Info("Repo %s %s: %s", rootDir, remote, msg)
		}
	}
}

func pushAllRepo() {
	cfg := Read()
	for _, repo := range cfg.Repos {
		out, err := runGit(repo.Path, "remote")
		if err != nil || out == "" {
			logger.Error("Repo %s has no remote", repo.Alias)
			continue
		}
		for _, remote := range strings.Split(out, "\n") {
			remote = strings.TrimSpace(remote)
			if remote == "" {
				continue
			}
			msg, err := runGitErr(repo.Path, "push", remote)
			if err != nil {
				logger.Error("Repo %s push %s error: %v\n%s", repo.Alias, remote, err, msg)
			} else if msg != "" {
				logger.Info("Repo %s %s: %s", repo.Alias, remote, msg)
			}
		}
	}
}

func pullRepo() {
	cfg := Read()
	for _, repo := range cfg.Repos {
		msg, err := runGitErr(repo.Path, "pull")
		if err != nil {
			logger.Error("Repo %s pull error: %v\n%s", repo.Alias, err, msg)
		} else {
			logger.Info("Repo %s: %s", repo.Alias, strings.ReplaceAll(msg, "\n", " | "))
		}
	}
}

func runLod(sorted bool) {
	entries, err := os.ReadDir(".")
	if err != nil {
		logger.Error("read dir error: %v", err)
		return
	}

	type dirCount struct {
		name  string
		count int
	}
	var results []dirCount
	maxLen := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) > maxLen {
			maxLen = len(name)
		}

		args := []string{"log", "--oneline", "--all"}
		if afterDate != "" {
			args = append(args, "--after="+afterDate)
		}
		if beforeDate != "" {
			args = append(args, "--before="+beforeDate)
		}
		args = append(args, "--", name)

		out, err := runGit(".", args...)
		if err != nil {
			results = append(results, dirCount{name, 0})
			continue
		}
		count := 0
		if out != "" {
			count = len(strings.Split(out, "\n"))
		}
		results = append(results, dirCount{name, count})
	}

	if sorted {
		sort.Slice(results, func(i, j int) bool {
			return results[i].count > results[j].count
		})
	}

	format := fmt.Sprintf("%%4d %%-%ds\n", maxLen)
	for _, r := range results {
		fmt.Printf(format, r.count, r.name)
	}
}

func checkRepoChange() {
	cfg := Read()
	for _, repo := range cfg.Repos {
		out, err := runGit(repo.Path, "status", "-s")
		if err != nil || out == "" {
			continue
		}

		fmt.Printf("\033[32m%-20s\033[0m \033[36m%s\033[0m\n", repo.Alias, repo.Path)
		for _, line := range strings.Split(out, "\n") {
			if line != "" {
				fmt.Printf("  %s\n", line)
			}
		}
		fmt.Println()
	}
}

func FindRootDir() string {
	dir, _ := os.Getwd()
	exist := ctool.IsFileExist(dir + "/.git")
	depth := 1
	for !exist {
		depth++
		if depth > 7 {
			logger.Error("Max depth exceeded")
			return ""
		}
		parts := strings.Split(dir, "/")
		parts = parts[:len(parts)-1]
		dir = strings.Join(parts, "/")
		if dir == "/" {
			logger.Error("None Any Git Repository")
			return ""
		}
		exist = ctool.IsFileExist(dir + "/.git")
	}
	return dir
}
func Read() *Config {
	if !ctool.IsFileExist(cfgFile) {
		err := os.WriteFile(cfgFile, []byte("{\"repos\":[]}"), 0644)
		if err != nil {
			logger.Fatal(err)
		}
	}
	file, err := os.ReadFile(cfgFile)
	if err != nil {
		logger.Error(err)
		return nil
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &cfg
}
func Write(cfg *Config) {
	marshal, err := json.Marshal(cfg)
	if err != nil {
		logger.Error(err)
	}
	err = os.WriteFile(cfgFile, marshal, 0644)
	if err != nil {
		logger.Error(err)
	}
}