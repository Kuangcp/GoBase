package ctool

import (
	"bytes"
	"io"
	"log"
	"sync"
)

type Adapter struct {
	pool sync.Pool
}

func New() *Adapter {
	return &Adapter{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 4096)
			},
		},
	}
}

// pooledReadCloser 包装 io.ReadCloser，在 Close 时归还 buffer 到 pool
type pooledReadCloser struct {
	io.ReadCloser
	adapter *Adapter
	buffer  []byte
	once    sync.Once
}

func (prc *pooledReadCloser) Close() error {
	err := prc.ReadCloser.Close()
	// 确保 buffer 只被归还一次
	prc.once.Do(func() {
		if prc.buffer != nil {
			prc.adapter.pool.Put(prc.buffer)
			prc.buffer = nil
		}
	})
	return err
}

// CopyStream 高效地复制流，使用内存池复用 buffer
// 返回完整的数据副本和可读取的 ReadCloser
// buffer 会在返回的 ReadCloser 被 Close 时归还到 pool
func (ada *Adapter) CopyStream(src io.ReadCloser) ([]byte, io.ReadCloser) {
	// 从 pool 获取 buffer 作为临时读取缓冲区
	buffer := ada.pool.Get().([]byte)
	clear(buffer)

	// 读取完整数据，使用 buffer 作为临时缓冲区
	bodyBt, err := ReadAll(src, buffer)
	if err != nil {
		log.Println(err)
		// 读取失败时立即归还 buffer
		ada.pool.Put(buffer)
		return nil, nil
	}

	// 创建独立的数据副本，不引用 pool 中的 buffer
	// 这样即使 buffer 被归还，返回的数据也是安全的
	dataCopy := make([]byte, len(bodyBt))
	copy(dataCopy, bodyBt)

	// 创建包装的 ReadCloser，在 Close 时归还 buffer
	reader := &pooledReadCloser{
		ReadCloser: io.NopCloser(bytes.NewBuffer(dataCopy)),
		adapter:    ada,
		buffer:     buffer,
	}

	return dataCopy, reader
}

// ReadAll 高效读取所有数据，使用提供的 buffer 作为临时缓冲区
// 返回的数据可能会引用 buffer，调用者需要复制数据以确保安全
func ReadAll(r io.Reader, b []byte) ([]byte, error) {
	// 重置 buffer 长度，但保持容量
	b = b[:0]

	for {
		// 计算可用空间
		available := cap(b) - len(b)
		if available == 0 {
			// 需要扩容，让 append 自动处理扩容策略
			b = append(b, 0)[:len(b)]
			available = cap(b) - len(b)
		}

		// 读取到可用空间
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}

func CopyStream(src io.ReadCloser) ([]byte, io.ReadCloser) {
	bodyBt, err := io.ReadAll(src)
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	return bodyBt, io.NopCloser(bytes.NewBuffer(bodyBt))
}
