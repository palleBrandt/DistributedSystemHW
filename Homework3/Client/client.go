package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gRPC "github.com/palleBrandt/DistributedSystemHW/tree/main/Homework3/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ServerConn *grpc.ClientConn
var server gRPC.ChittyChatClient

func main(){

	ConnectToServer()

	//Initialize the stream
	stream := JoinChittyChat();

	go Listen(stream);
	
	scanner := NewScanner(os.Stdin)

	for{
		inputText := scanner.Text()
		returnMessage, err := server.Publish(
			context.Background(),
			AuthorName: os.Hostname(),
			Text: inputText
		)
		if err != nil{
			fmt.Println(err)
		} else {
			fmt.Println(returnMessage)
		}
	}

}

// Calls the Join method to join the server and return the stream
func JoinChittyChat() gRPC.ChittyChat_JoinServer{
	stream, error := server.Join(
		context.Background(),
		&gRPC.Client{
			Name: os.Hostname()
		}
	)
	if err != nil {
		fmt.Println(err);
	}
	fmt.Println("You are now connected to Chitty-Chat, type away as hard as you are :)")
	return stream;
}


func Listen (stream gRPC.ChittyChat_JoinServer){
	for{
		message, err := stream.Recv()
		if err != nil {
			fmt.Println(err);
		}
		fmt.Println(message.AuthorName + ": " + message.Text);
	}
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



