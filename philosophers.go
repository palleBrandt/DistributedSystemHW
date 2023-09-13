package main
import "fmt"

func main() {
	//Channels (for forks??):
	var ch1 = make(chan bool);
	var ch2 = make(chan bool);
	var ch3 = make(chan bool);
	var ch4 = make(chan bool);
	var ch5 = make(chan bool);

	go phil(1);
	go phil(2);
	go phil(3);
	go phil(4);
	go phil(5);
}


func phil(no int) { 
	name := "philosopher " + string(no);
	forks := 0;
	for{
		select{
			case 
			*
		}	
	}

}

func fork(chan ch) {
	available bool := true;
	for {
		if available {
			ch <- 
		}
	}
	

}
