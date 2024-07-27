package server

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	ip               string
	hostName         string
	port             int
	cpuUtilisation   float32
	numOfCurrProcess int
	healthy          bool
}

func NewServer(ip string, hostName string, port int, cpuUtilisation float32, numOfCurrProcess int, healthy bool) *Server {
	return &Server{
		ip:               ip,
		hostName:         hostName,
		port:             port,
		cpuUtilisation:   cpuUtilisation,
		numOfCurrProcess: numOfCurrProcess,
		healthy:          healthy,
	}
}

type HTTPMethod string

// getter
func (s *Server) GetIP() string {
	return s.ip
}
func (s *Server) GetHostName() string {
	return s.hostName
}
func (s *Server) GetPort() int {
	return s.port
}
func (s *Server) GetCpuUtilisation() float32 {
	return s.cpuUtilisation
}
func (s *Server) GetNumOfCurrProcess() int {
	return s.numOfCurrProcess
}
func (s *Server) GetIsHealthy() bool {
	return s.healthy
}

// setter
func (s *Server) SetIP(ip string) {
	s.ip = ip
}
func (s *Server) SetHostName(hostName string) {
	s.hostName = hostName
}
func (s *Server) SetPort(port int) {
	s.port = port
}
func (s *Server) SetCpuUtilisation(cpuUtilisation float32) {
	s.cpuUtilisation = cpuUtilisation
}
func (s *Server) SetNumOfCurrProcess(numOfCurrProcess int) {
	s.numOfCurrProcess = numOfCurrProcess
}
func (s *Server) SetIsHealthyTrue() {
	s.healthy = true
}
func (s *Server) SetIsHealthyFalse() {
	s.healthy = false
}

// HeatlhCheck
func (s *Server) HealthCheck(wg *sync.WaitGroup) {
	defer wg.Done()
	portString := strconv.Itoa(s.port)
	url := s.ip + ":" + portString + "/health"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("error in http.NewRequest", err)
		return
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		s.SetIsHealthyFalse()
		log.Println("error in healthcheck of server", err)
		return
	}
	defer res.Body.Close()

	if _, err := io.ReadAll(res.Body); err != nil {
		log.Println("error in reading response", err)
	}

	if res.StatusCode == http.StatusOK {
		s.SetIsHealthyTrue()
		return
	}
	s.SetIsHealthyFalse()
}
