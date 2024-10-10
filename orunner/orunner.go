package orunner

import "OCluster_runner/grpc/src"

type Orunner struct {
}

func (o *Orunner) Run() {
	src.RunGRPC()
}
