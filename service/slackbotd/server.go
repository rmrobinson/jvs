package slackbotd

import (
	"fmt"
	"time"

	"faltung.ca/jvs/lib/proto-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	messageClientConn *grpc.ClientConn
	messageClient     proto.MessageManagerClient

	slackClient slackClient
}

func (s *Server) processMessageClientMessages() {
	// TODO: broadcast updates to Slack
}

func (s *Server) processSlackClientMessages(m message) string {
	if s.messageClientConn == nil {
		return "Not connected, ignoring"
	}

	req := proto.SendMessageRequest{
		Message: &proto.Message{
			Id:       m.event.EventTimestamp,
			User:     m.User,
			Contents: m.Contents,
		},
	}
	resp, err := s.messageClient.SendMessage(context.Background(), &req)

	if err != nil {
		return "Unable to send message: " + err.Error()
	}

	return resp.Message.Contents
}

func (s *Server) Run(slackToken string, commanderAddr string) {
	var err error
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	s.messageClientConn, err = grpc.Dial(commanderAddr, opts...)

	if err != nil {
		fmt.Printf("Unable to connect: %s\n", err.Error())
		return
	}

	s.messageClient = proto.NewMessageManagerClient(s.messageClientConn)

	messages := make(chan message)

	go s.slackClient.run(slackToken, "home", messages)

	// TODO: run goroutine monitoring messageClient

	for {
		m := <-messages

		m.Contents = s.processSlackClientMessages(m)
		m.Timestamp = time.Now().Format(time.RFC3339)

		s.slackClient.sendMessageReply(m)
	}
}
