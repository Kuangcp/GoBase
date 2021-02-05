package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
)

const (
	use = ".use"
	not = ".not"
)

var stateList = []string{use, not}

type (
	FileItemVO struct {
		Name    string `json:"name"`
		Use     bool   `json:"use"`
		Content string `json:"content,omitempty"`
	}
)

func SwitchFileState(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		ghelp.GinFailed(c)
		return
	}

	for _, s := range stateList {
		filePath := groupDir + file
		useState := filePath + s
		exists, err := isPathExists(useState)
		if err != nil {
			ghelp.GinFailedWithMsg(c, err.Error())
			return
		}

		if exists {
			// 当前为 not 才表示启用
			_, err := switchState(filePath, s == not)
			if err != nil {
				ghelp.GinFailedWithMsg(c, err.Error())
				return
			}

			err = generateHost()
			if err != nil {
				ghelp.GinFailedWithMsg(c, err.Error())
				return
			} else {
				ghelp.GinSuccessWith(c, "")
				return
			}
		}
	}

	ghelp.GinFailedWithMsg(c, "file not exist")
	return
}

func FileContent(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		ghelp.GinFailed(c)
		return
	}
	for _, s := range stateList {
		result := fillContentResult(c, file, s)
		if result {
			return
		}
	}
	ghelp.GinFailedWithMsg(c, "file not exist")
}

func fillContentResult(c *gin.Context, file string, state string) bool {
	filePath := groupDir + file + state
	exists, err := isPathExists(filePath)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return true
	}
	if !exists {
		return false
	}
	readFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return true
	}

	content := string(readFile)
	ghelp.GinSuccessWith(c, FileItemVO{Name: file, Content: content, Use: state == use})
	return true
}

func ListFile(c *gin.Context) {
	result := getFileList()
	ghelp.GinSuccessWith(c, result)
}

func getFileList() []FileItemVO {
	files, _ := ioutil.ReadDir(groupDir)

	var result []FileItemVO
	for _, f := range files {
		fileName := f.Name()
		if !strings.HasSuffix(fileName, use) && !strings.HasSuffix(fileName, not) {
			continue
		}
		//fmt.Println(fileName)
		result = append(result, FileItemVO{
			Name: fileName[:len(fileName)-4],
			Use:  fileName[len(fileName)-4:] == use,
		})
	}
	return result
}

func CreateOrUpdateFile(c *gin.Context) {
	var param FileItemVO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	if param.Name == "" || param.Content == "" || len(param.Name) > 30 || len(param.Content) > 3000 {
		ghelp.GinFailedWithMsg(c, "invalid param")
		return
	}

	targetFilePath, err := switchState(groupDir+param.Name, param.Use)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	err = ioutil.WriteFile(targetFilePath, []byte(param.Content), 0644)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	if param.Use {
		err := generateHost()
		if err != nil {
			ghelp.GinFailedWithMsg(c, err.Error())
			return
		}
	}

	ghelp.GinSuccessWith(c, "")
}

// return final file path that contain state
func switchState(absPath string, targetUse bool) (string, error) {
	origin := absPath + use
	target := absPath + not
	if targetUse {
		origin = absPath + not
		target = absPath + use
	}

	exists, err := isPathExists(origin)
	if !exists {
		return target, nil
	}
	if err != nil {
		logger.Error(err)
		return "", err
	}

	err = os.Rename(origin, target)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return target, nil
}

func generateHost() error {
	logger.Info("start reload host")
	list := getFileList()
	mergeResult := ""
	for _, vo := range list {
		if !vo.Use {
			continue
		}

		readFile, err := ioutil.ReadFile(groupDir + vo.Name + use)
		if err != nil {
			logger.Warn(err.Error())
			continue
		}

		mergeResult += "###########\n" +
			"#  " + vo.Name + "\n" +
			"###########\n" +
			string(readFile) + "\n\n\n"
	}

	err := ioutil.WriteFile(curHostFile, []byte(mergeResult), 0644)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf(err.Error())
	}
	return nil
}
