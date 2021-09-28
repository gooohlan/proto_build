# proto_build

摆脱繁琐编译命令，实现proto编译自由

## 准备：

确保你有完整的GRPC运行环境(`protoc`，`protoc-gen-go`，`protoc-gen-go-grpc`)，不完整的请自行安装，下面的安装方式可能有误，出现错误请查找网上教程

- **protoc**：[下载最新](https://github.com/protocolbuffers/protobuf/releases/)的`protoc`放入`bin`目录下
- **protoc-gen-go**：`go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26`
- **protoc-gen-go-grpc**：`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1`

## 使用：

1. 源码运行：下载源码到项目中任意目录，运行`main`函数即可
2. 命令行运行：下载源码编译或直接下载二进制包到项目文件，运行执行文件即可
3. 自动运行：[Goland自动编译proto文件](https://www.inkdp.cn/skill/back-end/49446.html)

## 运行：

正确运行后命令行会提示：`生成proto.go成功`

根据`proto`文件会生成`xxx.pb.go`和`xxx_grpc.pb.go`

![image-20210927194226535](https://cdn.jsdelivr.net/gh/inkdp/CDN@main/img/20210927194226.png)

## ⚠️注意⚠️

因为`proto`文件导入其他`proto`文件，以及文件目录和包名等一系列组合原因，可能出现正确编译，但生成的`go`文件出现包引用错误，定义`proto`文件时，`option go_package`请指定完整包名
