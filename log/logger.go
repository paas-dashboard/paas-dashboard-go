package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var Layout = "2006-01-02T15:04:05,000Z07:00"

type Formatter struct {
	TimestampFormat string
}

func InitLogger() {
	// lumberjack.Logger是满足io.Writer接口的，可以被用作logrus的输出
	logrus.StandardLogger().SetOutput(&lumberjack.Logger{
		Filename:   "/opt/paas-dashboard/logs/checker.log", // 日志文件路径
		MaxSize:    100,                                    // megabytes
		MaxBackups: 10,                                     // 最多保留3个备份
		MaxAge:     28,                                     // days
		Compress:   true,                                   // 是否压缩备份文件

	})
	logrus.StandardLogger().SetLevel(logrus.InfoLevel)
	level := os.Getenv("LOGLEVEL")
	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Errorf("parse level err: %s", level)
	} else {
		logrus.StandardLogger().SetLevel(l)
	}
	logrus.SetFormatter(&Formatter{})
	logrus.Info("init log done.")
}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var output bytes.Buffer

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = Layout
	}

	output.WriteString(entry.Time.Format(timestampFormat))
	output.WriteString("|")
	output.WriteString(strings.ToUpper(entry.Level.String()))
	output.WriteString("|")

	for _, v := range entry.Data {
		output.WriteString(fmt.Sprintf("%v|", v))
	}

	output.WriteString(entry.Message)
	output.WriteRune('\n')

	return output.Bytes(), nil
}
