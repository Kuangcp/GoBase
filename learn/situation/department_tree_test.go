package situation

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/stopwatch"
	"github.com/kuangcp/logger"
	"math/rand"
	"testing"
)

var list = buildDepList()

const repeatCount = 100
const departmentCount = 10
const writeFile = false

var runCounter = 0

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

	for i := 10; i < departmentCount; i++ {
		p := rand.Intn(8) + 1
		deps = append(deps, Dep{id: i, name: fmt.Sprint("a", i), parentId: p, level: p2Level[p]})
	}
	return deps
}

func buildTreeByLoop() *DepVO {
	levelMap := make(map[int][]Dep)
	depMap := make(map[int]*DepVO)
	maxLevel := 0
	minLevel := 1
	for _, v := range list {
		runCounter++
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
	//logger.Info(levelMap)
	minDep := levelMap[minLevel]
	if len(minDep) > 1 {
		logger.Warn("more than one root dep")
	}
	dep := minDep[0]
	rootDep := newDepVO(dep.id, dep.name)
	depMap[dep.id] = rootDep

	for i := minLevel; i < maxLevel; i++ {
		deps := levelMap[i]
		for _, v := range deps {
			runCounter++
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
				//logger.Warn("not exist", v)
			}
			depMap[v.id] = vo
		}
	}
	return rootDep
}

func findParentAndAppend(cid int, dep map[int]Dep, cache map[int]*DepVO) *DepVO {
	logger.Info("seek", cid)
	runCounter++
	d, ok := dep[cid]
	if !ok {
		logger.Warn("parent id is invalid")
		return nil
	}
	vo, ok := cache[cid]
	if ok {
		return vo
	}
	// 递归终点
	if cid == 1 {
		curDepVO := newDepVO(d.id, d.name)
		cache[cid] = curDepVO
		return curDepVO
	}

	parent := findParentAndAppend(d.parentId, dep, cache)
	if parent == nil {
		return nil
	}

	curDepVO := newDepVO(d.id, d.name)
	cache[cid] = curDepVO

	parent.appendChild(curDepVO)
	logger.Info(curDepVO.Id, "->", parent.Id)
	return curDepVO
}

func buildTreeByRecursive() *DepVO {
	depMap := make(map[int]Dep)
	for i := range list {
		runCounter++
		dep := list[i]
		depMap[dep.id] = dep
	}

	resultMap := make(map[int]*DepVO)

	for i := range list {
		runCounter++
		dep := list[i]
		//logger.Info(dep.id, dep.parentId)
		findParentAndAppend(dep.id, depMap, resultMap)
	}
	return resultMap[1]
}

func TestLoop(t *testing.T) {
	rootDep := buildTreeByLoop()

	if writeFile {
		marshal, _ := json.Marshal(rootDep)
		logger.Info(string(marshal))
		writer, _ := ctool.NewWriter("dep.json", true)
		defer writer.Close()
		writer.Write(marshal)
	}
	logger.Info(runCounter)
}

func TestRecursive(t *testing.T) {
	rootDep := buildTreeByRecursive()
	if writeFile {
		marshal, _ := json.Marshal(rootDep)
		writer, _ := ctool.NewWriter("dep2.json", true)
		defer writer.Close()
		writer.Write(marshal)
	}

	logger.Info(runCounter)
}

// 递归实现代码更简单，但是性能略差一些 可预料的是如果层级更深，差距更大
func TestCompare(t *testing.T) {
	watch := stopwatch.New()

	watch.Start("loop")
	runCounter = 0
	for i := 0; i < repeatCount; i++ {
		buildTreeByLoop()
	}
	logger.Info("loop run count", runCounter)
	watch.Stop()

	watch.Start("recursive")
	runCounter = 0
	for i := 0; i < repeatCount; i++ {
		buildTreeByRecursive()
	}
	logger.Info("recursive run count", runCounter)
	watch.Stop()

	logger.Info(watch.PrettyPrint())
}
