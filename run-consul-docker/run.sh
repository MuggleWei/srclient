#!/bin/bash

sudo docker run -d \
	--name consul-agent-0 \
	--net=host \
	-e 'CONSUL_LOCAL_CONFIG={"skip_leave_on_interrupt": true}' \
	-v /docker/consul/consul1/data:/consul/data \
	registry.docker-cn.com/library/consul:1.0.2 agent -server -bind=127.0.0.1 -client=127.0.0.1 -ui -bootstrap-expect=1 -node=agent-0 -enable-script-checks=true
