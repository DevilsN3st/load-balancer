package lb

import (
	"errors"
	"fmt"
	"io"
	"log"
	"main/server"
	"main/strategy"
	"net/http"
	"strconv"
	"sync"
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

func (lb *LoadBalancer) addServer(server *server.Server) bool {
	serverIp := server.GetIP()
	if _, ok := lb.ipMap[serverIp]; ok {
		return false
	}
	lb.serverList = append(lb.serverList, server)
	lb.ipMap[serverIp] = true
	return false
}

func (lb *LoadBalancer) removeServer(serverToRemove *server.Server) bool {
	serverIp := serverToRemove.GetIP()
	if _, ok := lb.ipMap[serverIp]; ok {
		lb.serverList = removeServerFromList(lb.serverList, serverToRemove)
		lb.ipMap[serverIp] = true
		return true
	}
	return false
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

func removeServerFromList(serverList []*server.Server, serverToRemove *server.Server) []*server.Server {
	var newServerList []*server.Server
	for _, curr_server := range serverList {
		if curr_server.GetIP() == serverToRemove.GetIP() {
			continue
		}
		newServerList = append(newServerList, curr_server)
	}
	return newServerList
}

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

func (lb *LoadBalancer) ListenAndServe() error {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		status := "healthy"
		if !lb.ready {
			status = "unhealthy"
		}
		fmt.Print(lb.ready)
		fmt.Fprintf(w, "Load balancer is %s", status)
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", lb.handleConnection)
	go lb.healthCheck()
	log.Println("Load balancer is listening on", lb.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", lb.port), nil)
}

func (lb *LoadBalancer) handleConnection(w http.ResponseWriter, r *http.Request) {

	if !lb.ready {
		http.Error(w, "Load balancer is offline", http.StatusServiceUnavailable)
		return
	}
	server, err := lb.getServer()
	fmt.Println("server", server)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	portString := strconv.Itoa(server.GetPort())
	req, err := http.NewRequest(r.Method, server.GetIP()+":"+portString+r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		log.Println("error in http.NewRequest", err)
		return
	}
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 3 * time.Second,
		},
	}
	log.Println("req at redirecting request", req)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		log.Println("error in client.Do", err.Error())
		return
	}
	defer resp.Body.Close()
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to write response body", http.StatusInternalServerError)
		log.Println("error in io.Copy", err)
		return
	}
}
