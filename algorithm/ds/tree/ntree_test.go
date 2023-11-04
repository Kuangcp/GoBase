package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"github.com/kuangcp/gobase/pkg/ctool/stream"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func appendSon(id string, all []ATable[string], depth int) []ATable[string] {
	if depth > 2 {
		return all
	}
	for i := 0; i < 2; i++ {
		ranInt := rand.Intn(9)
		if ranInt%4 == 3 {
			temId := ctool.RandomAlpha(4)
			all = append(all, ATable[string]{Id: temId, ParentId: id, Data: ctool.RandomAlpha(2)})
			all = appendSon(temId, all, depth+1)
		}
	}
	return all
}

func TestGenerateATable(t *testing.T) {
	var all []ATable[string]
	rand.Seed(uint64(time.Now().UnixMilli()))
	for i := 0; i < 4; i++ {
		id := ctool.RandomAlpha(4)
		node := ATable[string]{Id: id, Data: ctool.RandomAlpha(2)}
		all = append(all, node)

		all = appendSon(id, all, 0)
	}
	for i := range all {
		fmt.Println(all[i])
	}
	x, _ := json.Marshal(all)
	fmt.Println(string(x))
}

func TestTableToATree(t *testing.T) {
	var all []ATable[string]
	err := json.Unmarshal([]byte("[{\"id\":\"utuK\",\"parent_id\":\"\",\"data\":\"IE\"},{\"id\":\"zXxg\",\"parent_id\":\"utuK\",\"data\":\"qr\"},{\"id\":\"QXma\",\"parent_id\":\"zXxg\",\"data\":\"oY\"},{\"id\":\"jUMu\",\"parent_id\":\"zXxg\",\"data\":\"fW\"},{\"id\":\"Fdwp\",\"parent_id\":\"\",\"data\":\"DJ\"},{\"id\":\"VrGA\",\"parent_id\":\"Fdwp\",\"data\":\"jZ\"},{\"id\":\"SFFc\",\"parent_id\":\"VrGA\",\"data\":\"bu\"},{\"id\":\"DgjH\",\"parent_id\":\"\",\"data\":\"RF\"},{\"id\":\"ZBfK\",\"parent_id\":\"DgjH\",\"data\":\"lK\"},{\"id\":\"IGEl\",\"parent_id\":\"ZBfK\",\"data\":\"Va\"},{\"id\":\"ehNG\",\"parent_id\":\"IGEl\",\"data\":\"JK\"},{\"id\":\"lSie\",\"parent_id\":\"IGEl\",\"data\":\"AZ\"},{\"id\":\"HhYg\",\"parent_id\":\"\",\"data\":\"iI\"},{\"id\":\"bkDO\",\"parent_id\":\"HhYg\",\"data\":\"yC\"},{\"id\":\"zNAH\",\"parent_id\":\"bkDO\",\"data\":\"vO\"},{\"id\":\"YktD\",\"parent_id\":\"bkDO\",\"data\":\"fl\"}]"), &all)
	if err != nil {
		return
	}

	var tree []*ATree[string]
	layerCache := make(map[string][]ATable[string])
	idx := make(map[string]*ATree[string])

	stream.Just(all...).Group(func(item any) any {
		return item.(ATable[string]).ParentId
	}).ForEach(func(item any) {
		i := item.(stream.GroupItem)
		//fmt.Println(i)
		_, ok := layerCache[i.Key.(string)]
		if !ok {
			layerCache[i.Key.(string)] = []ATable[string]{}
		}
		for _, v := range i.Val {
			layerCache[i.Key.(string)] = append(layerCache[i.Key.(string)], v.(ATable[string]))
		}
	})

	topLevel := layerCache[""]
	var layers []string
	for _, t := range topLevel {
		node := &ATree[string]{Id: t.Id, Data: t.Data}
		idx[t.Id] = node
		tree = append(tree, node)
		layers = append(layers, t.Id)
	}
	for len(layers) > 0 {
		layers = fillNextLayer(layerCache, idx, layers)
	}

	for i, t := range tree {
		fmt.Println("start search")
		t.PrintJson()
		algo.WriteNMindMap(t, fmt.Sprint(i)+"atree.pu")
		s := t.Search("fW", true)
		if s {
			fmt.Println("MATCHED >>>>>>>>>>>")
			t.PrintJson()
		} else {
			fmt.Println("NOT FOUND >>>>>>>>>>>")
		}
		fmt.Print("\n\n\n\n")
	}

}

func fillNextLayer(layerCache map[string][]ATable[string], idx map[string]*ATree[string], layers []string) []string {
	var nextLayer []string
	for _, id := range layers {
		_, ok := layerCache[id]
		if !ok {
			continue
		}

		p := idx[id]
		next := layerCache[id]
		for _, t := range next {
			nextLayer = append(nextLayer, t.Id)
			node := &ATree[string]{Id: t.Id, Data: t.Data}
			idx[t.Id] = node
			p.Childes = append(p.Childes, node)
		}
	}
	return nextLayer
}

func TestDeleteSlice(t *testing.T) {
	x := []int{1, 2, 3, 4}
	x = append(x[:3], x[4:]...)
	fmt.Println(x)
}
