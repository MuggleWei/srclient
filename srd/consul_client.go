package service_discovery

import (
	"fmt"
	"log"
	"time"

	consul "github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	ID           string
	Name         string
	ConsulAgent  *consul.Agent
	ConsulHealth *consul.Health
	TTL          time.Duration
}

// NewConsul returns a Client interface for given consul address
func NewConsulClient(addr string) (ServiceDiscoveryClient, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	consulClient := &ConsulClient{
		ConsulAgent:  client.Agent(),
		ConsulHealth: client.Health(),
		TTL:          time.Second * 5,
	}
	return consulClient, nil
}

///////////////////// implement ServiceDiscoveryClient /////////////////////

func (this *ConsulClient) GetService(service, tag string) ([]*ServiceEntry, error) {
	passingOnly := true
	entries, _, err := this.ConsulHealth.Service(service, tag, passingOnly, nil)
	if len(entries) == 0 && err == nil {
		return nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, err
	}

	var serviceEntries []*ServiceEntry
	for _, entry := range entries {
		serviceEntries = append(serviceEntries, &ServiceEntry{
			Service: entry.Service.Service,
			ID:      entry.Service.ID,
			Addr:    entry.Service.Address,
			Port:    entry.Service.Port,
			Tags:    entry.Service.Tags,
		})
	}

	return serviceEntries, nil
}

func (this *ConsulClient) Register(registration *ServiceRegistration) error {
	reg := &consul.AgentServiceRegistration{
		ID:      registration.ID,
		Name:    registration.Name,
		Address: registration.Addr,
		Port:    registration.Port,
		Tags:    registration.Tag,
		Check: &consul.AgentServiceCheck{
			TTL: registration.TTL.String(),
			DeregisterCriticalServiceAfter: registration.TTL.String(),
		},
	}

	err := this.ConsulAgent.ServiceRegister(reg)
	if err != nil {
		return err
	}

	this.ID = registration.ID
	this.Name = registration.Name
	go this.UpdateTTL(registration.Check)

	return nil
}

func (this *ConsulClient) DeRegister(string) error {
	// TODO:
	return nil
}

///////////////////// self function /////////////////////
func (this *ConsulClient) UpdateTTL(check func() (bool, error)) {
	if check == nil {
		check = func() (bool, error) {
			return true, nil
		}
	}
	ticker := time.NewTicker(this.TTL / 3)
	for range ticker.C {
		this.update(check)
	}
}

func (this *ConsulClient) update(check func() (bool, error)) {
	ok, err := check()
	if !ok {
		agentErr := this.ConsulAgent.FailTTL("service:"+this.ID, err.Error())
		if agentErr != nil {
			log.Print(agentErr)
		}
	} else {
		agentErr := this.ConsulAgent.PassTTL("service:"+this.ID, "")
		if agentErr != nil {
			log.Print(agentErr)
		}
	}
}
