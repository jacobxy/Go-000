package main

import (
	"context"

	// "errors"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	SIGUSR1 = syscall.Signal(0xa)
	SIGUSR2 = syscall.Signal(0xc)
)

type MHandle struct {
}

func (m *MHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("helloword"))
}

func GetServer() *MServer {
	s := &MServer{}
	s.Reset()
	return s
}

func (m *MServer) Reset() {
	m.Server = &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &MHandle{},
	}

	m.Server.RegisterOnShutdown(func() {
		log.Println("http server shut down ")
	})
}

type MServer struct {
	*http.Server
	isRunning bool
}

func (m *MServer) ListenAndServe() error {
	if m.Server == nil {
		m.Reset()
	}
	if m.isRunning {
		return nil
	}
	m.isRunning = true
	log.Printf("[HTTP] 监听端口: %s\n", m.Server.Addr)
	return m.Server.ListenAndServe()
}

func (m *MServer) Shutdown(ctx context.Context) error {
	if err := m.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shut down error")
	}
	// m.Reset()
	m.Server = nil
	m.isRunning = false
	return nil
}

func main() {

	// ctx0, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	ctx0, cancel := context.WithCancel(context.Background())
	// ctx0 := context.Background()
	group, ctx := errgroup.WithContext(ctx0)

	s := GetServer()

	c := make(chan os.Signal)
	signal.Notify(c, SIGUSR1, SIGUSR2, os.Interrupt)

	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), "Signal Exit")
			case v := <-c:
				{
					log.Println(v)
					switch v {
					case os.Interrupt:
						cancel()
					case os.Kill:
						return errors.New("os kill")
					case SIGUSR1:
						group.Go(func() error {
							err := s.ListenAndServe()
							log.Println(err)
							return nil
						})
					case SIGUSR2:
						if err := s.Shutdown(ctx); err != nil {
							return s.Close()
							// return nil
						}
					}
				}
			}
		}
	})

	group.Go(func() error {
		//启动信号
		c <- SIGUSR1
		log.Println(GetHttpResponse("http://127.0.0.1:8080"))
		time.Sleep(2 * time.Second)
		//关闭信号
		c <- SIGUSR2
		time.Sleep(2 * time.Second)
		log.Println(GetHttpResponse("http://127.0.0.1:8080"))
		//启动信号
		c <- SIGUSR1
		time.Sleep(2 * time.Second)
		log.Println(GetHttpResponse("http://127.0.0.1:8080"))
		//关闭信号
		c <- SIGUSR2
		//退出信号
		c <- os.Interrupt
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Println("wait error ", err)
	}

	log.Println("process end")
}

func GetHttpResponse(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "http get error")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "http read error")
	}
	return string(b), nil
}
