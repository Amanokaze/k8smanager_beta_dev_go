package logger

import "time"

type RollingFileLogManager struct {
	LogManager
	rollingType        LogRollingType
	curRollingValue    int
	fileDateTimeFormat string
}

func CreateRollingFileLogManager() *RollingFileLogManager {
	rLogMgr := &RollingFileLogManager{
		LogManager:         *CreateLogManager(),
		rollingType:        rtDate,
		curRollingValue:    -1,
		fileDateTimeFormat: *getLogRollingTypeFormat(rtDate),
	}
	rLogMgr.fileNameChecker = rLogMgr

	return rLogMgr
}

func (rlog *RollingFileLogManager) SetLogRollingType(rtype LogRollingType) {
	if rlog.rollingType != rtype {
		rlog.rollingType = rtype
		rlog.fileDateTimeFormat = *getLogRollingTypeFormat(rtype)
		rlog.curRollingValue = -1
	}
}

func (rlog *RollingFileLogManager) checkFileName() {
	if rlog.rollingType == rtSingle {
		return
	}

	rVal := 0
	now := time.Now()
	if rlog.rollingType == rtDate {
		rVal = now.YearDay()
	} else if rlog.rollingType == rtHour {
		rVal = now.Hour()
	} else if rlog.rollingType == rtMinute {
		rVal = now.Minute()
	} else if rlog.rollingType == rtMonth {
		rVal = int(now.Month())
	}

	if rlog.curRollingValue != rVal {
		rlog.curRollingValue = rVal
		rlog.LogManager.logFileFullPath = rlog.filePath + "_" + now.Format(rlog.fileDateTimeFormat) + "." + rlog.fileExt
	}
}
