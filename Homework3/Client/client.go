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

 client := gRPC.NewChittyChatClient(conn)

    stream, err := client.Chat(context.Background())
    if err != nil {
        log.Fatalf("Error creating stream: %v", err)
    }
    defer stream.Close()

    // Start a goroutine to send messages to the server through the stream
    go func() {
        for {
            var inputText string
            fmt.Scanln(&inputText)

            message := &gRPC.Message{
                AuthorName: "YourName", // Your name or username
                Text:       inputText,
            }

            if err := stream.Send(message); err != nil {
                log.Fatalf("Error sending message: %v", err)
            }
        }
    }()

    // Start a goroutine to receive messages from the server through the stream
    go func() {
        for {
            message, err := stream.Recv()
            if err != nil {
                if err == io.EOF {
                    // The stream has been closed by the server.
                    break
                }
                log.Fatalf("Error receiving message: %v", err)
            }

            // Process and display the received message
            fmt.Printf("%s: %s\n", message.AuthorName, message.Text)
        }
    }()

    // Keep the main goroutine running
    select {}
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
			if err != nil {
				if err != io.EOF{
					fmt.Println(err);
				}
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
	
	
	
	