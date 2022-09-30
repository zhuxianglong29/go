package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

// MyHook ...
type MyHook struct {
}

// Levels 只定义 error 和 panic 等级的日志,其他日志等级不会触发 hook
func (h *MyHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.PanicLevel,
	}
}

// Fire 将异常日志写入到指定日志文件中
func (h *MyHook) Fire(entry *logrus.Entry) error {
	f, err := os.OpenFile("err.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(entry.Message)); err != nil {
		return err
	}
	return nil
}

func main() {
	/*
		logrus.SetLevel(logrus.TraceLevel)
		logrus.Trace("Trace msg")
		logrus.Debug("debug msg")
		logrus.Info("info msg")
		logrus.Warn("warn msg")
		logrus.Error("error msg")
		logrus.Fatal("fatal msg")
		logrus.Panic("panic msg")
	*/

	// logrus.SetReportCaller(true) //设置在输出日志中添加文件名和方法信息
	// logrus.Info("info msg")

	// //WithFields在日志中添加字段
	// logrus.WithFields(logrus.Fields{
	// 	"name": "dj",
	// 	"age":  18,
	// }).Info("info msg")
	// //批量添加用返回值的方式
	// a := logrus.WithFields(logrus.Fields{
	// 	"second": "abc",
	// 	"first":  123,
	// })
	// a.Info("info msg")
	// a.Error("error msg")

	// //json/text格式设置
	// logrus.SetLevel(logrus.TraceLevel)
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// logrus.Trace("trace msg")
	// logrus.Debug("debug msg")
	// logrus.Info("info msg")
	// logrus.Warn("warn msg")
	// logrus.Error("error msg")
	// logrus.Fatal("fatal msg")
	// logrus.Panic("panic msg")

	// //自设时间
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	ForceQuote:      true,                  //键值对加引号
	// 	TimestampFormat: "2006-01-02 15:04:05", //时间格式
	// 	FullTimestamp:   true,
	// })
	// logrus.WithField("name", "ball").WithField("say", "hi").Info("info log")

	//hook:通过设置，能独立处理hook设置的level,例如这里将error\panic用hook标定写入到err.log文件中
	logrus.AddHook(&MyHook{})
	logrus.Warn("Warn msg")
	logrus.Error("some errors\n")
	logrus.Panic("some panic\n")
	logrus.Print("hello world\n")

}
