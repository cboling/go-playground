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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"github.com/cboling/go-playground/grpc/example"
)

const (
	address = "localhost:50051"
)

// server is used to implement example.WorkerServer
type server struct{}

// RequestUnaryOperation implements example.WorkerServer unary response operation
func (s *server) RequestUnaryOperation(ctx context.Context, in *example.ExampleRequest) (*example.ExampleResponse, error) {

	// Swap the data in the payload and send it back

	myData := make([]byte, len(in.Payload))

	for val := 0; val < len(myData); val++ {
		myData[val] = in.Payload[len(myData)-val-1]
	}
	response := example.ExampleResponse{
		UtcTimestamp: in.UtcTimestamp,
		PonId:        in.PonId,
		OnuId:        in.OnuId,
		Payload:      myData,
	}
	return &response, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	example.RegisterWorkerServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
