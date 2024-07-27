package strategy

import (
	"errors"
	"fmt"
	"main/server"
)

type roundRobinStrategy struct {
	currIndex int
}

func NewRoundRobinStrategy() Strategy {
	return &roundRobinStrategy{
		currIndex: 0,
	}
}

func (r *roundRobinStrategy) GetServer(serverList []*server.Server) (*server.Server, error) {
	if len(serverList) == 0 {
		return nil, errors.New("list of servers connected :: nil")
	}
	totalLen := len(serverList)
	r.currIndex = (r.currIndex + 1) % totalLen
	fmt.Println("current Index at which server is making call ::::::::::::::::", serverList, r.currIndex)
	return serverList[r.currIndex], nil
}
