## srclient
srclient是一个服务注册与发现中心的client, 此项目分为两部分

1. srd: 服务注册与发现
1. clb: 客户端负载均衡

## 例子
1. 首先确保docker以及docker-compose已经正确安装, 接着运行`consul/run.sh`, 启动一个consul服务
1. 进入srd/example/hello-service，go build构建，之后运行run.sh，ps aux | grep hello-service，此时，可以看到已经运行了5个hello-service
1. 进入srd/example/watcher-service，go build构建，之后运行run.sh，此时，再运行curl_test.sh，可以看到，会将之前运行的hello-service集群的地址都打印出来
1. 进入clb/example/gate，go build构建，之后运行run.sh，此时运行 curl http://127.0.0.1:10102/hello，便可以观察到客户端负载的效果
