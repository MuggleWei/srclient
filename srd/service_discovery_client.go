package service_discovery

import "time"

type ServiceDiscoveryClient interface {
	// Get a Service from consul
	GetService(string, string) ([]*ServiceEntry, error)

	// Register a service with local agent
	Register(*ServiceRegistration) error

	// Deregister a service with local agent
	DeRegister(string) error
}

type ServiceRegistration struct {
	ID    string
	Name  string
	Addr  string
	Port  int
	Tag   []string
	TTL   time.Duration
	Check func() (bool, error)
}

type ServiceEntry struct {
	Service string
	ID      string
	Addr    string
	Port    int
	Tags    []string
}
