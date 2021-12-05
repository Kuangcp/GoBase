package ghelp

import (
	"embed"
	"io/fs"
	"path"
)

type StaticResource struct {
	// 静态资源
	StaticFS embed.FS
	// 设置embed文件到静态资源的相对路径，也就是embed注释里的路径
	Path string
}

// Open 静态资源被访问逻辑
func (_this *StaticResource) Open(name string) (fs.File, error) {
	fullName := path.Join(_this.Path, name)
	return _this.StaticFS.Open(fullName)
}
