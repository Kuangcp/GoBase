package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
)

const (
	use = ".use"
	not = ".not"

	titleMaxLen   = 20
	contentMaxLen = 3000
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
	success, err := switchFileState(file)
	if success {
		useState, err := fileUseState(file)
		if err != nil {
			logger.Warn("switch failed", err)
			ghelp.GinFailedWithMsg(c, err.Error())
			return
		}
		ghelp.GinSuccessWith(c, useState)
		return
	}

	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}
	ghelp.GinFailedWithMsg(c, "file not exist")
	return
}

func switchFileState(fileName string) (bool, error) {
	if fileName == "" {
		return false, nil
	}
	logger.Info("switch:", fileName)

	for _, s := range stateList {
		filePath := groupDir + fileName
		useState := filePath + s
		exists, err := isPathExists(useState)
		if err != nil {
			return false, err
		}

		if exists {
			// 当前为 not 才表示启用
			_, _, err := switchState(filePath, s == not)
			if err != nil {
				return false, err
			}

			err = generateHost()
			if err != nil {
				// rollback
				_, _, _ = switchState(filePath, s != not)
				return false, err
			} else {
				updateFileItemState(FileItemVO{Name: fileName, Use: s == not})
				return true, nil
			}
		}
	}
	return false, nil
}

func CurrentHosts(c *gin.Context) {
	exists, err := isPathExists(curHostFile)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}
	if !exists {
		return
	}
	readFile, err := ioutil.ReadFile(curHostFile)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	content := string(readFile)
	ghelp.GinSuccessWith(c, content)
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

func fileUseState(file string) (bool, error) {
	for _, s := range stateList {
		filePath := groupDir + file + s
		exists, err := isPathExists(filePath)
		if exists {
			return s == use, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, fmt.Errorf("file not exist")
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

func QueryMode(c *gin.Context) {
	switch SupportMode {
	case "":
		fallthrough
	case hostType:
		ghelp.GinSuccessWith(c, "text/x-hosts")
	case nginxType:
		ghelp.GinSuccessWith(c, "text/x-nginx-conf")
	}
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

	if param.Name == "" || param.Content == "" ||
		utf8.RuneCountInString(param.Name) > titleMaxLen || utf8.RuneCountInString(param.Content) > contentMaxLen {
		ghelp.GinFailedWithMsg(c, "invalid param")
		return
	}
	param.Name = strings.TrimSpace(param.Name)

	targetFilePath, hasSwitch, err := switchState(groupDir+param.Name, param.Use)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	finalContent := normalizeContent(param.Content)
	err = ioutil.WriteFile(targetFilePath, []byte(finalContent), 0644)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	if !hasSwitch {
		addFileItem(param, nil)
	} else {
		updateFileItemState(param)
	}

	if param.Use || hasSwitch {
		err := generateHost()
		if err != nil {
			ghelp.GinFailedWithMsg(c, err.Error())
			return
		}
	}

	ghelp.GinSuccessWith(c, "")
}

func normalizeContent(content string) string {
	result := strings.TrimSpace(content)
	result = strings.ReplaceAll(content, " ", " ")
	return result
}

// return final file path that contain state
func switchState(absPath string, targetUse bool) (string, bool, error) {
	origin := absPath + use
	target := absPath + not
	if targetUse {
		origin = absPath + not
		target = absPath + use
	}

	exists, err := isPathExists(origin)
	if !exists {
		already, _ := isPathExists(target)
		return target, already, nil
	}
	if err != nil {
		logger.Error(err)
		return "", false, err
	}

	err = os.Rename(origin, target)
	if err != nil {
		logger.Error(err)
		return "", false, err
	}
	return target, true, nil
}

func generateHost() error {
	logger.Info("Generate new:", curHostFile)
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

		mergeResult += buildFileBlock(vo.Name, string(readFile))
	}

	err := ioutil.WriteFile(curHostFile, []byte(mergeResult), 0644)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf(err.Error())
	}

	doAfterCmd()
	return nil
}

func doAfterCmd() {
	if GenerateAfterCmd == "" {
		return
	}
	fmt.Println(GenerateAfterCmd)
	//execCmdWithQuite(exec.Command("nginx", "-s", "reload"))
	fields := strings.Fields(GenerateAfterCmd)
	if len(fields) > 1 {
		execCmdWithQuite(exec.Command(fields[0], fields[1:]...))
	} else {
		execCmdWithQuite(exec.Command(fields[0]))
	}
}

// 静默执行 不关心返回值
func execCmdWithQuite(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		logger.Error(cmd, err)
		os.Exit(1)
	}
}

func buildFileBlock(name, content string) string {
	nameLen := utf8.RuneCountInString(name)
	padding := (titleMaxLen - nameLen) / 2
	paddingStr := strconv.Itoa(padding)
	return "" +
		"#━━━━━━━━━━━━━━━━━━━━\n" +
		"#" + fmt.Sprintf("%"+paddingStr+"s%s%"+strconv.Itoa(titleMaxLen-padding-nameLen)+"s", "", name, "") + "\n" +
		//"#━━━━━━━━━━━━━━━━━━━━\n" +
		content + "\n" +
		"\n"
}
