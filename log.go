package utils

import (
	"configs"
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log *logrus.Logger

type RotateLogsHandler struct {
	logger *rotatelogs.RotateLogs
}

func NewRotateLogsHandler() (*RotateLogsHandler, error) {
	logInfo := configs.Config.Log
	logFilePath := logInfo.Path
	fmt.Println("log file path:", logFilePath)
	logFileName := "sanicalclog"
	rl, err := rotatelogs.New(
		logFilePath+"/"+logFileName+".%Y%m%d.log",
		rotatelogs.WithLinkName(logFileName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}
	return &RotateLogsHandler{logger: rl}, nil
}

func (h *RotateLogsHandler) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = h.logger.Write([]byte(line))
	return err
}

func (h *RotateLogsHandler) Levels() []logrus.Level {
	return logrus.AllLevels
}

func InitLog() {
	handler, err := NewRotateLogsHandler()
	if err != nil {
		fmt.Println("rotate logs err:", err.Error())
		panic(err)
	}

	// 将自定义日志处理器添加到日志记录器中
	log = logrus.New()
	log.SetReportCaller(true)
	log.Hooks.Add(handler)
	log.Info("log set success")
}

func GetLog() *logrus.Logger {
	return log
}
