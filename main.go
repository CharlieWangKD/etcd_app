package main

import (
	"etcd_app/etcd"
	"etcd_app/g"
	"flag"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")

	g.ParseConfig(*cfg)

	client := etcd.GetIns()

	client.Init()

	go client.Campaign()

	go client.Register()

	go client.Watch()

	go client.LeaderFunc()

	select {}
}
