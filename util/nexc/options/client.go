package options

import (
	"fmt"
	"log"

	"gitlab.com/mergetb/tech/nex/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func dial() (*grpc.ClientConn, nex.NexClient) {

	conn, err := grpc.Dial(fmt.Sprintf("%s", Server), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return conn, nex.NewNexClient(conn)
}

func withClient(f func(nex.NexClient) error) error {

	conn, cli := dial()
	defer conn.Close()

	return f(cli)

}

func grpcFatal(err error) {

	s, ok := status.FromError(err)
	if !ok {
		log.Fatal(err)
	}
	log.Fatal(s.Message())

}
