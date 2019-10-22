package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func clone(url string) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
	if err != nil {

	}
	ref, err := r.Head()
	print(err, ref)
}

func main() {
	clone("https://github.com/Kuangcp/GoBase")
	//cuibase.AssertParamCount(1, "must input url")

}
