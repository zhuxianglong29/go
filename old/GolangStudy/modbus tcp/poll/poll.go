package main

import (
	"fmt"
	"time"

	modbus "github.com/thinkgos/gomodbus/v2"
)

var datemodbus []uint16

// modbus
func mbus() {
	p := modbus.NewTCPClientProvider("192.168.1.48:502", modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("connect failed, ", err)
		return
	}
	defer client.Close()

	fmt.Println("starting")
	for i := 0; i < 5; i++ {
		a, err := client.ReadHoldingRegisters(1, 0, 2)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(a)
		datemodbus = a
		//	fmt.Printf("ReadDiscreteInputs %#v\r\n", results)

		time.Sleep(time.Second)
	}
}

func main() {
	go mbus()
	for {

	}
}
