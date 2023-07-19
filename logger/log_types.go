package logger

type LogType uint8
type LogRollingType uint8

const (
	ltInfo LogType = iota
	ltConfig
	ltWarn
	ltError
	ltDebug
)

var logTypeTexts = [5]string{"[Info] ", "[Config] ", "[Warnnig] ", "[Error] ", "[Debug] "}

func getLogTypeText(lt LogType) *string {
	return &logTypeTexts[lt]
}

const (
	rtMonth LogRollingType = iota
	rtDate
	rtHour
	rtMinute
	rtSingle
)

const (
	yyMM        = "0601"
	yyMMdd      = yyMM + "02"
	yyMMddHH    = yyMMdd + "15"
	yyMMdd_HHmm = yyMMdd + "_" + "1504"
)

var logRollingTypeFormats = [5]string{yyMM, yyMMdd, yyMMddHH, yyMMdd_HHmm, ""}

func getLogRollingTypeFormat(lt LogRollingType) *string {
	return &logRollingTypeFormats[lt]
}
