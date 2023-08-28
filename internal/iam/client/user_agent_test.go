package client

import (
	"context"
	"os"
	"testing"
	"time"

	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	"github.com/go-kit/log"
)

func TestUserAgentClient(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	client, err := NewUserAgentClient([]string{"localhost:2379"}, logger)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 1; i++ {
		time.Sleep(time.Second)
		userAgent, err := client.UserAgentClient()
		if err != nil {
			logger.Log(err)
			t.Error(err)
			return
		}
		rsp, err := userAgent.Login(context.Background(), &pb.LoginRequest{
			Username: "admin",
			Password: "xxx",
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(rsp)
	}

}
