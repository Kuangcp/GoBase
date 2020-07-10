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

func (r RepoAlias) String() string {
	var nameColor = cuibase.Blue
	if r.ignore {
		nameColor = cuibase.Red
	}
	return fmt.Sprintf("%-30s %-50s %-10s",
		cuibase.Yellow.PrintNoEnd(r.alias),
		cuibase.Green.PrintNoEnd(r.path),
		nameColor.Print(r.name))
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
	info.PrintHelp()
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
		content += color.Printf("%c%c    %s\n", fileStatus.Staging, fileStatus.Worktree, filePath)
	}

	fmt.Println("▶ " + cuibase.Blue.Printf("%-17s", temps[len(temps)-1]) +
		cuibase.Green.Printf("%-45s", dir) +
		cuibase.Blue.Printf("M:%-3vA:%-3v", modify, add) + " ◀\n")
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
