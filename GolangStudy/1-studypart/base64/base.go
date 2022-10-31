package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	input := []byte("1")

	// 演示base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	fmt.Println(encodeString)

	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(decodeBytes))

	fmt.Println("///////////////////////////////")
	//var str string
	//str = "YCh7vACAGAAC8IaPQk0Q8Gufor95zXlimzz4+24="
	str := "gCh7vACAHwACFwxSv52724I=" //"A+gH0A=="
	fmt.Println(str)
	data := []byte(str)

	var strdata map[string]int

	err := json.Unmarshal(data, &strdata)
	//fmt.Println(err)

	destr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("edstr", string(destr))

	// // 如果要用在url中，需要使用URLEncoding
	// uEnc := base64.URLEncoding.EncodeToString([]byte(input))
	// fmt.Println(uEnc)

	// uDec, err := base64.URLEncoding.DecodeString(uEnc)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(string(uDec))
}
