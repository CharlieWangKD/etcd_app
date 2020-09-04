package etcd

import (
	"etcd_app/g"
	"github.com/coreos/etcd/pkg/transport"
	"go.etcd.io/etcd/clientv3"
	"net"
	"sync"
	"time"
)

var (
	mu         sync.Mutex
	EtcdClient *EtcdC
)

type EtcdC struct {
	EtcdClient  *clientv3.Client
	LocalIp     string
	WorkDir     string
	RegisterDir string
	LeaderKey   string
}

func GetIns() *EtcdC {
	if EtcdClient == nil {
		mu.Lock()
		defer mu.Unlock()
		if EtcdClient == nil {
			EtcdClient = &EtcdC{
				EtcdClient:  InitClient(),
				LocalIp:     InitLocalIp(),
				WorkDir:     g.Config().WorkDir,
				RegisterDir: g.Config().RegisterDir,
				LeaderKey:   g.Config().LeaderKey,
			}
		}
	}
	return EtcdClient
}

func InitTlsClient() *clientv3.Client {
	tlsInfo := transport.TLSInfo{
		CertFile:      `crt/etcd_client.pem`,
		KeyFile:       `crt/etcd_client-key.pem`,
		TrustedCAFile: `crt/ca.pem`,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		panic(err.Error())
		return nil
	}
	c, err := clientv3.New(clientv3.Config{
		Endpoints:            g.Config().EtcdCluster,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    40 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
		TLS:                  tlsConfig,
	})
	if err != nil {
		panic(err.Error())
		return nil
	}
	return c
}

func InitClient() *clientv3.Client {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:            g.Config().EtcdCluster,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    40 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
		return nil
	}
	return c
}

func InitLocalIp() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				return ip
			}
		}
	}
	return ""
}
