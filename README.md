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

### 运行：

正确运行后命令行会显示：

![image-20210927193647544](https://cdn.jsdelivr.net/gh/inkdp/CDN@main/img/20210927193647.png)

根据`proto`文件会生成`xxx.pb.go`和`xxx_grpc.pb.go`

![image-20210927194226535](https://cdn.jsdelivr.net/gh/inkdp/CDN@main/img/20210927194226.png)

## ⚠️注意⚠️

因为`proto`文件导入其他`proto`文件，以及文件目录和包名等一系列组合原因，可能出现正确编译，但生成的`go`文件出现包引用错误

![image-20210927214442275](https://cdn.jsdelivr.net/gh/inkdp/CDN@main/img/20210927214442.png)

![image-20210927214513593](https://cdn.jsdelivr.net/gh/inkdp/CDN@main/img/20210927214513.png)

此时我们需要手动修复包引用错误，修复方式根据自身运行目录以及业务场景各有不同，请自行修改源码编译后食用。贴一下个人的修复方案，在函数末尾添加如下脚本：

```go
cmd := exec.Command("/bin/bash", "-c", `find . -name '*.pb.go' | xargs perl -pi -e 's|cascade \"proto/cascade\"|cascade \"hitake/gateway/proto/cascade\"|g'`)
if _, err = cmd.Output(); err != nil {
   panic(err)
}
```

将引用的`proto/cascade`改为`hitake/gateway/proto/cascade`
