# DEV ENV
```bash
export GlideVendor="./vendor"
glide install
go install $GlideVendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install $GlideVendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
```


port from https://github.com/biolee/grpc-gateway-example