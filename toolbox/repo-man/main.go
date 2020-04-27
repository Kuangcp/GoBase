package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
)

type RepoAlias struct {
	alias  string
	path   string
	name   string
	ignore bool
}

func (this RepoAlias) String() string {
	var nameColor = cuibase.Blue
	if this.ignore {
		nameColor = cuibase.Red
	}
	return fmt.Sprintf("%v%-30s %v%-50s %v%-10v%v",
		cuibase.Yellow, this.alias,
		cuibase.Green, this.path,
		nameColor, this.name, cuibase.End)
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

func PushRepo(dir string, latch *sync.WaitGroup) {
	defer latch.Done()

	r, err := git.PlainOpen(dir)
	cuibase.CheckIfError(err)

	err = r.Push(&git.PushOptions{})
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

	content := ""
	modify := 0
	add := 0
	for filePath := range status {
		fileStatus := status.File(filePath)
		var color = cuibase.End
		if fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified {
			color = cuibase.Cyan
			modify++
		}
		if fileStatus.Staging == git.Deleted || fileStatus.Worktree == git.Deleted {
			color = cuibase.Red
		}
		if fileStatus.Staging == git.Added || fileStatus.Worktree == git.Added {
			color = cuibase.Green
			add++
		}
		if fileStatus.Staging == git.Untracked || fileStatus.Worktree == git.Untracked {
			color = cuibase.Yellow
			add++
		}
		content += fmt.Sprintf("   %v%c%c    %s%s\n",
			color, fileStatus.Staging, fileStatus.Worktree, filePath, cuibase.End)
	}
	fmt.Printf("%v▶ %-20v  %v%-50v %vM:%-3vA:%-3v ◀%v\n",
		cuibase.Blue, temps[len(temps)-1],
		cuibase.Green, dir,
		cuibase.Blue, modify, add,
		cuibase.End)

	fmt.Println(content)
}

func getRepoList() []interface{} {
	home, err := cuibase.Home()
	cuibase.CheckIfError(err)
	return cuibase.ReadFileLines(home+"/.repos.sh", func(s string) bool {
		return strings.Contains(s, "alias")
	}, func(s string) interface{} {
		temps := strings.Split(s, "'")
		return RepoAlias{alias: temps[0][5 : len(temps[0])-1],
			path:   temps[1][3:],
			name:   strings.TrimSpace(temps[2][2:]),
			ignore: strings.Contains(s, "+"),
		}
	})
}

func ParallelActionRepo(action func(string, *sync.WaitGroup)) {
	list := getRepoList()
	var latch sync.WaitGroup
	for _, s := range list {
		repoAlias := s.(RepoAlias)
		if repoAlias.ignore {
			continue
		}

		latch.Add(1)
		go action(repoAlias.path, &latch)
	}
	latch.Wait()
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
		"pla": func(_ []string) {
			ParallelActionRepo(PullRepo)
		},
		"pa": func(_ []string) {
			ParallelActionRepo(PushRepo)
		},
	}, func(_ []string) {
		ParallelActionRepo(ShowRepoStatus)
	})
}
