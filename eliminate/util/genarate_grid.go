package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type GenerateGrid struct {
}
type GridConfig struct {
	ID   int   `json:"id"`
	Row  int   `json:"row"`
	Col  int   `json:"col"`
	Data []int `json:"data"`
}

func (GenerateGrid) ReadConfig() []GridConfig {
	var datas []GridConfig
	fp, _ := os.Open("grid.json")
	dec := json.NewDecoder(fp)
	for {
		err := dec.Decode(&datas)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	for _, line := range datas {
		if line.Col*line.Row > len(line.Data) {
			log.Println("data field is invalid, id=", line.ID)
			return nil
		}
	}

	return datas
}

func (g GridConfig) toString() {
	fmt.Println("\nid: ", g.ID)

	for i := 0; i < g.Row; i++ {
		for j := 0; j < g.Col; j++ {
			fmt.Print(g.Data[i*g.Col+j], " ")
		}
		fmt.Println()
	}
}
