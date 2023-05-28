package log

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	// lumberjack.Logger是满足io.Writer接口的，可以被用作logrus的输出
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/checker/checker.log", // 日志文件路径
		MaxSize:    500,                            // megabytes
		MaxBackups: 3,                              // 最多保留3个备份
		MaxAge:     28,                             // days
		Compress:   true,                           // 是否压缩备份文件
	})
}
