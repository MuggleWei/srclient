#!/bin/bash

curl -H "Content-Type:application/json" -X POST -d '{"service": "hello-service"}' http://172.17.0.16:5001/watch | jq '.'
