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

	ConnectToServer()
	//Initialize the stream - setup streams for conecting with the server. This is to recieve broadcastet messages, and send messages.
	stream := SubscribeChittyChat();

	//Send an initial join message. Telling the server "hey this is my name, im the one on the other side of the stream"
	//this is so the server know who is joining / leaving. This is an internal event, which increments the lokal timestamp.
	t++;
	joinMessage := &gRPC.Message{
					AuthorName: userName,
					Text: "",
					LamportTimestamp: t}
		
	stream.Send(joinMessage)

	//Concurrently listen for incoming messages to log
	go Listen(stream);

	//Concurrently listen for typed messages in console that should be published to server.
	go Publish(stream);

	//This makes sure our client runs until we terminate it.
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
	//Loop listens for incoming broadcasted messages (from server)
	for{
		message, err := stream.Recv()
		if err != nil {
			if err != io.EOF{
				fmt.Println(err);
			}
			//Condition if it is a server messsage. This handles the formatting of how we look at join and leave messages from server.
		} else if message.AuthorName == "server"{
			//Formats increments the lokal timestamp and formats the message
			t = maxInt32(t, message.LamportTimestamp) + 1;
			fmt.Println(t);
			log.Println(message.Text);
		} else {
			//Formats increments the lokal timestamp and formats the message
			t = maxInt32(t, message.LamportTimestamp) + 1;
			fmt.Println(t);
			log.Println("Timestamp" , t , message.AuthorName , ": " , message.Text);
		}
	}
}

//published messages
func Publish (stream gRPC.ChittyChat_SubscribeClient){
	reader := bufio.NewReader(os.Stdin)

	//Listen for messages typed in console
	for{
		var inputText, _ = reader.ReadString('\n')
		inputText = strings.TrimRight(inputText,"\n")

		//Checks if message is less that 128 characters
		if 128 > utf8.RuneCountInString(inputText){
			t = t+1;
			publishMessage := &gRPC.Message{
				AuthorName: userName,
				Text: inputText,
				LamportTimestamp: t}
			
			//sends the message to the server.
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
	
	conn, err := grpc.Dial("192.168.43.179:5400", opts...)
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
	
	
	
	
	