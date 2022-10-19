package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}

	defer c.Close()
	_, err = c.Do("expire", "abc", 10) //加过期时间
	if err != nil {
		fmt.Println(err)
		return
	}
	//time.Sleep(time.Second * 1)
	// r, err := redis.Int(c.Do("Get", "abc"))
	// if err != nil {
	// 	fmt.Println("get abc failed,", err)
	// 	return
	// }
	// fmt.Println("get abc=", r)
}
