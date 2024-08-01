package util

import (
	sb_logger "github.com/SENERGY-Platform/go-service-base/logger"
	"github.com/y-du/go-log-level"
	"os"
)

var Logger *log_level.Logger

func InitLogger(c LoggerConfig) (out *os.File, err error) {
	Logger, out, err = sb_logger.New(c.Level, c.Path, c.FileName, c.Prefix, c.Utc, c.Terminal, c.Microseconds)
	Logger.SetLevelPrefix("ERROR ", "WARNING ", "INFO ", "DEBUG ")
	return
}
