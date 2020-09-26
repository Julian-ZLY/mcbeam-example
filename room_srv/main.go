package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/store/redis/v2"
	"github.com/wolfplus2048/mcbeam-example/room_srv/base"
	"github.com/wolfplus2048/mcbeam-example/room_srv/logic"
	"github.com/wolfplus2048/mcbeam-example/room_srv/mgr"
	"github.com/wolfplus2048/mcbeam-plus"
)

func main() {
	logger.Init(logger.WithLevel(logger.DebugLevel))
	service := mcbeam.NewService(
		mcbeam.Name("room"),
		mcbeam.Registry(etcd.NewRegistry()),
		mcbeam.MicroService(
			micro.NewService(
				micro.Store(redis.NewStore()),
				micro.Broker(nats.NewBroker()),
				micro.WrapHandler(mgr.WrapSession()),
			),
		),
	)
	if err := service.Init(); err != nil {
		logger.Fatal(err)
	}
	base.Start()
	service.Register(mgr.NewHandler(logic.NewMJPlayer, logic.NewMJRoom, service.Options().Service))
	service.Register(&logic.MJHandler{})

	s := service.Options().Service.Server()
	s.Subscribe(s.NewSubscriber("session.close", &mgr.Sub{}))
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
	base.Stop()
}
