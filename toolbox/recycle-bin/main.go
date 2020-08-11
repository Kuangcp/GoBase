package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wonderivan/logger"
)

var mainDir = ".config/app-conf/recycle-bin"
var trashDir = mainDir + "/trash"
var configDir = mainDir + "/conf"
var dbDir = configDir + "/db"

var info = cuibase.HelpInfo{
	Description: "Recycle bin",
	Version:     "1.0.0",
	VerbLen:     -3,
	ParamLen:    -10,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "Help info",
		}, {
			Verb:    "",
			Param:   "file",
			Comment: "Delete file",
		}, {
			Verb:    "-as",
			Param:   "suffix",
			Comment: "delete *.suffix",
		},
	}}

var (
	help   bool
	suffix string
)

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.StringVar(&suffix, "as", "", "")
}

func main() {
	flag.Parse()

	if help {
		info.PrintHelp()
		return
	}

	logger.SetLogPathTrim("recycle-bin")
	home, err := cuibase.Home()
	cuibase.CheckIfError(err)
	trashDir = home + "/" + trashDir

	if suffix != "" {
		DeleteFile(strings.Split(suffix, ","))
		return
	}

	DeleteFile(os.Args)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteFile(params []string) {
	var filepath string
	length := len(params)
	if length < 2 {
		return
	}

	if length == 2 {
		filepath = params[1]
	} else {
		filepath = params[2]
	}

	deleteFiles(filepath)
}

// deleteFies 移动文件到回收站，并通过LevelDB存储时间等信息
func deleteFiles(files ...string) {
	db, err := leveldb.OpenFile(dbDir, nil)
	defer db.Close()
	if err != nil {
		logger.Error(err)
		return
	}

	// TODO 多个同名文件处理方式
	for _, filepath := range files {
		logger.Info(filepath, "move to", trashDir)

		exists, err := PathExists(filepath)
		cuibase.CheckIfError(err)
		if !exists {
			logger.Error(filepath, "not found")
			return
		}

		cmd := exec.Command("mv", filepath, trashDir+"/"+filepath)
		var out bytes.Buffer

		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			logger.Error(err)
			return
		}

		db.Put([]byte(filepath), []byte(""), nil)
	}
}

func DeleteFileBySuffix(params []string) {
	fmt.Println(params)
}
