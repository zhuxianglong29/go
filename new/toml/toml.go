package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type (
	config struct {
		Os      string
		Version string
		Server  server
	}
	server struct {
		Addr  string
		Port  string
		Http  string
		Route string
	}
)

func main() {
	var tomlconfig config
	filePath := "test.toml"
	if _, err := toml.DecodeFile(filePath, &tomlconfig); err != nil {
		panic(err)
	}
	fmt.Println(tomlconfig)

}
