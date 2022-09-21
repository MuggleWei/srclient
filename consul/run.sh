#!/bin/bash

origin_dir=$(readlink -f "$(dirname "$0")")
cd $origin_dir

echo "---------------------------"
echo "- run consul"
echo "---------------------------"

sudo docker compose up -d
