package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {

	//Channels (for forks??):
	var ch1 = make(chan struct{})
	var ch2 = make(chan struct{})
	var ch3 = make(chan struct{})
	var ch4 = make(chan struct{})
	var ch5 = make(chan struct{})
	var ch6 = make(chan struct{})
	var ch7 = make(chan struct{})
	var ch8 = make(chan struct{})
	var ch9 = make(chan struct{})
	var ch10 = make(chan struct{})

	// go test()

	go fork(ch1, ch2)
	go fork(ch3, ch4)
	go fork(ch5, ch6)
	go fork(ch7, ch8)
	go fork(ch9, ch10)

	go phil(1, ch10, ch1)
	go phil(2, ch2, ch3)
	go phil(3, ch4, ch5)
	go phil(4, ch6, ch7)
	go phil(5, ch8, ch9)

	for {

	}
}

func phil(no int, left chan struct{}, right chan struct{}) {
	name := "philosopher " + strconv.Itoa(no)

	for {
		fmt.Println(name + " is thinking")
		//The possible deadlock in this aproach is that all philosophers grabs the right hand side fork.
		//This is avoided by making each philosopher think- and eat for a random amount of time.
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // time spent thinking
		right <- struct{}{}                                           //request fork
		<-right                                                       //recieve approval on request
		left <- struct{}{}
		<-left
		fmt.Println(name + " is eating")
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // time spent eating
		right <- struct{}{}                                           //send release to fork
		left <- struct{}{}
	}

}

func fork(left chan struct{}, right chan struct{}) {

	for {
		select {
		case <-right: //recieve request
			right <- struct{}{} //send approval
			<-right             //read release

		case <-left:
			left <- struct{}{}
			<-left

		}
	}
}
