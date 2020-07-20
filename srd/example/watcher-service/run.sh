#!/bin/bash

watcher_service_ip=127.0.0.1
consul_addr=127.0.0.1:8500

nohup ./watcher-service -ip=$watcher_service_ip -port=5001 -consul=$consul_addr &
