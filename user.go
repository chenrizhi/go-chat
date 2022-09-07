package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

// NewUser 创建用户的接口
func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name: addr,
		Addr: addr,
		C:    make(chan string),
		Conn: conn,
	}

	// 启动监听当前用户channel消息的goroutine
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前User channel的方法，一旦有消息，就发送给对应客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, err := u.Conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("send message failed, err: ", err)
		}
	}
}
