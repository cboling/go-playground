/*
 * Copyright (c) 2018 - present.  Boling Consulting Solutions (bcsw.net)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"fmt"
	"github.com/cboling/go-playground/grpc/example"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

// While the filename is 'source', think of this as the client application and the
// consumer is a server that does work for us
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := example.NewWorkerClient(conn)

	// Try the unary operation
	unaryOperation(c)

	// Try the server side streaming
	serverStreaming(c)

	// Try the client side streaming
	clientStreaming(c)

	// Try the server side streaming
	biDirectionalStreaming(c)
}

// This does a unary request/response operation
func unaryOperation(client example.WorkerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create request to send
	payload := make([]byte, 32)
	for val := 0; val < len(payload); val++ {
		payload[val] = byte(val)
	}
	now := time.Now()

	request := example.UnaryRequest{
		UtcTimestamp: now.UnixNano(),
		PonId:        0,
		OnuId:        127,
		Payload:      payload,
	}
	log.Printf("Unary Tx: %v", request)

	// Send it and wait for response, that's about it
	response, err := client.RequestUnaryOperation(ctx, &request)
	if err != nil {
		log.Fatalf("could not perform round-trip request/response: %v", err)
	}
	log.Printf("Unary Response: %v", response)
}

// This does a unary request/response operation
func serverStreaming(client example.WorkerClient) {
	// Send one request and get many responses

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var numMessages uint32 = 100
	request := example.ServerRequest{
		PleaseSend: numMessages,
	}
	log.Printf("Server Streaming Tx: want %v responses", numMessages)

	// Send it and wait for responses to stream in
	stream, err := client.RequestServerSideStream(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	count := 0
	log.Println("   Waiting for server responses")
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.RequestServerSideStream failure: %v", client, err)
		}
		count++
		if count%10 == 0 {
			log.Printf("    Got response # %v, %v", count, response)
		}
	}
	log.Printf("Server Streaming completed. Got %v of %v responses",
		count, numMessages)
}

// This does a unary request/response operation
func clientStreaming(client example.WorkerClient) {
	// Send several request and expect one back
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var numMessages uint32 = 100
	log.Println("Streaming client info to server")

	stream, err := client.RequestClientSideStream(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	var count uint32
	for count = 0; count < numMessages; count++ {
		request := example.ClientStreamRequest{
			MsgNumber: count,
		}
		if err := stream.Send(&request); err != nil {
			log.Fatalf("Failed to send client stream request: %v", err)
		}
		if count%10 == 0 {
			log.Printf("   Sent client stream request: %v", count)
		}
	}
	// Clean up and close the stream
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed on final close and rx: %v", err)
	}
	log.Printf("Receive response from client stream: %v", response)
}

// This does a unary request/response operation
func biDirectionalStreaming(client example.WorkerClient) {
	// For this example, we send 10 messages and the server will spit
	// back 3-5 responses.  When we finish, we (client side) will close
	// the connection.  When the server side sees our EOF, it will close
	// it's connection and our go routine for Rx will see that
	log.Println("Starting up bi-directional comms")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the connection
	stream, err := client.BiDirectional(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	// Use a done channel to signal server side has closed
	done := make(chan struct{})

	// Use a go-routine to receive our response stream from the server
	go func(d chan struct{}) {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(d)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a responjse : %v", err)
			}
			log.Printf("   Got a response: %v", response.PrintThisPlease)
		}
	}(done)

	// Now source our client requests
	for count := 0; count < 10; count++ {
		request := example.PeerMessage{
			PrintThisPlease: fmt.Sprintf("Client request %v", count),
		}
		if err := stream.Send(&request); err != nil {
			log.Fatalf("Failed to send client request: %v", err)
		}
	}
	stream.CloseSend()
	<-done

	log.Println("Ending bi-directional comms")
}

//// If here, the simple round trip worked.  Now lets create a
//// request generator that will periodically send out requests
//// to the consumer which will respond on it's own when it is ready
////
//// After a maximum number of requests, a shutdown is sent to have
//// the consumer abort any outstanding responces
//maxMessages := 100
//generator := requestGenerator(maxMessages, 1000*time.Millisecond)
//
//go func(gen <-chan example.UnaryRequest) {
//	for request := range gen {
//		go func(req example.UnaryRequest) {
//			// TODO: Send received messages to consumer
//			time.Sleep(5 * time.Millisecond)
//		}(request)
//	}
//	// TODO: Send the shutdown message
//}(generator)
//
//log.Println("Generator set up and active, starting response listener")
//// Now receive the responses
////responseReceiver := make(chan<- example.ExampleResponse)
//for {
//	// TODO: Receive the responses here  (use a buffered channel)
//	time.Sleep(5 * time.Second)
//	// responseReceiver <-
//}

//
//func requestGenerator(max int, interMessageDelay time.Duration) <-chan example.UnaryRequest {
//	output := make(chan example.UnaryRequest)
//
//	go func(maxRequests int) {
//		for maxRequests > 0 {
//			time.Sleep(interMessageDelay)
//
//			request := example.UnaryRequest{
//				PonId: uint32(rand.Int31n(16)),
//				OnuId: uint32(rand.Int31n(128)),
//			}
//			output <- request
//			maxRequests--
//		}
//		close(output)
//	}(max)
//
//	return output
//}
