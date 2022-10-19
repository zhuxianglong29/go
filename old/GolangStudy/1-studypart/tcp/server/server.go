package main

import (
	"fmt"
	"net"
	"strconv"
	_ "strconv"
	_ "strings"
	"time"
)

func process(conn net.Conn) {
	//循环接收
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		fmt.Printf("服务器在等待客户%s 发送信息\n", conn.RemoteAddr().String())
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("客户端退出err=%v", err)
			return
		}
		//3.显示客户端发送内容
		fmt.Print(string(buf[:n-1]))
		fmt.Printf("buf格式%T %v\n", buf, buf[:n])
		//**数据转换并处理
		strnub := string(buf[:n-2])
		fmt.Println("strnub=", strnub)
		k, _ := strconv.Atoi(strnub)
		fmt.Println("k=", k)
		j := strconv.Itoa(jisheng(k))

		// //强转
		// num int8 := int8(buf[:n-1])
		// j:=jisheng(num)

		nub := []byte(j)
		fmt.Println("阶乘\n", jisheng(k))
		conn.Write(nub)
		time.Sleep(time.Second)
		//a := []byte("\n格式无错")
		// b:=buf[0:n-1]
		// c:=append(nub,,b)

		fmt.Printf("j%T %v", j, j)
	}
}

func jisheng(n int) (result int) {
	if n > 1 {
		result = n * jisheng(n-1)
		return result
	}
	return 1
}

func main() {
	fmt.Println("服务器开始监")
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen err= ", err)
		return
	}
	defer listen.Close()

	//等待链接
	for {
		fmt.Println("等待连接...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Acceept err=", err)
		} else {
			fmt.Printf("Accept suc conn=%v 客户端ip=%v\n", conn, conn.RemoteAddr().String())
		}
		//协程-客户�?
		go process(conn)

	}

}
