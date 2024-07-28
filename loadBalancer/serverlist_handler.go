package lb

import (
	"main/server"
	"net/http"
)

func (lb *LoadBalancer) handleAddingServers(w http.ResponseWriter, r *http.Request) {

}

func (lb *LoadBalancer) handleRemovingServers(w http.ResponseWriter, r *http.Request) {

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
