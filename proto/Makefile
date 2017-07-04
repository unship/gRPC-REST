GOBINDATA_TARGET = embedded.go

.DEFAULT_GOAL := all
.PHONY: all
all: protoc $(GOBINDATA_TARGET)

.PHONY: protoc
protoc:
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. \
		service.proto
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		service.proto
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		service.proto
	go generate .

.PHONY: $(GOBINDATA_TARGET)
$(GOBINDATA_TARGET):protoc
	go-bindata -nometadata -pkg proto -o $@  service.swagger.json


