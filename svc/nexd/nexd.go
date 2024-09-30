package nexd

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	nex "gitlab.com/mergetb/tech/nex/pkg"
	"google.golang.org/grpc/reflection"
)

func (d *NexD) Run() {

	log.Printf("nexd %s\n", nex.Version)
	log.SetLevel(log.DebugLevel)

    nex.LoadConfig()
	go nex.RunLeaseManager()

	grpcServer := grpc.NewServer()
	nex.RegisterNexServer(grpcServer,d)
	reflection.Register(grpcServer)
	// err := nex.LoadConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	l, err := net.Listen("tcp", Listen)
	if err != nil {
		log.Errorf("failed to listen: %s", err.Error())
		return
	}

	log.Infof("Listening on tcp://%s", Listen)
	grpcServer.Serve(l)

}

/*** Interfaces *******************/
