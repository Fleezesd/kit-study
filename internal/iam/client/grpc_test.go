package client

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/fleezesd/kit-study/internal/iam/service"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	uuid "github.com/satori/go.uuid"
)

func TestGRPCClient(t *testing.T) {
	serviceAddress := "localhost:8081"
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("connect error")
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

func BenchmarkGRPCClient(b *testing.B) {

}

func ExampleTestGRPCClient() {

}
