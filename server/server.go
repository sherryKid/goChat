package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "../chat_msg"
)

// protoc -I chat_msg/ chat_msg/chat_msg.proto --go_out=plugins=grpc:chat_msg

type chatServer struct {
}

var streamMap map[string]pb.Chat_DoChatServer
var mutex *sync.Mutex

func (s *chatServer) DoChat(stream pb.Chat_DoChatServer) error { //chatServer may struct in this file

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		name := in.Name
		msg := in.Message
		if streamMap[name] == nil {
			mutex.Lock()
			if streamMap[name] == nil {
				streamMap[name] = stream
			}
			mutex.Unlock()
		}
		fmt.Println("msg:[", name, "]", msg)
		for k, v := range streamMap {
			if k != name {
				v.Send(in)
			}
		}
	}
}

func newServer() *chatServer {
	s := new(chatServer)
	return s
}
func main() {
	mutex = new(sync.Mutex)
	streamMap = make(map[string]pb.Chat_DoChatServer)
	lis, err := net.Listen("tcp", "localhost:10000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServer(grpcServer, newServer())
	fmt.Println("start server...")
	grpcServer.Serve(lis)
}
