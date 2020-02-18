package log

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

// Фильтр логирования
var LogFilter = &logutils.LevelFilter{
	Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
	MinLevel: logutils.LogLevel("ERROR"),
	Writer:   os.Stderr,
}

//PrintfDebug print message only if level = "DEBUG"
func PrintfDebug(format string, v ...interface{}) {
	if LogFilter.MinLevel == "DEBUG" {
		log.Printf(format, v...)
	}
}
