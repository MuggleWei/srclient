#!/bin/bash

hello_service_ip="172.17.0.16"
consul_addr="172.17.0.16:8500"

nohup ./hello-service -ip=$hello_service_ip -port=6001 -consul=$consul_addr -service.name=hello-service -service.id=hello-service-1 &
nohup ./hello-service -ip=$hello_service_ip -port=6002 -consul=$consul_addr -service.name=hello-service -service.id=hello-service-2 &
nohup ./hello-service -ip=$hello_service_ip -port=6003 -consul=$consul_addr -service.name=hello-service -service.id=hello-service-3 &
nohup ./hello-service -ip=$hello_service_ip -port=6004 -consul=$consul_addr -service.name=hello-service -service.id=hello-service-4 &
nohup ./hello-service -ip=$hello_service_ip -port=6005 -consul=$consul_addr -service.name=hello-service -service.id=hello-service-5 &
