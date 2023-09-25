a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?
    We donut use metadata, freaks

b) Does your implementation use threads or processes? Why is it not realistic to use threads?
   We use processes. For the simulation of the TCP these are sufficient enough to demonstrate how multiple 
   tasks can be managed concuretly. To create a more realalistic simulation we would need to use two 
   machines and actual network sockets - however this is a complex task.

c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
    We  observe the sequence numbers. If one is missing - we would simply send a request to resend the data. 

d) In case messages can be delayed or lost, how does your implementation handle message loss?


e) Why is the 3-way handshake important?
   The handshake is important to ensure both parties are ready to transfer data. Because this is a full duplex, we must
   ensure that both client and server is ready - we do this by requesting and acknowleding. The other thing we need to do
   is star the sequence for both client and server. 
   Clint sends request (sequenceStart)-> Server

   Connection is succesful
   Start the sequence 

