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
	"log"
	"math/rand"
	"time"
)

const (
	port = ":50051"
)

// This does a unary request/response operation
func unaryOperation(client example.WorkerClient) {
	fmt.Println(client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create request to send
	payload := make([]byte, 32)
	for val := 0; val < len(payload); val++ {
		payload[val] = byte(val)
	}
	now := time.Now()

	request := example.ExampleRequest{
		UtcTimestamp: now.UnixNano(),
		PonId:        0,
		OnuId:        127,
		Payload:      payload,
	}
	// Send it and wait for response
	response, err := client.RequestUnaryOperation(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(response)
}

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

	// If here, the simple round trip worked.  Now lets create a
	// request generator that will periodically send out requests
	// to the consumer which will respond on it's own when it is ready
	//
	// After a maximum number of requests, a shutdown is sent to have
	// the consumer abort any outstanding responces
	maxMessages := 100
	generator := requestGenerator(maxMessages, 1000*time.Millisecond)

	go func(gen chan<- example.ExampleRequest) {
		for request := range gen {
			go func(req example.ExampleRequest) {
				// TODO: Send received messages to consumer
				time.Sleep(5 * time.Millisecond)
			}(request)
		}
		// TODO: Send the shutdown message
	}(generator)

	// Now receive the responses
	//responseReceiver := make(chan<- example.ExampleResponse)
	for {
		// TODO: Receive the responses here  (use a buffered channel)
		time.Sleep(5 * time.Second)
		// responseReceiver <-
	}
}

func requestGenerator(max int, interMessageDelay time.Duration) chan<- example.ExampleRequest {
	output := make(chan example.ExampleRequest)

	go func(maxRequests int) {
		for maxRequests > 0 {
			time.Sleep(interMessageDelay)

			request := example.ExampleRequest{
				PonId: uint32(rand.Int31n(16)),
				OnuId: uint32(rand.Int31n(128)),
			}
			output <- request
			maxRequests--
		}
		close(output)
	}(max)

	return output
}
