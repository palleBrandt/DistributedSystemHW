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
	msg := &gRPC.Message{Message: "Dette er en test"}
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

//Implementation of Brodcast, which takes a list of messages from the server
//And displays them in the terminal
func () Broadcast (msgStream gRPC.ChittyChat_BroadcastServer) error{ // skaal kaldes a man og åbne en stream der ikke lukkes før severen bliver slukket, denne stream skal client tappe ind på og lytte på om der kommer en message
	for {
			message, err := msgStream.Recv()
			//The io.EOF error means that the stream has reached the end of its list
			//And we can close the stream and break the loop
			if err == io.EOF {
				return stream.Close
				}
			//Error is returend
			if err != nil {
					return err;
				}
			//If no error occur, the message is printed to the terminal
			log.Println(message)
	}
	return nil
}


