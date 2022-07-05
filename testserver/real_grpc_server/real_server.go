package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/pkg/etcd"
	"github.com/josexy/gw/testserver/real_grpc_server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RealServer struct {
	Addr       string
	ServiceReg *etcd.ServiceRegister
	pb.UnimplementedEchoServer
}

func (server *RealServer) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	logx.Error("--- ServerStreamingEcho ---")
	logx.Debug("request received: %v", in)
	// Read requests and send responses.
	logx.Debug("echo message %v", in.Message)
	err := stream.Send(&pb.EchoResponse{Message: in.Message})
	if err != nil {
		return err
	}
	return nil
}

func (server *RealServer) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	logx.Error("--- ClientStreamingEcho ---")
	// Read requests and send responses.
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			logx.Debug("echo last received message")
			return stream.SendAndClose(&pb.EchoResponse{Message: message})
		}
		message = in.Message
		logx.Debug("request received: %v, building echo", in)
		if err != nil {
			return err
		}
	}
}

func (server *RealServer) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	logx.Error("--- BidirectionalStreamingEcho ---")
	// Read requests and send responses.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		logx.Debug("request received %v, sending echo", in)
		if err := stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func (server *RealServer) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	logx.Error("--- UnaryEcho ---")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logx.Debug("miss metadata from context")
	}
	logx.Warn("metadata %v", md)
	logx.Debug("request received: %v, sending echo", in)
	return &pb.EchoResponse{Message: in.Message}, nil
}

func (server *RealServer) Run() {
	logx.Debug("start server at: %v", server.Addr)

	lis, err := net.Listen("tcp", server.Addr)
	if err != nil {
		logx.Fatal("net.Listen err: %v", err)
		return
	}

	// 开启一个goroutine监听服务
	go func() {
		s := grpc.NewServer(grpc.ChainStreamInterceptor())
		pb.RegisterEchoServer(s, server)
		if err = s.Serve(lis); err != nil {
			logx.Fatal("%v", err)
			return
		}
	}()

	// 服务注册
	server.registerService(server.Addr)
}

func (server *RealServer) unregisterService() {
	logx.Debug("revoke service and close server: %s", server.Addr)
	if server.ServiceReg != nil {
		_ = server.ServiceReg.Close()
	}
}

func (server *RealServer) registerService(addr string) {

	endpoints := []string{"127.0.0.1:2379"}
	prefix := constants.EtcdPrefix

	go func() {
		service, err := etcd.NewServiceRegister(endpoints, 5)
		if err != nil {
			logx.Error("connect server register: %v", err)
			return
		}
		server.ServiceReg = service

		key := prefix + addr
		logx.Debug("register service to etcd: %s", key)
		_ = service.RegisterService(key, addr)
	}()
}

func main() {

	var port1, port2 int
	var err error

	if len(os.Args) == 1 {
		port1 = 2003
		port2 = 2004
	} else {
		port1, err = strconv.Atoi(os.Args[1])
		if err != nil {
			logx.Fatal("%v", err)
			return
		}
		port2, err = strconv.Atoi(os.Args[2])
		if err != nil {
			logx.Fatal("%v", err)
			return
		}
	}

	rs1 := &RealServer{Addr: fmt.Sprintf("127.0.0.1:%d", port1)}
	rs2 := &RealServer{Addr: fmt.Sprintf("127.0.0.1:%d", port2)}

	rs1.Run()
	rs2.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)
	<-interrupt

	rs1.unregisterService()
	rs2.unregisterService()
}
