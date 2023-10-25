package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
	"unicode/utf8"
	"os/signal"
    "syscall"

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

	ClientUser := &gRPC.Client{Name: userName}

	ConnectToServer()
	server.Join(context.Background(),ClientUser)

	//Initialize the stream
	stream := SubscribeChittyChat();

	go Listen(stream);

	go Publish(stream);

	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        server.Leave(context.Background(),ClientUser)
        os.Exit(1)
    }()

	select {}

}
	
	// Calls the Join method to join the server and return the stream
	func SubscribeChittyChat() gRPC.ChittyChat_SubscribeClient{
		stream, err := server.Subscribe(
			context.Background(),
		)
		if err != nil {
			fmt.Println(err);
		}
		fmt.Println("You are now connected to Chitty-Chat, type away as hard as you are :)")
		return stream;
	}
	
	
	func Listen (stream gRPC.ChittyChat_SubscribeClient){
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

	func Publish (stream gRPC.ChittyChat_SubscribeClient){
		reader := bufio.NewReader(os.Stdin)

		for{
			var inputText, _ = reader.ReadString('\n')
			inputText = strings.TrimRight(inputText,"\n")

			if 128 > utf8.RuneCountInString(inputText){
				publishMessage := &gRPC.Message{
					AuthorName: userName,
					Text: inputText}
			
				stream.Send(publishMessage)
			} else {
				fmt.Println("!Maximum 128 characters allowed!")
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
	
	
	
	