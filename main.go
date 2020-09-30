package main

import (
	"fmt"
	"github.com/inetaf/tcpproxy"
	"github.com/pires/go-proxyproto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	connectionErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tcpproxy_connect_errors",
		},
	)
)

func init() {
	prometheus.MustRegister(connectionErrorCounter)
}

func main() {
	list, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}
	listener := &proxyproto.Listener{
		Listener: list,
	}
	proxy := &tcpproxy.Proxy{
		ListenFunc: func(net, laddr string) (net.Listener, error) {
			return listener, nil
		},
	}
	dstTarget := &tcpproxy.DialProxy{
		Addr:            os.Getenv("REMOTE_ADDR"),
		KeepAlivePeriod: -1,
		DialTimeout:     5 * time.Second,
		OnDialError: func(src net.Conn, dstDialErr error) {
			log.Printf("failed connecting to: %v: %v", src.RemoteAddr(), dstDialErr)
			connectionErrorCounter.Inc()
		},
		ProxyProtocolVersion: 0,
	}
	proxy.AddRoute(":8000", dstTarget)
	go createPrometheusEndpoint()
	log.Fatal(proxy.Run())
}

func createPrometheusEndpoint() {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	srv := &http.Server{
		Handler:      mux,
		Addr:         ":8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
