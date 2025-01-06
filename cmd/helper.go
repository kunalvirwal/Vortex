package main

import (
	"github.com/kunalvirwal/Vortex/internal/master"
)

func initSchedulers() {
	// checks for docker events: done
	// implement a desiered state
	// write container functions
	// implement custom health checks: done
	// implement gRPC routes

}

func InitTracker() {
	go master.ListenEvents()
}

func InitgRPCListener() {
	go master.StartGrpcServer()
}
