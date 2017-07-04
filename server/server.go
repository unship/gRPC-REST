package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	pb "github.com/biolee/gRPC-REST/proto"
	"github.com/biolee/swagger-ui-go-embed/ui"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type myService struct{}

func (m *myService) Echo(c context.Context, s *pb.EchoMessage) (*pb.EchoMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}

func newServer() *myService {
	return new(myService)
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func serveSwagger(mux *http.ServeMux) {
	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := ui.Handler
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func serveSwaggerJSON(mux *http.ServeMux) {
	data, err := pb.Asset("service.swagger.json")
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewReader(data))
	})
}

func serveGRPCGateWay(mux *http.ServeMux) {
	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: addr,
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	gwmux := runtime.NewServeMux()

	ctx := context.Background()

	err := pb.RegisterEchoServiceHandlerFromEndpoint(ctx, gwmux, addr, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}
	mux.Handle("/", gwmux)
}

func getHTTPMux() *http.ServeMux {
	mux := http.NewServeMux()

	serveSwagger(mux)
	serveSwaggerJSON(mux)
	serveGRPCGateWay(mux)

	return mux

}

func getGRPCServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certPool, addr)),
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterEchoServiceServer(grpcServer, newServer())

	return grpcServer
}

func getTCPServer() *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(getGRPCServer(), getHTTPMux()),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*keyPair},
			NextProtos:   []string{"h2"},
		},
	}
}

func ServeTCP() {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := getTCPServer()

	log.Fatal(srv.Serve(tls.NewListener(conn, srv.TLSConfig)))
}
