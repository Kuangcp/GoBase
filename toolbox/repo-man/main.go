package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
	"strings"
	"sync"
)

type RepoAlias struct {
	alias string
	path  string
	name  string
}

func (this RepoAlias) String() string {
	return fmt.Sprintf("%v%-30s %v%-50s %v%-10v%v",
		cuibase.Yellow, this.alias,
		cuibase.Green, this.path,
		cuibase.Blue, this.name, cuibase.End)
}

func HelpInfo(_ []string) {
	info := cuibase.HelpInfo{
		Description: "Git repository manager",
		VerbLen:     -3,
		ParamLen:    -5,
		Params: []cuibase.ParamInfo{
			{
				Verb:    "-h",
				Param:   "",
				Comment: "Help info",
			},
		}}
	cuibase.Help(info)
}
func PullRepo(dir string, latch *sync.WaitGroup) {
	defer latch.Done()

	r, err := git.PlainOpen(dir)
	cuibase.CheckIfError(err)
	worktree, err := r.Worktree()
	cuibase.CheckIfError(err)

	// TODO ssh bug
	err = worktree.Pull(&git.PullOptions{SingleBranch: true})
	logger.Error(dir, err)
}

func ShowRepoStatus(dir string, latch *sync.WaitGroup) {
	defer latch.Done()

	r, err := git.PlainOpen(dir)
	cuibase.CheckIfError(err)
	worktree, err := r.Worktree()
	cuibase.CheckIfError(err)
	status, err := worktree.Status()
	cuibase.CheckIfError(err)
	if status.IsClean() {
		return
	}

	temps := strings.Split(dir, "/")
	fmt.Printf("%v %-20v  %v%v %v\n", cuibase.Yellow, temps[len(temps)-1],
		cuibase.Green, dir, cuibase.End)

	for filePath := range status {
		fileStatus := status.File(filePath)
		var color = cuibase.End
		if fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified {
			color = cuibase.Cyan
		}
		if fileStatus.Staging == git.Deleted || fileStatus.Worktree == git.Deleted {
			color = cuibase.Red
		}
		if fileStatus.Staging == git.Added || fileStatus.Worktree == git.Added {
			color = cuibase.Green
		}
		if fileStatus.Staging == git.Untracked || fileStatus.Worktree == git.Untracked {
			color = cuibase.Yellow
		}
		fmt.Printf("   %v%c%c    %s%s\n", color, fileStatus.Staging, fileStatus.Worktree, filePath, cuibase.End)
	}
	fmt.Println()
}

func getRepoList() []interface{} {
	home, err := cuibase.Home()
	cuibase.CheckIfError(err)
	return cuibase.ReadFileLines(home+"/.repos.sh", func(s string) bool {
		return strings.Contains(s, "alias") && !strings.Contains(s, "+")
	}, func(s string) interface{} {
		temps := strings.Split(s, "'")
		return RepoAlias{alias: temps[0][5 : len(temps[0])-1], path: temps[1][3:], name: strings.TrimSpace(temps[2][2:])}
	})
}

func main() {
	logger.SetLogPathTrim("/toolbox/")
	cuibase.RunAction(map[string]func(params []string){
		"-h": HelpInfo,
		"l": func(params []string) {
			list := getRepoList()
			for i := range list {
				println(list[i].(RepoAlias).String())
			}
		},
		"pla": func(params []string) {
			list := getRepoList()
			var latch sync.WaitGroup
			for _, s := range list {
				latch.Add(1)
				go PullRepo(s.(RepoAlias).path, &latch)
			}
			latch.Wait()
		},
	}, func(params []string) {
		list := getRepoList()
		var latch sync.WaitGroup
		for _, s := range list {
			latch.Add(1)
			go ShowRepoStatus(s.(RepoAlias).path, &latch)
		}
		latch.Wait()
	})
}
