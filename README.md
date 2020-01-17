#### 1.1 特别需要注意的几个点 

- 参考文献:https://studygolang.com/articles/11923?p=1
- 参考文献:https://segmentfault.com/a/1190000013408485
- 参考文献:https://github.com/grpc-ecosystem/grpc-gateway

安装grpc-gateway:

```shell
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
```

项目目录结构：

```shell
└── src
    └── grpc-helloworld-gateway
        ├── gateway
        │   └── main.go
        ├── greeter_server
        │   └── main.go
        └── helloworld
            ├── helloworld.pb.go
            ├── helloworld.pb.gw.go
            └── helloworld.proto
```

google/api文件：

```shell
$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api
```

将文件拷贝到当前目录下:

```shell
# 拷贝以下两个目录
annotations.proto
http.proto
# 编译google.api [包含annotations.pb.go,http.pb.go文件]
cd helloworld
protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto
```

==还需要注意一点===

```shell
==将生成好 [包含annotations.pb.go,http.pb.go文件]== 拷贝一份放在这里
GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api
```

生成pb文件:

```shell
cd src/grpc-helloworld-gateway
# 生成helloworld.pb.go
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
helloworld/helloworld.proto 
# 生成helloworld.pb.gw.go
cd helloworld
protoc --grpc-gateway_out=logtostderr=true:. ./helloworld.proto
# 生成helloworld.swagger.json
cd src/grpc-helloworld-gateway

protoc -I/usr/local/include -I. \
-I$GOPATH/src  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--swagger_out=logtostderr=true:. \
helloworld/helloworld.proto
```

==这里需要注意一下== [grpc注册的端口地址一定要一致gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, "localhost:8181", opts)]

```go
mux :=  runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
		}
		//这里特别注意，注册端口一定要与grpc 服务端口一致 net.Listen("tcp", "localhost:8181")
		err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, "localhost:8181", opts)
		if err!=nil {
			panic(err)
		}
		//cors := cors.AllowAll()
		//serverMux := http.NewServeMux()
		//serverMux.Handle("/", cors.Handler(grpcMux))
		 http.ListenAndServe("localhost:8182", mux)
```

#### 1.2 etcd集群的搭建

```shell
$ docker pull quay.io/coreos/etcd
$ docker-compose up
$ docker ps 
# 验证集群
$ curl -L http://127.0.0.1:32787/v2/members
$ curl -L http://127.0.0.1:32789/v2/members
$ curl -L http://127.0.0.1:32791/v2/members
# 也可以使用etcdctl：
$ docker exec -t etcd1 etcdctl member list
```





