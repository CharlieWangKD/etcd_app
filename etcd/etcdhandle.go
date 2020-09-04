package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"log"
	"os"
	"time"
)

var leaderFlag bool

func (c EtcdC) LeaderFunc() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	for {
		select {
		case <-ticker.C:
			{
				if leaderFlag == true {
					log.Println("I amd leader. yeah!")
					// TODO: do something
					//return
				} else {
					log.Println("[info] Service is not master")
				}
			}
		}
	}
}

func (c EtcdC) Campaign() {
	log.Println("start campaign goroutine")
	for {
		s, err := concurrency.NewSession(c.EtcdClient, concurrency.WithTTL(15))
		if err != nil {
			log.Println(err)
			continue
		}
		e := concurrency.NewElection(s, c.LeaderKey)
		if err = e.Campaign(context.TODO(), c.LocalIp); err != nil {
			log.Println(err)
			continue
		}
		log.Println("elect: success")
		leaderFlag = true
		select {
		case <-s.Done():
			leaderFlag = false
			log.Println("elect: expired")
		}
	}
}

func (c EtcdC) Watch() {
	log.Println("start watch goroutine")
	rec := c.EtcdClient.Watch(context.Background(), c.WorkDir, clientv3.WithPrefix())
	for wresp := range rec {
		for _, ev := range wresp.Events {
			// TODO: do something
			log.Println(ev)
			log.Println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}

func (c EtcdC) Register() {
	log.Println("start register goroutine")
	lease := clientv3.NewLease(c.EtcdClient)
	leaseid, err := lease.Grant(context.Background(), 600)
	if err != nil {
		log.Println(err)
	}
	v, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	k := c.RegisterDir + c.LocalIp
	_, err = c.EtcdClient.Put(context.Background(), k, v, clientv3.WithLease(leaseid.ID))
	if err != nil {
		log.Println(err)
	}
	keepRespChan, err := lease.KeepAlive(context.TODO(), leaseid.ID)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepRespChan == nil {
					log.Println("lease already expired.")
					goto END
				} else {
					log.Println("receive auto renew lease.:", keepResp.ID)
				}
			}
		}
	END:
	}()
}

func (c EtcdC) Get(key string) {
	cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := c.EtcdClient.Get(cctx, key)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
	cancel()
}

func (c EtcdC) Init() {
	cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := c.EtcdClient.Get(cctx, c.WorkDir, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	}
	cancel()
	for _, ev := range resp.Kvs {
		// TODO: run init func
		log.Println("key: ", string(ev.Key))
		log.Println("value: ", string(ev.Value))
	}
	log.Println("init done.")
}
