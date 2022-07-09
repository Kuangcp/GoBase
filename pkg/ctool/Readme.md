## Debug

```go
import _ "net/http/pprof"

go func() {
    http.ListenAndServe("0.0.0.0:8897", nil)
}()

```

go tool pprof  http://localhost:8897/debug/pprof/heap

采样
go tool pprof -seconds=10 -http=:9981  http://localhost:8897/debug/pprof/heap