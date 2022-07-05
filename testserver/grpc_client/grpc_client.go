package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/testserver/real_grpc_server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	timestampFormat = time.UnixDate
)

func unaryCallWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("--- unary ---\n")

	// Create metadata and context.
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("msg", "你好")

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("failed to call UnaryEcho: %v", err)
	}
	fmt.Printf("response:%v\n", r.Message)
}

func serverStreamingWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("--- server streaming ---\n")

	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.ServerStreamingEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("failed to call ServerStreamingEcho: %v", err)
	}

	// Read all the responses.
	var rpcStatus error
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		logx.Debug(" - %s", r.Message)
	}
	if rpcStatus != io.EOF {
		log.Fatalf("failed to finish server streaming: %v", rpcStatus)
	}
}

func clientStreamWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("--- client streaming ---\n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("failed to call ClientStreamingEcho: %v\n", err)
	}

	// Send all requests to the server.
	if err := stream.Send(&pb.EchoRequest{Message: message}); err != nil {
		log.Fatalf("failed to send streaming: %v\n", err)
	}

	// Read the response.
	r, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to CloseAndRecv: %v\n", err)
	}
	fmt.Printf("response:%v\n", r.Message)
}

func bidirectionalWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("--- bidirectional ---\n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("failed to call BidirectionalStreamingEcho: %v\n", err)
	}

	go func() {
		// Send all requests to the server.
		if err := stream.Send(&pb.EchoRequest{Message: message}); err != nil {
			log.Fatalf("failed to send streaming: %v\n", err)
		}
		stream.CloseSend()
	}()

	// Read all the responses.
	var rpcStatus error
	fmt.Printf("response:\n")
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf(" - %s\n", r.Message)
	}
	if rpcStatus != io.EOF {
		log.Fatalf("failed to finish server streaming: %v", rpcStatus)
	}
}

const message = "hello world你好"

func main() {

	var port int
	if len(os.Args) == 1 {
		port = 2003
	} else {
		port, _ = strconv.Atoi(os.Args[1])
	}

	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	//一元方法
	unaryCallWithMetadata(c, message)
	time.Sleep(400 * time.Millisecond)

	//服务端流式
	serverStreamingWithMetadata(c, message)
	time.Sleep(1 * time.Second)

	//客户端流式
	clientStreamWithMetadata(c, message)
	time.Sleep(1 * time.Second)

	//双向流式
	bidirectionalWithMetadata(c, message)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)
	<-interrupt
}
