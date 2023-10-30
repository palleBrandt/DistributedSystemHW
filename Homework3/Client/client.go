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

	gRPC "github.com/palleBrandt/DistributedSystemHW/tree/main/Homework3/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ServerConn *grpc.ClientConn
var server gRPC.ChittyChatClient
var userName string
var t int32

func main(){
	// set lokal logical timestamp 
	t = 0;

	//Extracting the username
	user, err := user.Current()
	if err != nil{
		fmt.Println(err.Error())
	} else {
		userName = user.Username
	}

	//Setup log
	file, _:= os.Create("ChittyChatLog.txt")
	log.SetOutput(file);

	// ClientUser := &gRPC.Client{Name: userName}
	ConnectToServer()
	//Initialize the stream
	stream := SubscribeChittyChat();
	joinMessage := &gRPC.Message{
					AuthorName: userName,
					Text: ""}
			
				stream.Send(joinMessage)

	go Listen(stream);

	go Publish(stream);

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
	log.Println("You are now connected to Chitty-Chat, type away as hard as you are :)")
	return stream;
}

//Logs recived messages
func Listen (stream gRPC.ChittyChat_SubscribeClient){
	for{
		message, err := stream.Recv()
		if err != nil {
			if err != io.EOF{
				fmt.Println(err);
			}
			//Condition if it is a server messsage
		} else if message.AuthorName == "server"{
			t = maxInt32(t, message.LamportTimestamp) + 1;
			log.Println(message.Text);
			//Formats increments the lokal timestamp and formats the message
		} else {
			t = maxInt32(t, message.LamportTimestamp) + 1;
			log.Println("Timestamp" , t , message.AuthorName , ": " , message.Text);
		}
	}
}

func Publish (stream gRPC.ChittyChat_SubscribeClient){
	reader := bufio.NewReader(os.Stdin)

	for{
		var inputText, _ = reader.ReadString('\n')
		inputText = strings.TrimRight(inputText,"\n")

		if 128 > utf8.RuneCountInString(inputText){
			t = t+1;
			publishMessage := &gRPC.Message{
				AuthorName: userName,
				Text: inputText,
				LamportTimestamp: t}
		
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
	
	conn, err := grpc.Dial("localhost:5400", opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}
	
	server = gRPC.NewChittyChatClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())
}

func maxInt32(a, b int32) int32 {
		if a > b {
			return a
		}
		return b
	}

func conReady(s gRPC.ChittyChatClient) bool {
	return ServerConn.GetState().String() == "READY"
}
	
	
	
	
	