package client

import (
	"context"
	"flag"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/fleezesd/kit-study/internal/iam/service"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	uuid "github.com/satori/go.uuid"
)

var (
	addr = flag.String("addr", "localhost:9090", "The address to connect to.")
)

func TestGRPCClient(t *testing.T) {
	serviceAddress := "localhost:9090"
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	userClient := pb.NewUserClient(conn)
	UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
	md := metadata.Pairs(service.ContextReqUUid, UUID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	res, err := userClient.RpcUserLogin(ctx, &pb.LoginRequest{
		Username: "admin",
		Password: "xxx",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Token)
}
