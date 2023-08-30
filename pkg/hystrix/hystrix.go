package hystrix

import (
	"errors"
	"fmt"
	"sync"

	"github.com/afex/hystrix-go/hystrix"
)

type runFunc func() error

var config = hystrix.CommandConfig{
	Timeout:                5000, // 执行command的超时时间(毫秒)
	MaxConcurrentRequests:  8,    // command的最大并发量
	RequestVolumeThreshold: 5,    // 请求阈值 超过该值进行错误率计算
	SleepWindow:            1000, // 过多长时间，熔断器再次检测是否开启(毫秒)
	ErrorPercentThreshold:  30,   // 错误率
}

type Hystrix struct {
	loadMap        *sync.Map // 储存每个调用函数对应Hystrix
	degradationMsg string    // 降级信息
}

func NewHystrix(msg string) *Hystrix {
	return &Hystrix{
		loadMap:        &sync.Map{},
		degradationMsg: msg,
	}
}

func (h *Hystrix) Run(name string, run func() error) error {
	if _, ok := h.loadMap.Load(name); !ok {
		hystrix.ConfigureCommand(name, config)
		h.loadMap.Store(fmt.Sprintf("hystrix-%s", name), name)
	}

	// Do为同步方式 Go为Goroutine异步
	err := hystrix.Do(name, func() error {
		// other service
		return run()
	}, func(err error) error {
		// fallback message
		return errors.New(h.degradationMsg)
	})
	if err != nil {
		return err
	}
	return nil
}
