package main

import (

	// This has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	
	"net"
	
	"sync"

	gRPC "github.com/palleBrandt/DistributedSystemHW/tree/main/Homework3/proto"
	"google.golang.org/grpc"
)

func main(){
	list, _ := net.Listen("tcp", "localhost:5400")
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	server := &Server{
		savedMessages: make ([]*gRPC.Message, 0,200),
		clients: make ([]gRPC.ChittyChat_JoinServer, 0, 200)}
	
	gRPC.RegisterChittyChatServer(grpcServer, server)
	grpcServer.Serve(list)
}

type Server struct {
	//Maakes it possible to lock the server
	sync.Mutex
    // an interface that the server type needs to have
    gRPC.UnimplementedChittyChatServer
    savedMessages 					[]*gRPC.Message;

	// A list of all streams created between the clients and the server
	clients							[]gRPC.ChittyChat_JoinServer
}

// Lets all users know that a new user has join. Sends a stream to the newly
// join user, on which it now can recive new messages, and fills it with
// all registered messages.
func (s *Server) Join (client *gRPC.Client, stream gRPC.ChittyChat_JoinServer) error {

	s.Lock();
	s.clients = append(s.clients, stream);
	s.Unlock();


	joinMessage := &gRPC.Message{AuthorName: client.Name, Text: "Participant " + client.Name + " joined Chitty-Chat at:"};
	s.broadcast(joinMessage);

	//Fills the stream, which is returned to the client, with the list saved messages
	for _,message := range s.savedMessages {
		if err := stream.Send(message); err != nil {
			return  err;
		}
	}
	return nil
}

// In your Server type, add a new method for the Chat streaming RPC.
func (s *Server) Chat(stream gRPC.ChittyChat_ChatServer) error {

	s.Lock();
	s.clients = append(s.clients, stream);
	s.Unlock();


	joinMessage := &gRPC.Message{AuthorName: "server", Text: "Participant Markus, joined Chitty-Chat at:"};
	s.broadcast(joinMessage);

	//Fills the stream, which is returned to the client, with the list saved messages
	for _,message := range s.savedMessages {
		if err := stream.Send(message); err != nil {
			return  err;
		}
	}

    // Create a channel to handle incoming client messages.
    messageChan := make(chan *gRPC.Message)

    // Goroutine to handle incoming messages from the client.
    go func() {
        for {
            message, err := stream.Recv()
            if err != nil {
                close(messageChan)
                return
            }
			s.savedMessages = append(s.savedMessages, message);
            messageChan <- message
        }
    }()

    // Now, you can use the messageChan to receive messages from clients and broadcast them.
    for message := range messageChan {
        // Handle the received message and broadcast it to other clients.
        s.broadcast(message)
    }
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