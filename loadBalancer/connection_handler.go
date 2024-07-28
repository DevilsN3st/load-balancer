package lb

import (
	"log"
	"net/http"
)

func (lb *LoadBalancer) handleConnection(w http.ResponseWriter, r *http.Request) {

	if !lb.ready {
		http.Error(w, "Load balancer is offline", http.StatusServiceUnavailable)
		return
	}

	server, err := lb.getServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	req, client, err := lb.prepareOutgoingHttpRequest(w, r, server)
	if err != nil {
		log.Println("error in http.NewRequest", err)
		return
	}

	lb.handleHttpCalls(req, client, w)

}
