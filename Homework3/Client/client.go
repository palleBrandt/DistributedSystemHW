package main

import (
	"context"
	"fmt"
	"log"

	gRPC "github.com/palleBrandt/DistributedSystemHW/tree/main/Homework3/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ServerConn *grpc.ClientConn
var server gRPC.ChittyChatClient


func main(){
	ConnectToServer()
	msg := &gRPC.Message{Text: "Dette er en test"}
	returnMessage, error := server.Publish(context.Background(), msg)
	fmt.Println(returnMessage.Message , error)
}

func ConnectToServer(){
	var opts []grpc.DialOption
	opts = append(
    opts, grpc.WithBlock(), 
    grpc.WithTransportCredentials(insecure.NewCredentials()),	
	)

	conn, err := grpc.Dial(":5400", opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}

	server = gRPC.NewChittyChatClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())
}

func conReady(s gRPC.ChittyChatClient) bool {
	return ServerConn.GetState().String() == "READY"
}



