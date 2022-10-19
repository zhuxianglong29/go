package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	// //底层文件方式操作,中端输入输出
	// var buf [16]byte
	// os.Stdin.Read(buf[:])
	// os.Stdin.WriteString(string(buf[:]))

	// 只读方式打开当前目录下的main.go文件
	file, err := os.Create("./test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 5; i++ {
		file.WriteString("ab\n") //创建时才能写
		file.Write([]byte("cd\n"))
	}
	file.Close()

	file1, err := os.Open("./test.txt")
	if err != nil {
		fmt.Println("open file failed!, err:", err)
		return
	}
	fmt.Println("open success")

	var buf [128]byte
	var content []byte
	for {
		n, err := file1.Read(buf[:])
		if err == io.EOF {
			// 读取结束
			break
		}
		if err != nil {
			fmt.Println("read file err ", err)
			return
		}
		content = append(content, buf[:n]...)
	}
	fmt.Println(string(content))
	time.Sleep(time.Second * 2)

	file1.Close()
}
