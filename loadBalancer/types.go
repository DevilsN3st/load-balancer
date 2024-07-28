package lb

import (
	"errors"
	"main/server"
	"main/strategy"
	"time"
)

type LoadBalancer struct {
	port                  int
	serverList            []*server.Server
	ipMap                 map[string]bool
	loadBalancingStrategy strategy.Strategy
	ready                 bool
	heartBeatInterval     time.Duration
}

var loadBalancerInstance *LoadBalancer

func load(port int, loadBalancingStrategy strategy.Strategy) *LoadBalancer {
	return &LoadBalancer{
		port:                  port,
		serverList:            make([]*server.Server, 0),
		ipMap:                 make(map[string]bool),
		loadBalancingStrategy: loadBalancingStrategy,
		ready:                 true,
		heartBeatInterval:     5,
	}
}

func GetLoadBalancer(port int, loadBalancingStrategy strategy.Strategy) *LoadBalancer {
	if loadBalancerInstance == nil {
		loadBalancerInstance = load(port, loadBalancingStrategy)
		return loadBalancerInstance
	}
	return loadBalancerInstance
}

func (lb *LoadBalancer) SetServers(nodes []*server.Server) {
	for _, n := range nodes {
		lb.serverList = append(lb.serverList, server.NewServer(n.GetIP(), n.GetHostName(), n.GetPort(), n.GetCpuUtilisation(), n.GetNumOfCurrProcess(), n.GetIsHealthy()))
	}
}

func (lb *LoadBalancer) getServer() (*server.Server, error) {
	c := 0
	for c < len(lb.serverList) {
		server, err := lb.loadBalancingStrategy.GetServer(lb.serverList)
		if err != nil {
			return nil, errors.New("failed :: error occrued in loadbalancer")
		}
		if server.GetIsHealthy() {
			return server, nil
		}
		c += 1
	}
	lb.ready = false
	return nil, errors.New("no server is available to serve request")
}
