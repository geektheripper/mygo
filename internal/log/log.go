package log

import (
	"os"

	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stderr)

func GetLogger() *log.Logger {
	return logger
}
