package main

import (
	"context"
	mproto "myService/api/newinfo"

	"google.golang.org/grpc"
	"fmt"
	"sync"
	"io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
	NewAddr string `yaml:"new-grpc-addr"`
}

var _config Config

func init(){
	configContent , err := ioutil.ReadFile("../../config/config.yaml")	
	if err != nil {
		fmt.Println(err)
	}
	yaml.Unmarshal(configContent,&_config)
}

var ClientPool = sync.Pool {
	New : func()interface{} {
		conn, _ := grpc.Dial(_config.NewAddr,grpc.WithInsecure())
		return conn
	},
}

func GetClient() *grpc.ClientConn {
	obj ,ok := ClientPool.Get().(*grpc.ClientConn)
	if ok {
		return obj
	}
	return nil
}


func GetNews(ids []int32) *mproto.NewArrayInfo {
	conn := GetClient()
	defer ClientPool.Put(conn)
	client := mproto.NewAskNewsClient(conn)
	r , err := client.AskNews(context.Background(),&mproto.NewIds{ News:ids})
	if err !=nil {
		fmt.Println(err)
		return nil
	}
	return r
}

func main() {
	// conn ,_ := grpc.Dial(":14010",grpc.WithInsecure())
	// defer conn.Close()
	// client := mproto.NewAskNewsClient(conn)
	// client.CreateNews(context.Background(),&mproto.NewInfo{
	// 	Title: "test-title",
	// 	Content: "test",
	// })
	ids := []int32{1,2,3,4,5}
	r := GetNews(ids)
	fmt.Println(r)
}