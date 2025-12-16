package ctool

import (
	"bytes"
	"io"
	"testing"
)

// go test -bench=. -benchmem -benchtime=3s base_benchmark_test.go base.go

// 生成测试数据
func generateTestData(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

// 基准测试：使用 Adapter.CopyStream (带内存池)
func BenchmarkAdapterCopyStream_Small(b *testing.B) {
	data := generateTestData(1024) // 1KB
	adapter := New()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src := io.NopCloser(bytes.NewReader(data))
		result, rc := adapter.CopyStream(src)
		if result == nil || rc == nil {
			b.Fatal("CopyStream failed")
		}
		rc.Close()
	}
}

func BenchmarkAdapterCopyStream_Medium(b *testing.B) {
	data := generateTestData(64 * 1024) // 64KB
	adapter := New()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src := io.NopCloser(bytes.NewReader(data))
		result, rc := adapter.CopyStream(src)
		if result == nil || rc == nil {
			b.Fatal("CopyStream failed")
		}
		rc.Close()
	}
}

func BenchmarkAdapterCopyStream_Large(b *testing.B) {
	data := generateTestData(1024 * 1024) // 1MB
	adapter := New()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src := io.NopCloser(bytes.NewReader(data))
		result, rc := adapter.CopyStream(src)
		if result == nil || rc == nil {
			b.Fatal("CopyStream failed")
		}
		rc.Close()
	}
}

// 基准测试：使用普通函数 CopyStream (无内存池)
func BenchmarkFuncCopyStream_Small(b *testing.B) {
	data := generateTestData(1024) // 1KB
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src := io.NopCloser(bytes.NewReader(data))
		result, rc := CopyStream(src)
		if result == nil || rc == nil {
			b.Fatal("CopyStream failed")
		}
		rc.Close()
	}
}

func BenchmarkFuncCopyStream_Medium(b *testing.B) {
	data := generateTestData(64 * 1024) // 64KB
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src := io.NopCloser(bytes.NewReader(data))
		result, rc := CopyStream(src)
		if result == nil || rc == nil {
			b.Fatal("CopyStream failed")
		}
		rc.Close()
	}
}

func BenchmarkFuncCopyStream_Large(b *testing.B) {
	data := generateTestData(1024 * 1024) // 1MB
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src := io.NopCloser(bytes.NewReader(data))
		result, rc := CopyStream(src)
		if result == nil || rc == nil {
			b.Fatal("CopyStream failed")
		}
		rc.Close()
	}
}

// 并发场景基准测试
func BenchmarkAdapterCopyStream_Parallel(b *testing.B) {
	data := generateTestData(64 * 1024) // 64KB
	adapter := New()
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			src := io.NopCloser(bytes.NewReader(data))
			result, rc := adapter.CopyStream(src)
			if result == nil || rc == nil {
				b.Fatal("CopyStream failed")
			}
			rc.Close()
		}
	})
}

func BenchmarkFuncCopyStream_Parallel(b *testing.B) {
	data := generateTestData(64 * 1024) // 64KB
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			src := io.NopCloser(bytes.NewReader(data))
			result, rc := CopyStream(src)
			if result == nil || rc == nil {
				b.Fatal("CopyStream failed")
			}
			rc.Close()
		}
	})
}

