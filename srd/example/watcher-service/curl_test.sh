#!/bin/bash

curl -H "Content-Type:application/json" -X POST -d '{"service": "hello-service"}' http://127.0.0.1:5001/watch | jq '.'
