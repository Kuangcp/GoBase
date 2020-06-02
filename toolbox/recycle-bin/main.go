package main

import (
	"fmt"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
	"os"
)

var trashDir = ".config/app-conf/recycle-bin"

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
			Handler: DeleteFile,
		}, {
			Verb:    "-as",
			Param:   "suffix",
			Comment: "delete *.suffix",
			Handler: DeleteFileBySuffix,
		},
	}}

func main() {
	logger.SetLogPathTrim("recycle-bin")
	home, err := cuibase.Home()
	cuibase.CheckIfError(err)
	trashDir = home + "/" + trashDir

	cuibase.RunActionFromInfo(info, DeleteFile)
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
	logger.Info(filepath, "move to", trashDir)

	exists, err := PathExists(filepath)
	cuibase.CheckIfError(err)
	if !exists {
		logger.Error(filepath, "not found")
		return
	}

	err = os.Rename(filepath, trashDir+"/"+filepath)
	if err != nil {
		logger.Error(err)
	}
}

func DeleteFileBySuffix(params []string) {
	fmt.Println(params)
}
