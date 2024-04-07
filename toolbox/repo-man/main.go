package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/sizedpool"
	"github.com/kuangcp/logger"
	"github.com/pkg/errors"
	"io/fs"
	"path/filepath"
	"strings"
)

type RepoAlias struct {
	alias  string
	path   string
	name   string
	ignore bool
}

func (r RepoAlias) String() string {
	var nameColor = ctool.Blue
	if r.ignore {
		nameColor = ctool.Red
	}
	return fmt.Sprintf("%-30s %-50s %-10s",
		ctool.Yellow.PrintNoEnd(r.alias),
		ctool.Green.PrintNoEnd(r.path),
		nameColor.Print(r.name))
}

func PullRepo(dir string) {
	r, err := git.PlainOpen(dir)
	ctool.CheckIfError(err)
	worktree, err := r.Worktree()
	ctool.CheckIfError(err)

	//status, err := worktree.Status()
	//ctool.CheckIfError(err)
	//if !status.IsClean() {
	//	logger.Warn("exist change")
	//	for k, v := range status {
	//		logger.Info(k, string(v.Staging), string(v.Worktree), v.Extra)
	//	}
	//	return
	//}

	remote, err := findRemote(dir)
	ctool.CheckIfError(err)
	//logger.Info(remote)
	err = worktree.Pull(&git.PullOptions{SingleBranch: true, RemoteName: remote})
	if err.Error() == "already up-to-date" {
		return
	}
	logger.Error(dir, remote, err)
}

func findRemote(dir string) (string, error) {
	curRemote := ""
	logger.Info(dir)
	last := ""
	remoteBase := dir + "/.git/refs/remotes/"
	err := filepath.Walk(remoteBase, func(path string, info fs.FileInfo, err error) error {
		tmp := path[len(remoteBase):]
		if strings.HasSuffix(path, "HEAD") {
			curRemote = strings.TrimRight(tmp, "/HEAD")
		}
		if !strings.Contains(tmp, "/") && tmp != "" {
			last = tmp
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if curRemote == "" && last == "" {
		return "", errors.Errorf("No remote setup")
	}
	if curRemote == "" {
		return last, nil
	}
	return curRemote, nil
}

func findAllRemote(dir string) ([]string, error) {
	var all []string
	err := filepath.Walk(dir+"/.git/refs/remotes/", func(path string, info fs.FileInfo, err error) error {
		logger.Info(path, info.Name())
		if info.IsDir() && info.Name() != "remotes" {
			all = append(all, info.Name())
		}
		return nil
	})
	if err != nil {
		return all, err
	}
	return all, nil
}

func PushRepo(dir string) {
	r, err := git.PlainOpen(dir)
	ctool.CheckIfError(err)

	remote, err := findRemote(dir)
	ctool.CheckIfError(err)
	err = r.Push(&git.PushOptions{RemoteName: remote})
	if err != nil {
		logger.Error(dir, err)
	}
}

func PushAllRemote(dir string) {
	r, err := git.PlainOpen(dir)
	ctool.CheckIfError(err)

	remote, err := findAllRemote(dir)
	ctool.CheckIfError(err)
	logger.Info(remote, r)

	for _, rm := range remote {
		err = r.Push(&git.PushOptions{RemoteName: rm})
		if err != nil {
			logger.Error(dir, err)
		}
	}
}

func ShowRepoStatus(dir string) {
	watch := ctool.NewStopWatch()

	watch.Start("open")
	r, err := git.PlainOpen(dir)
	ctool.CheckIfError(err)
	watch.Stop()

	watch.Start("worktree")
	worktree, err := r.Worktree()
	ctool.CheckIfError(err)
	watch.Stop()

	//TODO 重大性能问题
	watch.Start("status")
	status, err := worktree.Status()
	ctool.CheckIfError(err)
	watch.Stop()
	if status.IsClean() {
		return
	}

	temps := strings.Split(dir, "/")

	content := ""
	modify := 0
	add := 0
	for filePath := range status {
		fileStatus := status.File(filePath)
		var color = ctool.End
		if fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified {
			color = ctool.Cyan
			modify++
		}
		if fileStatus.Staging == git.Deleted || fileStatus.Worktree == git.Deleted {
			color = ctool.Red
		}
		if fileStatus.Staging == git.Added || fileStatus.Worktree == git.Added {
			color = ctool.Green
			add++
		}
		if fileStatus.Staging == git.Untracked || fileStatus.Worktree == git.Untracked {
			color = ctool.Yellow
			add++
		}
		content += color.Printf("%c%c    %s\n", fileStatus.Staging, fileStatus.Worktree, filePath)
	}

	fmt.Print("\033[48;5;244m▶ " + fmt.Sprintf("\033[48;5;244m%-17s\033[0m", temps[len(temps)-1]) +
		fmt.Sprintf("\u001B[48;5;244m%-45s\u001B[0m", dir) +
		fmt.Sprintf("\u001B[48;5;244mM:%-3vA:%-3v", modify, add) + " ◀\033[0m\n")
	fmt.Println(content)
	fmt.Println(watch.PrettyPrint())
}

func getRepoList() []RepoAlias {
	home, err := ctool.Home()
	ctool.CheckIfError(err)
	return ctool.ReadLines[RepoAlias](home+"/.repos.sh", func(s string) bool {
		return strings.Contains(s, "alias")
	}, func(s string) RepoAlias {
		temps := strings.Split(s, "'")
		return RepoAlias{alias: temps[0][5 : len(temps[0])-1],
			path:   temps[1][3:],
			name:   strings.TrimSpace(temps[2][2:]),
			ignore: strings.Contains(s, "+"),
		}
	})
}

func ParallelActionRepo(action func(string)) {
	list := getRepoList()

	pool, _ := sizedpool.NewQueuePool(len(list))

	//var latch sync.WaitGroup
	for _, repoAlias := range list {
		if repoAlias.ignore {
			continue
		}
		if repo != "" && repoAlias.name != repo {
			continue
		}

		path := repoAlias.path
		pool.Run(func() {
			action(path)
		})

		//latch.Add(1)
		//go action(repoAlias.path, &latch)
	}
	pool.Wait()
}

func main() {
	helpInfo.Parse()

	if list {
		list := getRepoList()
		for i := range list {
			println(list[i].String())
		}
		return
	} else if fetch {
		ParallelActionRepo(PullRepo)
		return
	} else if status {
		ParallelActionRepo(ShowRepoStatus)
		return
	} else if push {
		ParallelActionRepo(PushRepo)
		return
	} else if pushAll {
		ParallelActionRepo(PushAllRemote)
		return
	}

	ParallelActionRepo(ShowRepoStatus)
}
