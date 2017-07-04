GOBINDATA_TARGET = embedded.go

.PHONY: all
all: generate


.PHONY: cert
cert:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -key server.key -out server.pem -days 3650

.PHONY: generate
generate:$(GOBINDATA_TARGET) cert


.PHONY: $(GOBINDATA_TARGET)
$(GOBINDATA_TARGET):
	go-bindata -nometadata -pkg cert -o $@  server.key server.pem
