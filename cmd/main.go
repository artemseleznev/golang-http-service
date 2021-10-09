package main

import (
	"context"
	"fmt"
	"golang.org/x/net/netutil"
	"http-service-example/env"
	"http-service-example/server"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	env, err := env.InitEnv()
	if err != nil {
		log.Fatalf("could not initialize env vars: %v", err)
	}

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	s := server.New(env.ReqTimeout)
	http.HandleFunc("/contents", s.Contents)
	srv := &http.Server{}

	go func() {
		if err := httpServe(srv, env); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	osCall := <-done
	log.Printf("system call: %+v", osCall)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		log.Print("http server stopped")
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}

func httpServe(srv *http.Server, env *env.Env) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", env.HttpPort))
	if err != nil {
		return err
	}
	listener = netutil.LimitListener(listener, env.MaxConns)

	log.Printf("starting http server on %s...", env.HttpPort)
	return srv.Serve(listener)
}
