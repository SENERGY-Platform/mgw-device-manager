package util

import (
	sb_util "github.com/SENERGY-Platform/go-service-base/util"
	log_level "github.com/y-du/go-log-level"
	"os"
)

var Logger *log_level.Logger

func InitLogger(config sb_util.LoggerConfig) (out *os.File, err error) {
	Logger, out, err = sb_util.NewLogger(config)
	Logger.SetLevelPrefix("ERROR ", "WARNING ", "INFO ", "DEBUG ")
	return
}
