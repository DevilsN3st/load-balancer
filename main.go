package main

import (
	"flag"
	"log"
	lb "main/loadBalancer"
	"main/server"
	"main/strategy"
	"strings"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	PORT = flag.Int("p", 8080, "Load balancer PORT")
)

// hardcoded for now otherwise we could create a separate fetch from separate endpoint
func getServerListSlice() []*server.Server {
	var serverList []*server.Server
	serverList = append(serverList, server.NewServer("http://localhost", "my-pc", 3000, 0, 0, true))
	serverList = append(serverList, server.NewServer("http://localhost", "my-pc", 3001, 0, 0, true))
	serverList = append(serverList, server.NewServer("http://localhost", "my-pc", 3002, 0, 0, true))
	serverList = append(serverList, server.NewServer("http://localhost", "my-pc", 3003, 0, 0, true))
	return serverList
}

func main() {
	var nodes stringSlice
	flag.Var(&nodes, "n", "List of servers to balance load")
	flag.Parse()

	loadBalancer := lb.GetLoadBalancer(*PORT, strategy.NewRoundRobinStrategy())
	serverListSlice := getServerListSlice()
	loadBalancer.SetServers(serverListSlice)

	if err := loadBalancer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
