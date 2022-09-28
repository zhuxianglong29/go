package main

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	//writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755) //os.O_APPEND设置为打开文件
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	logrus.SetOutput(io.MultiWriter(writer2, writer3)) //writer1,

	logrus.Info("info msg")
	logrus.Info("info msg")
	logrus.Info("info msg")
	logrus.Println("这是一条很普通的日志。")
}
