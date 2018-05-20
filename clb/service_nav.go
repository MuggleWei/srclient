package client_loadbalancer

import (
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	srd "github.com/MuggleWei/go-toy/srd"
)

type ServiceNavigation struct {
	ServiceDiscoveryClient srd.ServiceDiscoveryClient

	ServiceName string
	Services    []*srd.ServiceEntry
	RollIndex   uint32
	Mtx         sync.RWMutex
}

func NewServiceNavigation(serviceDiscoveryClient srd.ServiceDiscoveryClient, serviceName string, interval time.Duration) (*ServiceNavigation, error) {
	nav := &ServiceNavigation{
		ServiceDiscoveryClient: serviceDiscoveryClient,
		ServiceName:            serviceName,
		Services:               []*srd.ServiceEntry{},
		RollIndex:              0,
	}

	err := nav.UpdateServices()
	if err != nil {
		return nil, err
	}
	go nav.update(interval)

	return nav, nil
}

func (this *ServiceNavigation) GetService() string {
	this.Mtx.RLock()
	defer this.Mtx.RUnlock()

	if len(this.Services) == 0 {
		return ""
	}

	idx := this.getIndex(uint32(len(this.Services)))
	return this.Services[idx].Addr + ":" + strconv.Itoa(this.Services[idx].Port)
}

func (this *ServiceNavigation) UpdateServices() error {
	services, err := this.ServiceDiscoveryClient.GetService(this.ServiceName, "")
	if err != nil {
		log.Print(err)
		return err
	}

	this.Mtx.Lock()
	defer this.Mtx.Unlock()

	this.Services = services

	return nil
}

func (this *ServiceNavigation) update(interval time.Duration) {
	if interval <= 0 {
		interval = time.Second * 3
	}
	ticker := time.NewTicker(interval)
	for range ticker.C {
		this.UpdateServices()
	}
}

func (this *ServiceNavigation) getIndex(module uint32) uint32 {
	for {
		old := this.RollIndex
		idx := (old + 1) % module
		if atomic.CompareAndSwapUint32(&this.RollIndex, old, idx) {
			return idx
		}
	}
}
