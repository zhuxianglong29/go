package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8888/chat", nil)
	if err != nil {
		fmt.Println("错误信息:", err)
	}
	wg.Add(2)
	go read(conn)
	go writeM(conn)
	wg.Wait()
}
func read(conn *websocket.Conn) {
	defer wg.Done()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("错误信息:", err)
			break
		}
		if err == io.EOF {
			continue
		}
		fmt.Println("获取到的信息:", string(msg))
	}
}
func writeM(conn *websocket.Conn) {
	defer wg.Done()
	for {
		fmt.Print("请输入:")
		reader := bufio.NewReader(os.Stdin)
		data, _ := reader.ReadString('\n')
		conn.WriteMessage(1, []byte(data))
	}
}
