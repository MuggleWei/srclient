package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	clb "github.com/MuggleWei/srclient/clb"
	srd "github.com/MuggleWei/srclient/srd"
)

var transport *http.Transport

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)

	transport = &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}
}

func main() {
	ptrIp := flag.String("ip", "127.0.0.1", "run ip address")
	ptrPort := flag.Int("port", 10102, "listen port")
	ptrConsulAddr := flag.String("consul", "127.0.0.1:8500", "consul address")

	flag.Parse()

	serviceDiscoveryClient, err := srd.NewConsulClient(*ptrConsulAddr)
	if err != nil {
		panic(err)
	}
	clientLB := clb.NewClientLoadBalancer(serviceDiscoveryClient, time.Second*3)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		addr, err := clientLB.GetService("hello-service")
		if err != nil {
			panic(err)
		}

		url := "http://" + addr + "/hello"
		client := &http.Client{Transport: transport}
		rsp, err := client.Get(url)
		if err != nil {
			panic(err)
		}
		defer rsp.Body.Close()

		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(w, string(bytes))
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", *ptrIp, *ptrPort), nil)
}
