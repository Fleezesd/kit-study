package hystrix

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

// hystrix 包是一个延迟和容错库，旨在隔离对远程系统、服务和第 3 方库的访问点，阻止级联故障
func HystrixDemo() {
	output := make(chan bool, 1)

	// goroutine
	errors := hystrix.Go("", func() error {
		// talk to other services
		output <- true
		return nil
	}, func(err error) error {
		// do this when services are down
		return errors.New("xxx services error")
	})

	//  CommandConfig is used to tune circuit settings at runtime
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	// 同步API
	err := hystrix.Do("my_command", func() error {
		// talk to other services
		return nil
	}, nil)
	log.Println(err)

	select {
	case <-output:
		// success
	case <-errors:
		// failure
	}

	// dashboard metrics
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
}
