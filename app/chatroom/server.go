package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//关播消息Channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// message 监听
func (s *Server) listenMessager() {
	for {
		msg := <-s.Message

		s.mapLock.Lock()
		for _, client := range s.OnlineMap {
			client.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) broadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)
	s.Message <- sendMsg

}

func (s *Server) handler(conn net.Conn) {
	//handle
	// fmt.Println("connect success", conn.RemoteAddr().String())

	//用户上线
	user := NewUser(conn)

	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	//广播消息

	s.broadCast(user, "已上线")
}

func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()
	fmt.Printf("server start: %s:%d\n", s.IP, s.Port)

	//启动message listen
	go s.listenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		//do handler
		go s.handler(conn)

	}

}
