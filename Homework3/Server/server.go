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
    streamMessage 					[]*gRPC.Message;
    // here you can implement other fields that you want
}

func (s *Server) Publish (ctx context.Context, message *gRPC.Message) (*gRPC.Message, error) {
    // some code here

	if 128 > utf8.RuneCountInString(message.Message){
		s.streamMessage = append(s.streamMessage, message);
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
