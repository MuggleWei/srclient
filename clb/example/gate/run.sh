#!/bin/bash

nohup ./gate -ip=0.0.0.0 -port=10102 -consul=127.0.0.1:8500 &
