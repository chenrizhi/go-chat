package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个Server的接口
func NewServer(ip string, port int) *Server {
	s := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return s
}

// BroadCast 广播消息
func (s Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// 处理业务
	fmt.Println("connection success!")
	user := NewUser(conn)
	// 用户上线，将用户加入到OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	// 广播用户上线消息
	s.BroadCast(user, "已上线")

	// 接受用户端发送的消息
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)

			if n == 0 {
				s.BroadCast(user, "已下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read failed, err: ", err)
				return
			}

			// 提取用户消息，去除\n
			msg := string(buf)[:n-1]
			s.BroadCast(user, msg)
		}
	}()
}

// ListenMessage 监听Message广播消息channel的goroutine，一旦有消息，就发送给全部的在线用户
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		// 将msg发送给全部的在线用户
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// Start 启动服务器的接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("server listen failed, err: ")
		return
	}
	defer listener.Close()

	// 启动监听Message的goroutine
	go s.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed, err: ", err)
			return
		}

		// do handler
		go s.Handler(conn)
	}
}
