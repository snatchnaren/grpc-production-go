package grpc_client

import (
	"context"
	"github.com/apssouza22/grpc-server-go/util"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func TimeoutLogExample() {
	clientBuilder := GrpcClientBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithContext(context.Background())
	clientBuilder.WithStreamInterceptors(util.GetDefaultStreamClientInterceptors())
	clientBuilder.WithUnaryInterceptors(util.GetDefaultUnaryClientInterceptors())
	cc, err := clientBuilder.GetConn("localhost", "50051")

	defer cc.Close()
	ctx := context.Background()
	md := metadata.Pairs("user", "user", "pass", "123")
	ctx = metadata.NewOutgoingContext(ctx, md)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	timeout := time.Minute * 1
	ctx, cancel := context.WithTimeout(ctx, timeout)
	client := helloworld.NewGreeterClient(cc)
	request := &helloworld.HelloRequest{
		Name: "mike",
	}
	healthClient := grpc_health_v1.NewHealthClient(cc)
	response, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%v", response)
	helloReply, err := client.SayHello(ctx, request)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%v", helloReply)

	defer cancel()
}
