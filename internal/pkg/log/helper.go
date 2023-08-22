package log

// 后续会加入viper 读取配置

func LogOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             "debug",
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}
