package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	srd "github.com/MuggleWei/go-toy/srd"
)

type ReqWatch struct {
	Service string `json:"service"`
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)
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

	http.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var reqWatch ReqWatch
		err = json.Unmarshal(body, &reqWatch)
		if err != nil {
			panic(err)
		}

		log.Println(reqWatch.Service)

		services, err := serviceDiscoveryClient.GetService(reqWatch.Service, "")
		if err != nil {
			log.Print(err)
			panic(err)
		}

		bytes, err := json.Marshal(services)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(w, string(bytes))
	})

	log.Print(fmt.Sprintf("%s:%d", *ptrIp, *ptrPort))
	http.ListenAndServe(fmt.Sprintf("%s:%d", *ptrIp, *ptrPort), nil)
}
