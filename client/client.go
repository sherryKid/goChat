package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"

	pb "../chat_msg"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Your nick name:")
	var name string
	fmt.Scanf("%s", &name)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("127.0.0.1:10000", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewChatClient(conn)
	stream, err := client.DoChat(context.Background())
	defer stream.CloseSend()
	if err != nil {
		log.Fatalf("%v.DoChat(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("[%s] %s", in.Name, in.Message)
		}
	}()
	stdin := bufio.NewReader(os.Stdin)
	for {
		msg, error := stdin.ReadString('\n')
		if error != nil {
			log.Fatalf("Failed to read line:%v", error)
		}
		note := &pb.ChatMsg{Name: name, Message: msg}
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}

}
