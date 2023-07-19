package logger

import (
	"os"
	"strings"
	"testing"
)

func TestCreateOnTuneLogManager(t *testing.T) {
	oLogMgr := &OnTuneLogManager{
		RollingFileLogManager: *CreateRollingFileLogManager(),
	}
	oLogMgr.fileNameChecker = oLogMgr

	err := CheckOrCreateFilePathDir(oLogMgr.RollingFileLogManager.LogManager.GetFullFilePath(), "")
	if err != nil {
		t.Errorf("Directory Creation Error: %s", err.Error())
	}

	oLogMgr.WriteLog("Log Write Test")

	err = RemoveDir(oLogMgr.RollingFileLogManager.LogManager.GetFullFilePath())
	if err != nil {
		t.Errorf("Directory Remove Error: %s", err.Error())
	}
}

func RemoveDir(filePath string) error {
	tmp := strings.Replace(filePath, "\\", "/", -1)
	//tmp := filePath
	nPos := 0
	for {
		pos := strings.Index(tmp, "/")
		if pos > 0 {
			dir := filePath[:pos+nPos]
			err := os.RemoveAll(dir)
			if err != nil {
				return err
			}

			tmp = tmp[pos+1:]
			nPos += pos + 1
		} else if pos == 0 {
			tmp = tmp[pos+1:]
			nPos += pos + 1
		} else {
			break
		}
	}

	return nil
}
