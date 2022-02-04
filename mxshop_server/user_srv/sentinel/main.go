package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
	"math/rand"
	"time"
)

func main() {
	err := sentinel.InitDefault()
	if err != nil {
		// 初始化 Sentinel 失败
		log.Fatalf("初始化setinel失败: %v", err)
		return
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.WarmUp, //冷启动
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              1000, //1k并发
			WarmUpPeriodSec:        30, //30s内达到Threshold(1k)并发
		},
		{
			Resource:               "some-test2",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})

	if err != nil {
		// 加载规则失败，进行相关处理
		log.Fatalf("加载规则失败: %v", err)
		return
	}

	var gloabTotal int //总数
	var passTotal int  //通过数
	var blockTotal int //拒绝数

	ch := make(chan interface{})
	for i := 0; i < 100; i++ {
		go func() {
			for {
				gloabTotal++
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					//fmt.Println("限流中")
					blockTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					//fmt.Println("通过")
					passTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}
			}
		}()
	}

	go func() {
		var oldTotal int //过去1s总数
		var oldPass int  //过去1s通过数
		var oldBlock int //过去1s拒绝数
		for {
			oneSecondTotal := gloabTotal - oldTotal
			oldTotal = gloabTotal

			onePassTotal := passTotal - oldPass
			oldPass = passTotal

			oneBlockTotal := blockTotal - oldBlock
			oldBlock = blockTotal

			time.Sleep(time.Second)
			fmt.Printf("total: %d, pass: %d, block: %d\n", oneSecondTotal, onePassTotal, oneBlockTotal)
		}
	}()

	<-ch

	//for i := 0; i <= 12; i++ {
	//	e, b := sentinel.Entry("some-test2", sentinel.WithTrafficType(base.Inbound))
	//	if b != nil {
	//		fmt.Println("限流中")
	//	} else {
	//		fmt.Println("通过")
	//		e.Exit()
	//	}
	//}

}
