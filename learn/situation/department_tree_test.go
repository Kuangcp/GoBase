package situation

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"math/rand"
	"testing"
)

type (
	Dep struct {
		id       int
		name     string
		parentId int
		level    int
	}
	DepVO struct {
		Id    int      `json:"id"`
		Name  string   `json:"name"`
		Child []*DepVO `json:"child"`
	}
)

func newDepVO(id int, name string) *DepVO {
	var child []*DepVO
	return &DepVO{Id: id, Name: name, Child: child}
}

func (d *DepVO) appendChild(c *DepVO) {
	d.Child = append(d.Child, c)
}

func buildDepList() []Dep {
	var deps []Dep
	deps = append(deps, Dep{id: 1, name: "a1", parentId: 0, level: 1})
	deps = append(deps, Dep{id: 2, name: "a2", parentId: 1, level: 2})
	deps = append(deps, Dep{id: 3, name: "a3", parentId: 1, level: 2})
	deps = append(deps, Dep{id: 4, name: "a4", parentId: 3, level: 3})
	deps = append(deps, Dep{id: 5, name: "a5", parentId: 4, level: 4})
	deps = append(deps, Dep{id: 6, name: "a6", parentId: 4, level: 4})
	deps = append(deps, Dep{id: 7, name: "a7", parentId: 5, level: 5})
	deps = append(deps, Dep{id: 8, name: "a8", parentId: 6, level: 6})

	p2Level := make(map[int]int)
	p2Level[1] = 2
	p2Level[2] = 3
	p2Level[3] = 3
	p2Level[4] = 4
	p2Level[5] = 5
	p2Level[6] = 5
	p2Level[7] = 6
	p2Level[8] = 7

	for i := 10; i < 600; i++ {
		p := rand.Intn(8) + 1
		deps = append(deps, Dep{id: i, name: fmt.Sprint("a", i), parentId: p, level: p2Level[p]})
	}
	return deps
}

func TestLoop(t *testing.T) {
	list := buildDepList()

	levelMap := make(map[int][]Dep)
	depMap := make(map[int]*DepVO)
	maxLevel := 0
	minLevel := 1
	for _, v := range list {
		if maxLevel < v.level {
			maxLevel = v.level
		}
		if minLevel > v.level {
			minLevel = v.level
		}

		l, ok := levelMap[v.level]
		if ok {
			l = append(l, v)
			levelMap[v.level] = l
		} else {
			var tmp []Dep
			tmp = append(tmp, v)
			levelMap[v.level] = tmp
		}
	}
	logger.Info(levelMap)
	minDep := levelMap[minLevel]
	if len(minDep) > 1 {
		logger.Warn("")
	}
	dep := minDep[0]
	rootDep := newDepVO(dep.id, dep.name)
	depMap[dep.id] = rootDep

	for i := minLevel; i < maxLevel; i++ {
		deps := levelMap[i]
		for _, v := range deps {
			var vo *DepVO
			if v.level == minLevel {
				vo = rootDep
			} else {
				vo = newDepVO(v.id, v.name)
			}

			parent, ok := depMap[v.parentId]
			if ok {
				parent.appendChild(vo)
			} else {
				logger.Warn("not exist", v)
			}
			depMap[v.id] = vo
		}
	}

	logger.Info(rootDep)

	marshal, err := json.Marshal(rootDep)
	if err != nil {
		return
	}
	//logger.Info(string(marshal))
	writer, _ := ctool.NewWriter("dep.json", true)
	defer writer.Close()
	writer.Write(marshal)
}

// TODO
func TestRecursive(t *testing.T) {
	//list := buildDepList()
	
}
