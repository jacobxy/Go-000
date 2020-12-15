package main

import (
	"context"
	"log"
	mproto "myService/api/newinfo"
	"myService/models"
	"net"

	"google.golang.org/grpc"
	"io/ioutil"
    "gopkg.in/yaml.v2"
	"fmt"
)

// var port = "14010"

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

func main() {
	server := grpc.NewServer()
	mproto.RegisterAskNewsServer(server, &NewService{})
	lis, _ := net.Listen("tcp", _config.NewAddr)
	server.Serve(lis)
}

type NewService struct{}

func (this *NewService) AskOneNew(ctx context.Context, r *mproto.NewId) (*mproto.NewInfo, error) {
	res := models.ReadNew(int(r.Newid))
	return &res.NewInfo, nil
}

func (this *NewService) AskNews(ctx context.Context, r *mproto.NewIds) (*mproto.NewArrayInfo, error) {
	rr := models.ReadNews(r.News)
	if len(rr) == 0{
	 	return &mproto.NewArrayInfo{}, nil
	}
	res := &mproto.NewArrayInfo{
		NewsInfo:make([]*mproto.NewInfo,0,100),
	}
	for _,v := range rr {
		if v == nil {
			continue
		}	
		res.NewsInfo = append(res.NewsInfo, &v.NewInfo)
	}
	return res, nil
}

func (this *NewService) AskNewBeginAndEnd(ctx context.Context, r *mproto.NewBeginAndEnd) (*mproto.NewArrayInfo, error) {
	return &mproto.NewArrayInfo{}, nil
}

func (this *NewService) CreateNews(ctx context.Context, r *mproto.NewInfo) (*mproto.NewInfo, error) {
	log.Println("CreateNews")
	if r == nil {
		return nil, nil
	}
	models.InsertNews(&models.New{NewInfo: *r})
	return &mproto.NewInfo{}, nil
}
