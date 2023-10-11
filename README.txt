a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?
   Our packages are slices. The idea is to send slices of length 2 containing sequence and acknowledgment.

b) Does your implementation use threads or processes? Why is it not realistic to use threads?
   We use Threads using goroutines. For the simulation of the TCP these are sufficient enough to demonstrate how multiple 
   tasks can be managed concuretly. To create a more realalistic simulation we would need to use two 
   machines and actual network sockets - however this is a complex task. Threads are unrealisic because there will not be an issue
   with loosing data, corruption or reordering.

c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
    We would observe the sequence numbers, and use these to order them correctly.

d) In case messages can be delayed or lost, how does your implementation handle message loss?
   Our implemeenation is only of the threeway handshake (Only establishing the connection, not transmitting data), however if a request to the server or resonse back to the client is delayed, we take acount for that by 
   having a timeout of 2 secounds. If the timeout exprires then it is considered it lost and then one could implent a restart of the three way handsshake again.

   The same strategy could be used when transmitting data, but when the timeout exprires the one sending data, would resend it agin. And by using sequence number we could check for
   which data that wasn't recieved.

e) Why is the 3-way handshake important?
   The handshake is important to ensure both parties are ready to transfer data. Because this is a full duplex, we must
   ensure that both client and server is ready - we do this by requesting and acknowleding. The other thing we need to do
   is star the sequence for both client and server. 
   This way both members knows traffic is recieved at the other end. 
   And by this we ensure that TCP is a connection-oriented, reliable, streaming-protocol