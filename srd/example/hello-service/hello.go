package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	srd "github.com/MuggleWei/go-toy/srd"
)

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
		ID:    *ptrServiceID,
		Name:  *ptrServiceName,
		Addr:  *ptrIp,
		Port:  *ptrPort,
		Tag:   []string{*ptrServiceTag},
		TTL:   time.Second * 3,
		Check: func() (bool, error) { return true, nil },
	}
	err = serviceDiscoveryClient.Register(&registration)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		rsp := "hello, this is " + *ptrServiceID
		fmt.Fprint(w, rsp)
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", registration.Addr, registration.Port), nil)
}
