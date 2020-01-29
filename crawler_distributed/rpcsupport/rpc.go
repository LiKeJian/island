package rpcsupport

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"reptiles/crawler_distributed/config"
	pb "reptiles/crawler_distributed/proto"
)

func ServeRpc(
	host string, service pb.ReptilesServer) error {

	s := grpc.NewServer()
	pb.RegisterReptilesServer(s, service)

	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", host)

	if err := s.Serve(listener);err != nil{
		log.Fatalf("failed to listen: %v", err)
	}
	return nil
}

func NewClient(port string) (pb.ReptilesClient, error) {
	conn, err := grpc.Dial(config.Host+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewReptilesClient(conn)
	return client, nil
}
