package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HandleSingal() context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGKILL, os.Interrupt)
	go func() {
		select {
		case <-signalChannel:
			cancel()
		}
	}()
	return ctx
}

type Message struct {
	Con net.Conn
	Msg string
}

type Server struct {
	Listen  net.Listener
	MsgChan chan Message
}

func (s *Server) HandleMessage(ctx context.Context) {
	for {
		select {
		case v := <-s.MsgChan:
			answer := fmt.Sprintf("I Get the Message %s already:", v.Msg)
			fmt.Fprintln(v.Con, answer)
		case <-ctx.Done():
			log.Println("Server HandleMessage Exit")
			return
		}
	}
}

func (s *Server) StartServer() {
	var err error
	s.Listen, err = net.Listen("tcp", ":10001")
	if err != nil {
		log.Fatalln(err)
	}
	ctx := HandleSingal()
	go s.HandleMessage(ctx)
	go s.RecConn(ctx)
	select {
	case <-ctx.Done():
		s.Listen.Close()
		log.Println("Server Exit")
	}
}
func (s *Server) RecConn(ctx context.Context) {
	for {
		conn, err := s.Listen.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go s.HandleConn(ctx, conn)
	}
}

func (s *Server) HandleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Println("Read Error", err)
			break
		}
		str := string(buf[:n])
		select {
		case s.MsgChan <- Message{Con: conn, Msg: str}:
		case <-ctx.Done():
			log.Println("Server HandleConn Exit")
			return
		}
	}
}

func main() {
	s := Server{
		MsgChan: make(chan Message, 10),
	}
	s.StartServer()
	time.Sleep(5 * time.Second)
}
