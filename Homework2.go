package main
import "fmt"
import "math/rand"
import "time"

//Convention: sequnceX, sequnceY, data

func main() {
	var IPv4 = make(chan []int);

	go Client(IPv4);
	go Server(IPv4);
	for {}
}

func Client (IPv4 chan []int) {
	sequenceX := rand.Intn(3000000000);
	shake1 := []int{sequenceX, 0}
	loop:
	for{
		select{
		case IPv4 <- shake1:
			fmt.Println("sendt")
			break loop
		case <- time.After(2 * time.Second):
			fmt.Println("client not response")
			break loop
		}
	}
	shake2 := <- IPv4
	if shake2[0] == sequenceX + 1  {
		sequenceX := sequenceX + 1
		ackY := shake2[1] + 1 //creates val for sequenceY on client-side
		shake3 := []int{sequenceX, ackY}
		IPv4 <- shake3;
		fmt.Println("Client ready")
	} else {fmt.Println("Page not found: error 404", shake2)}
}

func Server (IPv4 chan []int ) {
	sequenceY := rand.Intn(3000000000);
	shake1 := <- IPv4
	ackX := shake1[0] + 1;
	shake2 := []int{ackX, sequenceY};
	loop:
	for{
		select{
		case IPv4 <- shake2:
			fmt.Println("sendt2")
			break loop
		case <- time.After(2 * time.Second):
			fmt.Println("client not response")
			break loop
		}
	}
	shake3 := <- IPv4
	if shake3[1] == sequenceY + 1 { //Connection established
		ackX = ackX + 1
		fmt.Println("Connection established succesfully")
		// request1 := <- IPv4
	} else {fmt.Println("Connection access denied: error 406")}

	

	// Rækkefølge, ikke være dobbelt, tjekke alle pakker er der
}

