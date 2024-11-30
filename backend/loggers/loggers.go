package loggers

import (
	"fmt"
	"log"
	"os"
)

const (
	RESET  = "\033[0m"
	RED    = "\033[41m"
	YELLOW = "\033[43m"
	CYAN   = "\033[46m"
	GREEN  = "\033[42m"
)

var (
	Warning = initLogger(fmt.Sprintf("%sWARNING:%s ", YELLOW, RESET))
	Info    = initLogger(fmt.Sprintf("%sINFO:%s ", CYAN, RESET))
	Error   = initLogger(fmt.Sprintf("%sERROR:%s ", RED, RESET))
	Debug   = initLogger(fmt.Sprintf("%sDEBUG:%s ", GREEN, RESET))
)

func initLogger(prefix string) *log.Logger {
	logOptions := log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix

	return log.New(os.Stderr, prefix, logOptions)
}
