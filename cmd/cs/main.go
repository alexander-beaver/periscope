package main

import (
	"CloudScan/pkg/manager"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	//Get file path as the first CLI argument
	configFile, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error reading config file")
		os.Exit(1)
	}
	c := &manager.Config{}
	err = json.Unmarshal(configFile, &c)
	if err != nil {
		return
	}
	manager.Run(*c)
}
