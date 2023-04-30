package appContextConfig

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// logEntry
// --------------------------------------------------------------------------------------------
func LogEntry(logMsg string, runtimeSkip int) string {
	pc, file, line, ok := runtime.Caller(runtimeSkip)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return fmt.Sprintf(" file: %s, function: %s, msg: %s", filename, fn, logMsg)
}

// logInfo
// --------------------------------------------------------------------------------------------
func (app *Application) LogInfo(msg string) {
	app.SysLog.Info(LogEntry(msg, 2))
}
