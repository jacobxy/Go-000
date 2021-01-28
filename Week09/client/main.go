package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	Conn net.Conn
}

func (c *Client) Rec() {
	local := c.Conn.LocalAddr().String()
	for {
		buf := [1024]byte{}
		n, err := c.Conn.Read(buf[:])
		if err != nil {
			log.Println("client Read Error", err)
			break
		}
		str := string(buf[:n])
		fmt.Println(local, "RecMessage ", str)
	}
}

func (c *Client) SendMessage() {
	connect := fmt.Sprintf(" Hello Say At %d", time.Now().Unix())
	_, err := c.Conn.Write([]byte(connect))
	if err != nil {
		log.Println("send message error ", err)
		return
	}
	local := c.Conn.LocalAddr().String()
	fmt.Println(local, "SendMessage ", connect)
}

func NewClient() *Client {
	conn, err := net.Dial("tcp", "localhost:10001")
	if err != nil {
		log.Println("connect error")
		return nil
	}
	c := &Client{
		Conn: conn,
	}
	go c.Rec()
	return c
}

func main() {
	mc := [5]*Client{}
	for k := range mc {
		mc[k] = NewClient()
	}
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		for _, v := range mc {
			if v.Conn == nil {
				continue
			}
			v.SendMessage()
		}
	}
	time.Sleep(5 * time.Second)
}
