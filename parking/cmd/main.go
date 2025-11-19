package main

import (
	"github.com/h4x4d/parking_net/parking/cmd/grpc"
	"github.com/h4x4d/parking_net/parking/cmd/rest"
	"sync"
)

func main() {
	group := sync.WaitGroup{}
	group.Add(2)
	go rest.StartServer(&group)
	go grpc.StartServer(&group)

	group.Wait()
}
