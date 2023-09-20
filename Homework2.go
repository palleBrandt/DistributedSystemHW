package main

import (
	"fmt"
	"math/rand"
	"time"
)

//Convention: sequnceX, sequnceY, data

func main() {
	var connection bool = false
	var IPv4 = make(chan []int, 100)

	for connection == false {
		go Client(IPv4)
		go Server(IPv4)
	}
}

func Client(IPv4 chan []int) {
	sequenceX := rand.Intn(3000000000)
	shake1 := []int{sequenceX, 0}
	IPv4 <- shake1
	var shake2 []int
	for {
		select {
		case shake2 = <-IPv4:
			if shake2[0] == sequenceX+1 {
				sequenceX := sequenceX + 1
				ackY := shake2[1] + 1 //creates val for sequenceY on client-side
				shake3 := []int{sequenceX, ackY}
				IPv4 <- shake3
				fmt.Println("Client ready")
			} else {
				fmt.Println("Page not found: error 404", shake2)
				return // ack denied
			}
		case <-time.After(2 * time.Second):
			return
		}
	}
}

func Server(IPv4 chan []int) {
	sequenceY := rand.Intn(3000000000)
	shake1 := <-IPv4
	ackX := shake1[0] + 1
	shake2 := []int{ackX, sequenceY}
	IPv4 <- shake2
	var shake3 []int
	for {
		select {
		case shake3 = <-IPv4:
			if shake3[1] == sequenceY+1 {
				sequenceY = sequenceY + 1
				fmt.Println("Connection established!")
			} else {
				fmt.Println("Page not found: error 404", shake3)
				return // ack denied
			}
		case <-time.After(2 * time.Second):
			return
		}
	}
}
