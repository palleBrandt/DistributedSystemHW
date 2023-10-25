package main

import (

	// This has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	"context"
	"net"
	"unicode/utf8"
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
func (s *Server) Join (stream gRPC.ChittyChat_JoinServer) error {

	s.Lock();
	s.clients = append(s.clients, stream);
	s.Unlock();


	//Fills the stream, which is returned to the client, with the list saved messages
	for _,message := range s.savedMessages {
		if err := stream.Send(message); err != nil {
			return  err;
		}
	}

	for {
        message, err := stream.Recv() // Receive a chat message from the client
        if err != nil {
            s.Lock()
            for i, client := range s.clients {
                if client == stream {
                    s.clients = append(s.clients[:i], s.clients[i+1:]...) // Remove the disconnected client
                    break
                }
            }
            s.Unlock()
            return err
        }
        s.Lock()
        s.savedMessages = append(s.savedMessages, message) // Store the new message in the chat history
        s.Unlock()
        s.broadcast(message) // Broadcast the new message to all connected clients
    }
}


// Recives messages from Users and provoke the broadcast methode to send the message
// To all users
func (s *Server) Publish (ctx context.Context, message *gRPC.Message) (*gRPC.Message, error) { 

	if 128 > utf8.RuneCountInString(message.Text){
		s.savedMessages = append(s.savedMessages, message);

		succesMessage := &gRPC.Message{
			AuthorName: "Server",
			Text: "All good"}
		
		//Sends the message to all user, via the brodcast methode
		s.broadcast(message);

		return succesMessage, nil
	} else {
		errorMessage := &gRPC.Message{
			AuthorName: "Server",
    		Text: "Your shit is too long mf"}
		return errorMessage, nil
	}
}

// Sends the message to all streams in the Cliens list.
func (s *Server) broadcast (message *gRPC.Message) error{ // skaal kaldes a man og åbne en stream der ikke lukkes før severen bliver slukket, denne stream skal client tappe ind på og lytte på om der kommer en message
	for _, client := range s.clients {
			if err := client.Send(message); err != nil {
				return err
			}
	}
	return nil
}