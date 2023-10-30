package main

import (

	// This has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.

	"net"
	"sync"
	"fmt"
	"strconv"

	gRPC "github.com/palleBrandt/DistributedSystemHW/tree/main/Homework3/proto"
	"google.golang.org/grpc"
)

func main(){
	list, _ := net.Listen("tcp", "10.26.26.4:5400")
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)


	server := &Server{
		t: 0,
		clients: make ([]gRPC.ChittyChat_SubscribeServer, 0, 200)}
		
	
	gRPC.RegisterChittyChatServer(grpcServer, server)
	grpcServer.Serve(list)
}



type Server struct {
	//Maakes it possible to lock the server
	sync.Mutex
    // an interface that the server type needs to have
    gRPC.UnimplementedChittyChatServer

	// A list of all streams created between the clients and the server
	clients	[]gRPC.ChittyChat_SubscribeServer
	t int32;
}

// Lets all users know that a new user has join. Sends a stream to the newly
// join user, on which it now can recive new messages, and fills it with
// all registered messages.
func (s *Server) Subscribe (stream gRPC.ChittyChat_SubscribeServer) error {

	s.Lock();
	s.clients = append(s.clients, stream);
	s.Unlock();

	 clientMessage, err := stream.Recv() // Receive a chat message from the client
        if err != nil {
            fmt.Println(err);
        }
	s.Join(clientMessage);

	for {
        message, err := stream.Recv() // Receive a chat message from the client
        if err != nil {
            s.Lock()
            for i, client := range s.clients {
                if client == stream {
                    s.clients = append(s.clients[:i], s.clients[i+1:]...) // Remove the disconnected client
					s.Leave(clientMessage)
                    break
                }
            }
            s.Unlock()
            return err
        }
        s.Lock()
		//Increments timestamp for recieving a message
		s.t = maxInt32(s.t, message.LamportTimestamp) + 1;
        s.Unlock()
		//Increments timestamp for sending a message
		s.t ++;
        s.broadcast(message) // Broadcast the new message to all connected clients
    }
}

// Sends the message to all streams in the Cliens list.
func (s *Server) Join (message *gRPC.Message) error{
	//Increments timestamp for when a client joins the server
	s.t ++;
	joinMessage := &gRPC.Message{AuthorName: "server", Text: "Participant " + message.AuthorName + " joined Chitty-Chat at Lamport time: " + strconv.FormatInt(int64(s.t), 10)};
	s.broadcast(joinMessage);
	return nil
}

// Sends the message to all streams in the Cliens list.
func (s *Server) Leave (message *gRPC.Message) error{ 
	//Increments timestamp for when a client leaves the server
	s.t ++;
	leaveMessage := &gRPC.Message{AuthorName: "server", Text: "Participant " + message.AuthorName + " left Chitty-Chat at Lamport time: " + strconv.FormatInt(int64(s.t), 10)};
	s.broadcast(leaveMessage);
	return nil
}

// Sends the message to all streams in the Cliens list.
func (s *Server) broadcast (message *gRPC.Message) error{
	for _, client := range s.clients {
			if err := client.Send(message); err != nil {
				return err
			}
	}
	return nil
}

func maxInt32(a, b int32) int32 {
		if a > b {
			return a
		}
		return b
	}