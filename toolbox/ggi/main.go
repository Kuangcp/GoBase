package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"os"
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
		dir := FindRootDir()
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
	if jumpRepo != "" {
		cfg := Read()
		for _, r := range cfg.Repos {
			if r.Alias == jumpRepo {
				fmt.Print(r.Path)
				return
			}
		}
		return
	}

	if delRepo != "" {
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
		return
	}

	if listRepo {
		cfg := Read()
		for _, r := range cfg.Repos {
			fmt.Println(r.Alias, r.Path)
		}
		return
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
