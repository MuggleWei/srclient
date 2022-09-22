package clb

import (
	"errors"
	"log"
	"testing"
	"time"

	srd "github.com/MuggleWei/srclient/srd"
)

var consulAddr = "127.0.0.1:8500"
var service = "hello-service"

func TestNav(t *testing.T) {
	serviceDiscoveryClient, err := srd.NewConsulClient(consulAddr)
	if err != nil {
		panic(err)
	}

	nav, err := NewServiceNavigation(serviceDiscoveryClient, service, time.Second*3)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		addr := nav.GetService()
		if addr == "" {
			panic(errors.New("failed get service"))
		}
		log.Println(addr)
	}
}

func TestLB(t *testing.T) {
	serviceDiscoveryClient, err := srd.NewConsulClient(consulAddr)
	if err != nil {
		panic(err)
	}

	clientLB := NewClientLoadBalancer(serviceDiscoveryClient, time.Second*3)

	for i := 0; i < 10; i++ {
		addr, err := clientLB.GetService(service)
		if err != nil {
			panic(err)
		}

		log.Println(addr)
	}
}
