package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	fileCount := map[string]int{
		"cpp": 10,
		"js":  8,
	}
	bytes, _ := json.Marshal(fileCount)
	fmt.Println(string(bytes))

	var Window map[string]int

	err := json.Unmarshal([]byte(bytes), &Window)
	if err != nil {
		fmt.Println("JSON decode error!")
		return
	}
	fmt.Println(Window)

}
