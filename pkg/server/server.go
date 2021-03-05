// Package main implements a simple gRPC server that implements the esbbridge rpc server described in esbbridge_rpc.proto
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/spritkopf/esb-bridge/pkg/esbbridge"
	pb "github.com/spritkopf/esb-bridge/pkg/server/service"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type esbBridgeServer struct {
	pb.UnimplementedEsbBridgeServer
}

// GetFeature returns the feature at the given point.
func (s *esbBridgeServer) Transfer(ctx context.Context, msg *pb.EsbMessage) (*pb.EsbMessage, error) {

	// simple echo for now
	log.Printf("Transfer Message: %v\n", msg)
	return &pb.EsbMessage{Addr: msg.Addr, Cmd: msg.Cmd, Payload: msg.Payload}, nil
}

// Listen starts to listen for a specific messages and streams incoming messages to the client
func (s *esbBridgeServer) Listen(listener *pb.Listener, messageStream pb.EsbBridge_ListenServer) error {

	log.Printf("Start listening: %v, %v", listener.Addr, listener.Cmd)

	listenAddr := [5]byte{}
	copy(listenAddr[:5], listener.Addr)

	lc := make(chan esbbridge.EsbMessage, 1)
	esbbridge.AddListener(listenAddr, listener.Cmd[0], lc)

	// TODO: only 3 cycles for testing purpose, use context as abort criterium
	for i := 0; i < 3; i++ {
		msg := <-lc
		log.Printf("Incoming Message: %v\n", msg)
		err := messageStream.Send(&pb.EsbMessage{Addr: msg.Address, Cmd: []byte{msg.Cmd}, Payload: msg.Payload})
		if err != nil {
			return err
		}

	}

	log.Println("Done listening")

	return nil
}

func newServer() *esbBridgeServer {
	s := &esbBridgeServer{}
	return s
}

func main() {
	flag.Parse()

	err := esbbridge.Open("/dev/ttyACM0")
	if err != nil {
		log.Fatalf("Could not open connection to esb-bridge device: %v", err)
	}
	fwVersion, err := esbbridge.GetFwVersion()
	if err != nil {
		log.Fatalf("Error reading Firmware version of esb-bridge device: %v", err)
	}
	log.Printf("esb-bridge firmware version: %v", fwVersion)
	defer esbbridge.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	log.Printf("Serving on port %v\n", *port)
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEsbBridgeServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
