package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	gw "grpc-helloworld-gateway/helloworld"
)

func main() {
	conn, err := grpc.Dial("localhost:8181", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := gw.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(),&gw.HelloRequest{
		Name:                 "111",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("-->>", resp.Message)
}
