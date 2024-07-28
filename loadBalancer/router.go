package lb

import (
	"fmt"
	"log"
	"net/http"
)

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
	http.HandleFunc("/addServer", lb.handleAddingServers)
	http.HandleFunc("/removeServer", lb.handleRemovingServers)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", lb.handleConnection)
	go lb.healthCheck()
	log.Println("Load balancer is listening on", lb.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", lb.port), nil)
}
