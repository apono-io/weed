package main

import (
	"context"
	"flag"
	"github.com/apono-io/weed/pkg/core"
	"github.com/apono-io/weed/pkg/k8s"
	log "k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var tlsCert, tlsKey string
	var port int
	flag.StringVar(&tlsCert, "tls-cert", "/etc/certs/tls.crt", "Path to the TLS certificate")
	flag.StringVar(&tlsKey, "tls-key", "/etc/certs/tls.key", "Path to the TLS key")
	flag.IntVar(&port, "port", 8443, "HTTPS Port to listen on")
	flag.Parse()

	log.Infof("Starting IAM Enforcer v%s (commit: %s, built at: %s)", core.Version, core.Commit, core.BuildDate)

	ctx, cancel := context.WithCancel(context.Background())
	server := k8s.NewServer(ctx, port)
	go func() {
		if err := server.ListenAndServeTLS(tlsCert, tlsKey); err != nil {
			log.Errorf("Failed to listen and serve: %v", err)
			os.Exit(1)
		}
	}()

	log.Infof("Server running on port: %d", port)

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	cancel()
	log.Infof("Shutdown gracefully...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
}
