package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	utils "github.com/fleezesd/kit-study/pkg/hystrix"
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
	// 调用服务名称
	serviceName := "login"
	// 服务熔断
	hy := utils.NewHystrix("调用服务降级")
	circuit, _, _ := hystrix.GetCircuit(serviceName)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		userAgent, err := client.UserAgentClient()
		if err != nil {
			logger.Log(err)
			t.Error(err)
			return
		}
		// 调用封装好的hystrix
		err = hy.Run(serviceName, func() error {
			rsp, err := userAgent.Login(context.Background(), &pb.LoginRequest{
				Username: "admin",
				Password: "xxx",
			})
			if err != nil {
				t.Error(err)
				return err
			}
			t.Log(rsp)
			return nil
		})
		t.Log("熔断器开启状态:", circuit.IsOpen(), "请求是否允许：", circuit.AllowRequest())
		if err != nil {
			t.Log(err)
		}
	}

}
