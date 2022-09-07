package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个Server的接口
func NewServer(ip string, port int) *Server {
	s := &Server{
		Ip:   ip,
		Port: port,
	}
	return s
}

func (s *Server) Handler(conn net.Conn) {
	// 处理业务
	fmt.Println("connection success!")
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
