package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	gw "grpc-helloworld-gateway/helloworld"
	"os"
	"os/signal"
	"syscall"
)

type Greeter struct {

}

func (g *Greeter) SayHello(ctx context.Context,reqest *gw.HelloRequest) (*gw.HelloReply,error){
	return &gw.HelloReply{
		Message:              "我是测试数据",
	},nil
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		lis, err := net.Listen("tcp", "localhost:8181")
		if err != nil {
			log.Fatalf("Failed start TCP Server on %s,  %v", "localhost:8181")
		}
		grpcServer := grpc.NewServer()
		gw.RegisterGreeterServer(grpcServer,new(Greeter))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start GRPC Server on %s : %v", "localhost:8181", err)
		}
	}()
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux :=  runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
		}
		//这里特别注意，注册端口一定要与grpc 服务端口一致 net.Listen("tcp", "localhost:8181")
		err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, "localhost:8181", opts)
		if err!=nil {
			panic(err)
		}
		//cors := cors.AllowAll()
		//serverMux := http.NewServeMux()
		//serverMux.Handle("/", cors.Handler(grpcMux))
		 http.ListenAndServe("localhost:8182", mux)
	}()
	<-stop
}
