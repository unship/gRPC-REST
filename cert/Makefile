GOBINDATA_TARGET = embedded.go

.PHONY: all
all: generate


.PHONY: cert
cert:clean
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -key server.key -out server.pem -days 3650

.PHONY: generate
generate:$(GOBINDATA_TARGET)


.PHONY: $(GOBINDATA_TARGET)
$(GOBINDATA_TARGET):cert
	go-bindata -nometadata -pkg cert -o $@  server.key server.pem


.PHONY: clean
clean:
	rm -f server.key server.pem embedded.go