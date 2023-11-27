package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/galihgpr/proto-client/gen/protobuf"
	pb "github.com/galihgpr/proto-client/gen/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type IUserServiceGrpc interface {
	GreetUser(ctx context.Context, req *pb.GreetingRequest) (res *pb.GreetingResponse, err error)
}

type UserServiceGrpcOption struct {
	Host                  string
	Port                  int
	Timeout               time.Duration
	MaxConcurrentRequests int
	ErrorPercentThreshold int
	Tls                   bool
	PemPath               string
	Secret                string
	Realtime              bool
}

type UserServiceGrpc struct {
	Option     UserServiceGrpcOption
	GrpcClient protobuf.UserServiceClient
}

func NewUserServiceGrpc(opt UserServiceGrpcOption) (iUserService IUserServiceGrpc, err error) {
	var opts []grpc.DialOption

	if opt.Tls {
		var pemServerCA []byte
		pemServerCA, err = ioutil.ReadFile(opt.PemPath)
		if err != nil {
			return
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemServerCA) {
			err = errors.New("failed to add server ca's certificate")
			return
		}

		// Create the credentials and return it
		config := &tls.Config{
			RootCAs: certPool,
		}

		tlsCredentials := credentials.NewTLS(config)

		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	opts = append(opts, grpc.WithReturnConnectionError())
	opts = append(opts, grpc.FailOnNonTempDialError(true))

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 8080),
		opts...,
	)
	if err != nil {
		return
	}

	gRPCClient := pb.NewUserServiceClient(conn)

	iUserService = UserServiceGrpc{
		Option:     opt,
		GrpcClient: gRPCClient,
	}
	return
}

func (o UserServiceGrpc) GreetUser(ctx context.Context, req *pb.GreetingRequest) (res *pb.GreetingResponse, err error) {
	res, err = o.GrpcClient.GreetUser(context.TODO(), req)
	return
}
