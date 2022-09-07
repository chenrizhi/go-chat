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

	server *Server
}

// NewUser 创建用户的接口
func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name:   addr,
		Addr:   addr,
		C:      make(chan string),
		Conn:   conn,
		server: server,
	}

	// 启动监听当前用户channel消息的goroutine
	go user.ListenMessage()

	return user
}

// Online 用户上线
func (u *User) Online() {
	// 用户上线，将用户加入到OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	// 广播用户上线消息
	u.server.BroadCast(u, "已上线")
}

// Offline 用户下线
func (u *User) Offline() {
	// 用户下线，将用户从OnlineMap删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	// 广播用户下线消息
	u.server.BroadCast(u, "已下线")
}

// DoMessage 处理消息
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前用户列表
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineUsers := "[" + user.Addr + "] " + user.Name + " 在线\n"
			u.SendMessage(onlineUsers)
		}
		u.server.mapLock.Unlock()
	} else {
		u.server.BroadCast(u, msg)
	}
}

// SendMessage 发送消息
func (u *User) SendMessage(msg string) {
	_, err := u.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("user send message failed, err: ", err)
		return
	}
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
