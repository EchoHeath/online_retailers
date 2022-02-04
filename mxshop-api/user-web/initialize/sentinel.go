package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func Sentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		// 初始化 Sentinel 失败
		zap.S().Fatalf("初始化setinel失败: %v", err)
		return
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test2",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              100,
			StatIntervalInMs:       1000,
		},
	})

	if err != nil {
		// 加载规则失败，进行相关处理
		zap.S().Fatalf("加载规则失败: %v", err)
		return
	}
}
