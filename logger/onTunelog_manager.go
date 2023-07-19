package logger

import (
	"bytes"
	"log"
)

type OnTuneLogManager struct {
	RollingFileLogManager
	DebugLog bool
}

func CreateOnTuneLogManager(debug bool) *OnTuneLogManager {
	oLogMgr := &OnTuneLogManager{
		RollingFileLogManager: *CreateRollingFileLogManager(),
		DebugLog:              debug,
	}
	oLogMgr.fileNameChecker = oLogMgr

	err := CheckOrCreateFilePathDir(oLogMgr.RollingFileLogManager.LogManager.GetFullFilePath(), "")
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return oLogMgr
}

func (otlog *OnTuneLogManager) writeExpLog(text *string, ltype LogType) {
	var buf bytes.Buffer
	otlog.setTimeLogHeader(&buf)
	buf.WriteString(*getLogTypeText(ltype))
	buf.WriteString(*text)
	buf.WriteString("\r\n")
	otlog.WriteBufLog(&buf)
}

func (otlog *OnTuneLogManager) Info(text string) {
	otlog.writeExpLog(&text, ltInfo)
}

func (otlog *OnTuneLogManager) Config(text string) {
	otlog.writeExpLog(&text, ltConfig)
}

func (otlog *OnTuneLogManager) Warn(text string) {
	otlog.writeExpLog(&text, ltWarn)
}

func (otlog *OnTuneLogManager) Error(text string) {
	otlog.writeExpLog(&text, ltError)
}

func (otlog *OnTuneLogManager) Debug(text string) {
	if otlog.DebugLog {
		otlog.writeExpLog(&text, ltDebug)
	}
}
