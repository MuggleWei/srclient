# go toy 微服务
[english](readme_en.md)

## 描述
当前工程包含了

1. srd: 服务注册与发现
1. clb: 客户端负载均衡

服务注册与发现中心，暂时只支持consul

## 例子
在运行例子之前，首先确保consul已经正确的运行，如果你暂时没有玩过consul，建议可以按照以下步骤，先将consul运行起来
1. 安装docker
1. 运行run-consul-docker文件夹中的run.sh，注意根据自己的情况，调整脚本中的ip地址
1. sudo docker ps -a，查看consul是否已经正常运行

当consul已经顺利运行之后，可以试着运行一下本repo中的例子 (注意根据自己的情况，调整脚本中的ip地址)
1. 进入srd/example/hello-service，go build构建，之后运行run.sh，ps aux | grep hello-service，此时，可以看到已经运行了5个hello-service
1. 进入srd/example/watcher-service，go build构建，之后运行run.sh，此时，再运行curl_test.sh，可以看到，会将之前运行的hello-service集群的地址都打印出来
1. 进入clb/example/gate，go build构建，之后运行run.sh，此时运行 curl http://127.0.0.1:10102/hello，便可以观察到客户端负载的效果