# webchat box
将微信在容器中运行，暴露出HTTP REST接口

# hook目录

## vcpkg
```
git clone https://github.com/microsoft/vcpkg
.\vcpkg\bootstrap-vcpkg.bat

vcpkg install protobuf[zlib]:x86-windows-static
vcpkg install spdlog:x86-windows-static
vcpkg install nng:x86-windows-static
vcpkg install magic-enum:x86-windows-static
vcpkg integrate install
```

## protobuf
```
cd hook\rpc\proto
..\tool\protoc --nanopb_out=. wcf.proto
```

# rest目录
```
cd rest/internal/rpc

protoc --go_out=. rpc.proto
```

