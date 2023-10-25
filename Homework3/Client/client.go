package main

import (
"context"
"fmt"
"log"
"os/user"
"io"

gRPC "github.com/palleBrandt/DistributedSystemHW/tree/main/Homework3/proto"

"google.golang.org/grpc"
"google.golang.org/grpc/credentials/insecure"
)

var ServerConn *grpc.ClientConn
var server gRPC.ChittyChatClient
var userName string

func main(){

user, err := user.Current()
if err != nil{
	fmt.Println(err.Error())
} else {
	userName = user.Username
}

ConnectToServer()

//Initialize the stream
stream := JoinChittyChat();

go Listen(stream);

for{
	var inputText string
	fmt.Scanln(&inputText)

	publishMessage := &gRPC.Message{
			AuthorName: userName,
			Text: inputText}
	
	returnMessage, err := server.Publish(
		context.Background(),
		publishMessage,
	)
		if err != nil{
			fmt.Println(err)
			} else {
				fmt.Println(returnMessage)
			}
	}
}
	
	// Calls the Join method to join the server and return the stream
	func JoinChittyChat() gRPC.ChittyChat_JoinClient{
		stream, err := server.Join(
			context.Background(),
			&gRPC.Client{
				Name: userName,
			},
		)
		if err != nil {
			fmt.Println(err);
		}
		fmt.Println("You are now connected to Chitty-Chat, type away as hard as you are :)")
		return stream;
	}
	
	
	func Listen (stream gRPC.ChittyChat_JoinClient){
		for{
			message, err := stream.Recv()
			if err != nil && err != io.EOF {
				fmt.Println(err);
			} else {
				fmt.Println(message.AuthorName , ": " , message.Text);
			}
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
	
	
	
	