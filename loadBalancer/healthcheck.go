package lb

import (
	"sync"
	"time"
)

func (lb *LoadBalancer) healthCheck() {
	wg := &sync.WaitGroup{}
	for {
		for _, server := range lb.serverList {
			wg.Add(1)
			go server.HealthCheck(wg)
		}
		Inactive := true
		for _, server := range lb.serverList {
			if server.GetIsHealthy() {
				Inactive = false
				lb.ready = true
				break
			}
		}
		if Inactive {
			lb.ready = false
		}
		wg.Wait()
		time.Sleep(lb.heartBeatInterval * time.Second)
	}
}
