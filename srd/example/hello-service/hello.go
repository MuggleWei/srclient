package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	srd "github.com/MuggleWei/go-toy/srd"
)

const (
	ServiceStatus_Pass           = 1
	ServiceStatus_ReadyToOffline = 2
)

var serviceStatus = ServiceStatus_Pass

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	ptrIp := flag.String("ip", "127.0.0.1", "run ip address")
	ptrPort := flag.Int("port", 10102, "listen port")
	ptrServiceName := flag.String("service.name", "echo-service", "service name")
	ptrServiceID := flag.String("service.id", "echo-service-0", "id in service group")
	ptrServiceTag := flag.String("service.tag", "", "service tags")
	ptrConsulAddr := flag.String("consul", "127.0.0.1:8500", "consul address")

	flag.Parse()

	serviceDiscoveryClient, err := srd.NewConsulClient(*ptrConsulAddr)
	if err != nil {
		panic(err)
	}
	registration := srd.ServiceRegistration{
		ID:   *ptrServiceID,
		Name: *ptrServiceName,
		Addr: *ptrIp,
		Port: *ptrPort,
		Tag:  []string{*ptrServiceTag},
		TTL:  time.Second * 3,
		Check: func() (bool, error) {
			if serviceStatus == ServiceStatus_Pass {
				return true, nil
			}
			return false, errors.New("ready to offline")
		},
	}
	err = serviceDiscoveryClient.Register(&registration)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		rsp := "hello, this is " + *ptrServiceID
		fmt.Fprint(w, rsp)
	})

	http.HandleFunc("/offline", func(w http.ResponseWriter, r *http.Request) {
		serviceStatus = ServiceStatus_ReadyToOffline
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", registration.Addr, registration.Port), nil)
}
