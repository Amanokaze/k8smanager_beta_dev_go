package logger

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getRollingTypeFormatFileName(otlog *OnTuneLogManager, rtype LogRollingType) string {
	return otlog.filePath + "_" + time.Now().Format(*getLogRollingTypeFormat(rtype)) + "." + otlog.fileExt
}

func getLogTypeHeaderText(otlog *OnTuneLogManager, header *string, text string, ltype LogType) string {
	*header = *getLogTypeText(ltype)
	return strings.Trim(string(ReadFileLastBytes(otlog.logFileFullPath, int64(len(*header)+len(text)+2))), "\r\n")
}

func assertLogRollingtypeFiles(otlog *OnTuneLogManager, assert *assert.Assertions, rtype LogRollingType) {
	otlog.SetMaxLogMegabytes(1)
	otlog.SetLogRollingType(rtype)
	filename := getRollingTypeFormatFileName(otlog, rtype)
	otlog.WriteLog("111")
	assert.Equal(IsExistFile(filename), true)

	otlog.maxLogMegaByte = 1
	for iCnt := 1; iCnt < 6; iCnt++ {
		otlog.WriteLog(strconv.Itoa(iCnt * 111))
		exist := (iCnt < int(otlog.bakRollingCount))
		filename = otlog.getBakLogFilePath(int16(iCnt))
		assert.Equal(IsExistFile(filename), exist)
	}
}

func TestLog_manager(t *testing.T) {
	IsConsoleLog = false

	assert := assert.New(t)

	//Test Properies
	logMgr := CreateOnTuneLogManager(false)
	logMgr.SetLogRollingType(rtSingle)
	logMgr.SetFilePath("log/onTuneKube")
	logMgr.SetFileExt("txt")
	logMgr.SetBakDivStr("^")
	logMgr.SetBakFileExt("bak")
	logMgr.SetMaxLogMegabytes(10)
	logMgr.SetBakRollingCount(5)
	logMgr.SetLogDateFormat("2006-01-02 15-04-05")

	assert.Equal(logMgr.logFileFullPath, "log/onTuneKube.txt")
	assert.Equal(logMgr.getBakLogFilePath(3), "log/onTuneKube^3.bak")
	assert.Equal(logMgr.maxLogMegaByte, int64(1048576*10))
	assert.Equal(logMgr.bakRollingCount, int16(5))
	assert.Equal(IsExistFile("log"), true)
	DeleteFile(logMgr.logFileFullPath)

	//Test CreateFile and FileExist
	logMgr.WriteLog("111")
	assert.Equal(IsExistFile("log/onTuneKube.txt"), true)

	//Test Bakup CreateFile and FileExist
	logMgr.maxLogMegaByte = 1
	for iCnt := 1; iCnt < 6; iCnt++ {
		logMgr.WriteLog(strconv.Itoa(iCnt * 111))
		exist := (iCnt < int(logMgr.bakRollingCount))
		assert.Equal(IsExistFile("log/onTuneKube^"+strconv.Itoa(iCnt)+".bak"), exist)
	}

	assert.Equal(IsExistFile("log/onTuneKube^5.bak"), false)

	//Test onTunelog_manager write contets > Info(), Warn(), Config(), Error(), Debug()
	logMgr.SetMaxLogMegabytes(1)

	var header string
	text := "testInfo"
	logMgr.Info(text)

	lText := getLogTypeHeaderText(logMgr, &header, text, ltInfo)
	assert.Equal(lText, header+text)

	text = "testWarn"
	logMgr.Warn(text)
	lText = getLogTypeHeaderText(logMgr, &header, text, ltWarn)
	assert.Equal(lText, header+text)

	text = "testConfig"
	logMgr.Config(text)
	lText = getLogTypeHeaderText(logMgr, &header, text, ltConfig)
	assert.Equal(lText, header+text)

	text = "testError"
	logMgr.Error(text)
	lText = getLogTypeHeaderText(logMgr, &header, text, ltError)
	assert.Equal(lText, header+text)

	logMgr.DebugLog = false
	logMgr.Debug("testDebug")
	assert.Equal(lText, header+text)

	logMgr.DebugLog = true
	text = "testDebug"
	logMgr.Debug(text)
	lText = getLogTypeHeaderText(logMgr, &header, text, ltDebug)
	assert.Equal(lText, header+text)

	assert.Equal(IsExistFile("log/onTuneKube.txt"), true)

	//Test createFile at RollingLogType > rtDate
	assertLogRollingtypeFiles(logMgr, assert, rtDate)

	//Test createFile at RollingLogType > rtMinute
	assertLogRollingtypeFiles(logMgr, assert, rtMinute)

	//Test createFile at RollingLogType > rtMonth
	assertLogRollingtypeFiles(logMgr, assert, rtMonth)
}
