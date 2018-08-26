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
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	gs       *grpc.Server
	address  string
	port     int32
	secure   bool
	services []func(*grpc.Server)

	*GrpcSecurity
}

type GrpcSecurity struct {
	KeyFile  string
	CertFile string
	CaFile   string
}

/*
Instantiate a GRPC server data structure
*/
func NewServer(
	address string,
	port int32,
	certs *GrpcSecurity,
	secure bool,
) *Server {
	server := &Server{
		address:      address,
		port:         port,
		secure:       secure,
		GrpcSecurity: certs,
	}
	return server
}

/*
Start prepares the GRPC server and starts servicing requests
*/
func (s *Server) Start(ctx context.Context) {
	host := strings.Join([]string{
		s.address,
		strconv.Itoa(int(s.port)),
	}, ":")

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if s.secure {
		creds, err := credentials.NewServerTLSFromFile(s.CertFile, s.KeyFile)
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		s.gs = grpc.NewServer(grpc.Creds(creds))

	} else {
		log.Println("In DEFAULT\n")
		s.gs = grpc.NewServer()
	}

	// Register all required services
	for _, service := range s.services {
		service(s.gs)
	}

	if err := s.gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

/*
Stop servicing GRPC requests
*/
func (s *Server) Stop() {
	s.gs.Stop()
}

/*
AddService appends a generic service request function
*/
func (s *Server) AddService(
	registerFunction func(*grpc.Server, interface{}),
	handler interface{},
) {
	s.services = append(s.services, func(gs *grpc.Server) { registerFunction(gs, handler) })
}

func main() {

}
