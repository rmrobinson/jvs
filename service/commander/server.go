package commander

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"faltung.ca/jvs/lib/proto-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server

	rootCommand RootCommand
}

func NewServer() *Server {
	s := &Server{}

	return s
}

func (s *Server) SetRootCommand(cmd RootCommand) {
	s.rootCommand = cmd
}

func (s *Server) Run(port int) error {
	// Run the gRPC listener.
	// This is assumed to never return.
	s.grpcServer = grpc.NewServer()

	proto.RegisterMessageManagerServer(s.grpcServer, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer lis.Close()

	if err != nil {
		return err
	}

	s.grpcServer.Serve(lis)

	return nil
}

func (s *Server) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	resp := &proto.SendMessageResponse{}
	resp.Message = &proto.Message{}

	contents := strings.Split(req.Message.Contents, " ")

	result, err := s.rootCommand.Execute(ctx, contents)

	if err != nil {
		return nil, err
	}

	resp.Message.Contents = result
	return resp, nil
}

func (s *Server) WatchMessages(req *proto.WatchMessagesRequest, stream proto.MessageManager_WatchMessagesServer) error {
	return errors.New("Not Implemented")
}
