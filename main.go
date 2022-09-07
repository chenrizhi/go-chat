package main

import "fmt"

func main() {
	if err := loadConfig("./go-chat.yaml"); err != nil {
		fmt.Println("config init failed, err: ", err)
		return
	}
	bind := configData.Server.Bind
	port := configData.Server.Port
	server := NewServer(bind, port)
	server.Start()
}
