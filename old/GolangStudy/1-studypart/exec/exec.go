package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// //执行命令并返回结果
	// res, err := exec.Command("ls", "-l").Output() /*ls -l /users...(执行路径，空为当前路径)//Output 和 CombinedOutput两种返回方法（func (c *Cmd) Output() ([]byte, error)
	// //执行Cmd中包含的命令，并返回标准输出的切片
	// func (c *Cmd) CombinedOutput() ([]byte, error)
	// //执行Cmd中包含的命令，并返回标准输出与标准错误合并后的切片）*/
	// if err != nil {
	// 	fmt.Println("执行出错：", err)
	// }
	// fmt.Println(string(res))

	//区分stdin、stderr
	// cmd := exec.Command("ls", "-l", "/var/log/*.log")
	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout // 标准输出
	// cmd.Stderr = &stderr // 标准错误
	// err := cmd.Run()
	// outStr, errStr := stdout.String(), stderr.String()
	// fmt.Printf("标准输出:\n%s\n标准错误:\n%s\n", outStr, errStr)
	// if err != nil {
	// 	fmt.Println("执行出错：", err)
	// }

	//多条命令，要使用管道
	// c1 := exec.Command("grep", "ERROR", "/var/log/messages")
	// c2 := exec.Command("wc", "-l")
	// c2.Stdin, _ = c1.StdoutPipe()
	// c2.Stdout = os.Stdout
	// _ = c2.Start()
	// _ = c1.Run()
	// _ = c2.Wait()

	//设置命令级别的环境变量
	os.Setenv("NAME", "wangbm")
	cmd := exec.Command("echo", os.ExpandEnv("$NAME"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("%s", out)

}
