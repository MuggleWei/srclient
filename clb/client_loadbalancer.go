package client_loadbalancer

import (
	"sync"
	"time"

	srd "github.com/MuggleWei/go-toy/srd"
)

type ClientLoadBalancer struct {
	ServiceDiscoveryClient srd.ServiceDiscoveryClient

	Mtx            sync.RWMutex
	ServiceNavs    map[string]*ServiceNavigation
	UpdateInterval time.Duration
}

func NewClientLoadBalancer(serviceDiscoveryClient srd.ServiceDiscoveryClient, interval time.Duration) *ClientLoadBalancer {
	return &ClientLoadBalancer{
		ServiceDiscoveryClient: serviceDiscoveryClient,
		ServiceNavs:            make(map[string]*ServiceNavigation),
		UpdateInterval:         interval,
	}
}

func (this *ClientLoadBalancer) GetService(service string) (string, error) {
	nav := this.getServiceNav(service)
	if nav == nil {
		newNav, err := this.newServiceNav(service)
		if err != nil {
			return "", err
		}
		nav = newNav
	}

	return nav.GetService(), nil
}

func (this *ClientLoadBalancer) getServiceNav(service string) *ServiceNavigation {
	this.Mtx.RLock()
	defer this.Mtx.RUnlock()

	nav, ok := this.ServiceNavs[service]
	if !ok {
		return nil
	}
	return nav
}

func (this *ClientLoadBalancer) newServiceNav(service string) (*ServiceNavigation, error) {
	this.Mtx.Lock()
	defer this.Mtx.Unlock()

	_, ok := this.ServiceNavs[service]
	if !ok {
		nav, err := NewServiceNavigation(this.ServiceDiscoveryClient, service, this.UpdateInterval)
		if err != nil {
			return nil, err
		}

		this.ServiceNavs[service] = nav
	}

	return this.ServiceNavs[service], nil

}
