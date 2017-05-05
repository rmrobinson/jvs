package main

import (
	"google.golang.org/grpc"
	"faltung.ca/jvs/lib/proto-go"
	"flag"
	"golang.org/x/net/context"
	"fmt"
)

var (
	addr = flag.String("server", "127.0.0.1:10004", "The server IP:port to connect to")
)
func main() {
	flag.Parse()

	var err error
	var conn *grpc.ClientConn
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err = grpc.Dial(*addr, opts...)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	msg := proto.NewMessageManagerClient(conn)

	req := proto.SendMessageRequest{
		Message: &proto.Message{
			Id: "1",
			User: "Robert",
			Contents: "device get light/1",
		},
	}
	resp, err := msg.SendMessage(context.Background(), &req)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Response: %s\n", resp.Message.Contents)

}
