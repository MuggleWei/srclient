package srd

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
	AgentReg     *consul.AgentServiceRegistration
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
	this.AgentReg = &consul.AgentServiceRegistration{
		ID:      registration.ID,
		Name:    registration.Name,
		Address: registration.Addr,
		Port:    registration.Port,
		Tags:    registration.Tag,
		Check: &consul.AgentServiceCheck{
			TTL:                            registration.TTL.String(),
			DeregisterCriticalServiceAfter: registration.TTL.String(),
		},
	}

	err := this.ConsulAgent.ServiceRegister(this.AgentReg)
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

// /////////////////// self function /////////////////////
func (this *ConsulClient) UpdateTTL(check func() (bool, error)) {
	if check == nil {
		check = func() (bool, error) {
			return true, nil
		}
	}

	ticker := time.NewTicker(this.TTL / 3)
	var err error
	for range ticker.C {
		if err == nil {
			err = this.update(check)
			if err != nil {
				log.Printf("ttl update failed: %v\n", err)
			}
		} else {
			err = this.ConsulAgent.ServiceRegister(this.AgentReg)
			if err != nil {
				log.Printf("try register failed: %v\n", err)
			} else {
				log.Printf("success re-register\n")
			}
		}
	}
}

func (this *ConsulClient) update(check func() (bool, error)) error {
	ok, err := check()

	var agentErr error
	if !ok {
		agentErr = this.ConsulAgent.FailTTL("service:"+this.ID, err.Error())
	} else {
		agentErr = this.ConsulAgent.PassTTL("service:"+this.ID, "")
	}

	if agentErr != nil {
		log.Printf("consul client send ttl failed: %v", agentErr)
	}

	return agentErr
}
