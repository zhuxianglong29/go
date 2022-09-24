package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pelletier/go-toml"
)

func main() {
	config, _ := toml.LoadFile("./test.toml")
	key := config.Get("os").(string)
	addr := config.Get("server.addr").(string)
	a := fmt.Sprint(key, addr)
	fmt.Println(a)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-c
	fmt.Println("sigal return=", s)

}
