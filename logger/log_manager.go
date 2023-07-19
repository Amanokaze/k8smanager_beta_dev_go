package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileNameChecker interface {
	checkFileName()
}

type LogManager struct {
	logFileFullPath string
	filePath        string
	fileExt         string
	bakFileExt      string
	bakDivStr       string
	dateTimeFormat  string
	bakRollingCount int16
	maxLogMegaByte  int64
	fileNameChecker FileNameChecker
}

var IsConsoleLog = false

func CreateLogManager() *LogManager {
	fullpath := "log/onTunelog.log"

	logManager := &LogManager{
		logFileFullPath: fullpath,
		filePath:        "log/onTunelog",
		fileExt:         "log",
		bakFileExt:      "log",
		bakDivStr:       "_",
		dateTimeFormat:  "2006/01/02 15:04:05",
		bakRollingCount: 0,
		maxLogMegaByte:  0,
		fileNameChecker: nil,
	}

	return logManager
}

func (l *LogManager) SetFilePath(filePath string) {
	if l.filePath != filePath {
		l.filePath = filePath
		l.makeFullFilePath()
		err := CheckOrCreateFilePathDir(l.filePath, "ontune")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (l *LogManager) makeFullFilePath() {
	l.logFileFullPath = l.filePath + "." + l.fileExt
}

func (l *LogManager) GetFullFilePath() string {
	if l.fileNameChecker != nil {
		l.fileNameChecker.checkFileName()
	}

	return l.logFileFullPath
}

func (l *LogManager) SetFileExt(fileExt string) {
	if l.fileExt != fileExt {
		l.fileExt = fileExt
		l.makeFullFilePath()
	}
}

func (l *LogManager) SetBakFileExt(bakFileExt string) {
	l.bakFileExt = bakFileExt
}

func (l *LogManager) SetBakDivStr(bakDivStr string) {
	l.bakDivStr = bakDivStr
}

func (l *LogManager) SetBakRollingCount(bakRollingCount int16) {
	l.bakRollingCount = bakRollingCount
	if (l.bakRollingCount > 0) && (l.maxLogMegaByte == 0) {
		l.SetMaxLogMegabytes(10)
	}
}

func (l *LogManager) SetMaxLogMegabytes(maxLogMegaByte int64) {
	l.maxLogMegaByte = 1048576 * maxLogMegaByte
}

func (l *LogManager) SetLogDateFormat(dateTimeFormat string) {
	l.dateTimeFormat = dateTimeFormat
}

func (l *LogManager) getBakLogFilePath(index int16) string {
	extIndex := strings.LastIndex(l.logFileFullPath, ".")
	var bBuf bytes.Buffer
	bBuf.WriteString(l.logFileFullPath[:extIndex])
	bBuf.WriteString(l.bakDivStr)
	bBuf.WriteString(strconv.Itoa(int(index)))
	bBuf.WriteString(".")
	bBuf.WriteString(l.bakFileExt)

	return bBuf.String()
}

func (l *LogManager) setTimeLogHeader(buf *bytes.Buffer) {
	buf.WriteString("[")
	buf.WriteString(time.Now().Format(l.dateTimeFormat))
	buf.WriteString("] ")
}

func (l *LogManager) WriteLog(text string) {
	var buf bytes.Buffer
	l.setTimeLogHeader(&buf)
	buf.WriteString(text)
	buf.WriteString("\r\n")
	l.WriteBufLog(&buf)
}

func (l *LogManager) WriteBufLog(buf *bytes.Buffer) {
	if IsConsoleLog {
		fmt.Println(buf.String())
	}

	if l.fileNameChecker != nil {
		l.fileNameChecker.checkFileName()
	}

	var file, err = os.OpenFile(l.logFileFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if CheckError(err) {
		fmt.Printf("Logfile Open Error - %s %v", l.logFileFullPath, err.Error())
		panic(err.Error())
	}
	logger := *log.New(file, "", 0)

	defer l.checkBackupLog()
	defer file.Close()

	err = logger.Output(1, buf.String())
	if CheckError(err) {
		fmt.Printf("Logfile Write Error - %s %v", l.logFileFullPath, err.Error())
		return
	}
}

func (l *LogManager) checkBackupLog() {
	fs, err := os.Stat(l.logFileFullPath)
	if CheckError(err) || (fs.Size() < l.maxLogMegaByte) {
		return
	}

	if l.bakRollingCount > 0 {
		var nBakCount int16 = 1
		bakFilePath := l.getBakLogFilePath(nBakCount)
		var preModifiedtime int64

		for {
			file, err := os.Stat(bakFilePath)
			if os.IsNotExist(err) || file.IsDir() || (preModifiedtime > file.ModTime().UnixMilli()) {
				break
			} else {
				preModifiedtime = file.ModTime().UnixMilli()
				nBakCount++
				bakFilePath = l.getBakLogFilePath(nBakCount)
			}
		}

		if l.bakRollingCount <= nBakCount {
			nBakCount = 1
			bakFilePath = l.getBakLogFilePath(nBakCount)
		}
		DeleteFile(bakFilePath)
		RenameFile(l.logFileFullPath, bakFilePath)
	}
}
