#!/bin/bash

watcher_service_ip=172.17.0.16
consul_addr=172.17.0.16:8500

nohup ./watcher-service -ip=$watcher_service_ip -port=5001 -consul=$consul_addr &
