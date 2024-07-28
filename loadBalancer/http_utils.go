package lb

import (
	"errors"
	"fmt"
	"io"
	"log"
	"main/server"
	"net/http"
	"strconv"
	"time"
)

func (lb *LoadBalancer) prepareOutgoingHttpRequest(w http.ResponseWriter, r *http.Request, server *server.Server) (*http.Request, *http.Client, error) {
	portString := strconv.Itoa(server.GetPort())
	url := server.GetIP() + ":" + portString + r.URL.String()
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return nil, nil, errors.New(fmt.Sprintf("failed to create request: %d", http.StatusInternalServerError))
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
	return req, client, nil
}

func (lb *LoadBalancer) handleHttpCalls(req *http.Request, client *http.Client, w http.ResponseWriter) {

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
