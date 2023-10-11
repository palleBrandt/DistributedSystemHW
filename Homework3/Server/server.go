package main

import (

	// This has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	"context"
	"net"
	"unicode/utf8"
	"fmt"

	gRPC "github.com/palleBrandt/DistributedSystemHW/tree/main/Homework3/proto"
	"google.golang.org/grpc"
)

func main(){
	list, _ := net.Listen("tcp", "localhost:5400")
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)


	server := &Server{
		streamMessage: make ([]*gRPC.Message, 0,2)}
	
	gRPC.RegisterChittyChatServer(grpcServer, server)
	grpcServer.Serve(list)
}



type Server struct {
    // an interface that the server type needs to have
    gRPC.UnimplementedChittyChatServer
    savedMessages 					[]*gRPC.Message;
    // here you can implement other fields that you want
}

func (s *Server) Publish (ctx context.Context, message *gRPC.Message) (*gRPC.Message, error) { //ændre det til at den sender til en stream istedet for en liste
    // some code here

	if 128 > utf8.RuneCountInString(message.Message){
		s.savedMessages = append(s.savedMessages, message);
		succesMessage := &gRPC.Message{
    	Message: "fjing fjong ding dong"}
		fmt.Println(message.Message)
		return succesMessage, nil
	} else {
		errorMessage := &gRPC.Message{
    	Message: "Your shit is too long mf"}
		return errorMessage, nil
	}
}

func (s *Server) Broadcast (msgStream gRPC.ChittyChat_BroadcastServer) error{ // skaal kaldes a man og åbne en stream der ikke lukkes før severen bliver slukket, denne stream skal client tappe ind på og lytte på om der kommer en message
	for _, message := range s.savedMessages {
			if err := stream.Send(message); err != nil {
				return err
			}
	}
	return nil
}