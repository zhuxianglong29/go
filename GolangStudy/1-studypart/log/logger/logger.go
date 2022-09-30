package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("./xx.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetOutput(logFile)
	log.Println("这是一条很普通的日志。")
	v := "Printf输出的"
	log.Printf("这是一条%s日志。\n", v)
	//log.Fatalln("这是一条会触发fatal的日志。") //fatal会直接跳过出程序，与panic不同，不会执行defer,而panic会先执行defer
	//log.Panicln("这是一条会触发panic的日志。")
	log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)
	log.Println("这是设置后的日志")
}
