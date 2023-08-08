package ctool

import (
	"log"
	"testing"
)

func TestPostJson(t *testing.T) {
	type v struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	rsp, err := PostJson("http://127.0.0.1:9911/webapi/ec/EC_GetRecord",
		v{A: "33", B: 43})
	if err != nil {
		return
	}

	log.Println(string(rsp))

}

func TestGoFound(t *testing.T) {
	type Req struct {
		Id   int         `json:"id"`
		Text string      `json:"text"`
		Doc  interface{} `json:"document"`
	}

	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Addr  string `json:"addr"`
	}

	lines := ReadStrLinesNoFilter("ssssss.log")

	for i := 0; i < len(lines); i++ {
		text := lines[i]
		//println(text)
		var req = &Req{
			Id:   i,
			Text: text,
			Doc: User{
				Name:  text,
				Email: RandomAlpha(13),
				Addr:  RandomAlpha(20),
			},
		}
		PostJson("http://localhost:8080/api/index?database=B1", req)
	}

}
