package main

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io"
	"log"
	"net/http"
)

func main() {
	//compileAndRun()
	//defFuncInvoke()
	dynamicServer()
}

func dynamicServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			a := recover()
			log.Println(a)
			writer.Write([]byte("error"))
		}()

		bts, err := io.ReadAll(request.Body)
		if len(bts) == 0 {
			writer.Write([]byte("empty"))
			return
		}
		if err != nil {
			log.Println(err)
			writer.Write([]byte("error"))
			return
		}
		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)

		_, err = i.Eval(string(bts))
		if err != nil {
			log.Println(err)
			writer.Write([]byte("error"))
			return
		}

		v, err := i.Eval("cus.Handle")
		if err != nil {
			log.Println(err)
			writer.Write([]byte("error"))
			return
		}

		handler := v.Interface().(func() string)
		writer.Write([]byte(handler()))

		//handler := v.Interface().(func(writer http.ResponseWriter, request *http.Request))
		//handler(writer, request)
	})
	http.ListenAndServe(":9092", nil)

}
func defFuncInvoke() {
	const src = `package foo
func Bar(s string) string { return s + "-Foo" }`

	i := interp.New(interp.Options{})

	_, err := i.Eval(src)
	if err != nil {
		panic(err)
	}

	v, err := i.Eval("foo.Bar")
	if err != nil {
		panic(err)
	}

	bar := v.Interface().(func(string) string)

	r := bar("Kung")
	println(r)
}

func compileAndRun() {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)

	_, err := i.Eval(`import "fmt"`)
	if err != nil {
		panic(err)
	}

	_, err = i.Eval(`fmt.Println("Hello Yaegi")`)
	if err != nil {
		panic(err)
	}
}
