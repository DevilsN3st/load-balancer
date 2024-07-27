package strategy

import "main/server"

type Strategy interface {
	GetServer(serverList []*server.Server) (*server.Server, error)
}
