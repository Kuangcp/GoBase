package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/sizedpool"
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
	} else if jumpRepo != "" {
		cfg := Read()
		for _, r := range cfg.Repos {
			if r.Alias == jumpRepo {
				fmt.Print(r.Path)
				return
			}
		}
		return
	} else if delRepo != "" {
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
	} else if listRepo {
		cfg := Read()
		for _, r := range cfg.Repos {
			fmt.Println(r.Alias, r.Path)
		}
		return
	}

	if push {
		if allRepo {
			cfg := Read()
			for _, repo := range cfg.Repos {
				r, err := git.PlainOpen(repo.Path)
				if err != nil {
					logger.Error("Repo %s not found: %v", repo.Alias, err)
					continue
				}
				list, err := r.Remotes()
				if err != nil {
					logger.Error("Repo %s open error: %v", repo.Alias, err)
					return
				}
				for _, remote := range list {
					name := remote.Config().Name
					err := r.Push(&git.PushOptions{RemoteName: name})
					if err != nil {
						if err.Error() == "already up-to-date" {
							logger.Info("Repo %s %s already up-to-date", repo.Alias, name)
						} else {
							logger.Error("Repo %s push %s error %v", repo.Alias, name, err)
						}
					}
				}
			}
		} else {
			rootDir := FindRootDir()
			if rootDir == "" {
				logger.Error("None git repo")
				return
			}

			r, err := git.PlainOpen(rootDir)
			if err != nil {
				logger.Error("Repo %s not found: %v", rootDir, err)
				return
			}

			list, err := r.Remotes()
			if err != nil {
				logger.Error("Repo %s open error: %v", rootDir, err)
				return
			}
			group, _ := sizedpool.New(sizedpool.PoolOption{
				Size: 4,
			})
			for _, remote := range list {
				name := remote.Config().Name
				group.Run(func() {
					err := r.Push(&git.PushOptions{RemoteName: name})
					if err != nil {
						if err.Error() == "already up-to-date" {
							logger.Info("Repo %s %s already up-to-date", rootDir, name)
						} else {
							logger.Error("Repo %s push %s error %v", rootDir, name, err)
						}
					}
				})
			}
			group.Wait()
		}
	}
	if pull {
		cfg := Read()
		group, _ := sizedpool.New(sizedpool.PoolOption{
			Size: 4,
		})
		for _, repo := range cfg.Repos {
			group.Run(func() {
				r, err := git.PlainOpen(repo.Path)
				if err != nil {
					logger.Error("Repo %s not found: %v", repo.Alias, err)
					return
				}
				w, err := r.Worktree()
				if err != nil {
					logger.Error("Repo %s open error: %v", repo.Alias, err)
					return
				}

				// chmod 600 ~/.ssh/id_rsa
				// ssh-add ~/.ssh/id_rsa
				// Pull the latest changes from the origin remote and merge into the current branch
				logger.Info("Try pull repo %s", repo.Alias)
				err = w.Pull(&git.PullOptions{})
				if err != nil {
					if err.Error() == "already up-to-date" {
						logger.Info("Repo %s already up-to-date", repo.Alias)
					} else {
						logger.Error("Repo %s pull error %v", repo.Alias, err)
					}
				}
			})
		}
		group.Wait()
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
