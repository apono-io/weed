package k8s

import (
	"context"
	"fmt"
	"github.com/apono-io/weed/pkg/k8s/addmissions/permissions"
	"github.com/apono-io/weed/pkg/k8s/handlers"
	"net/http"
)

// NewServer creates and return a http.Server
func NewServer(ctx context.Context, port int) *http.Server {
	// Routers
	ah := handlers.AdmissionHandler(ctx)
	mux := http.NewServeMux()
	mux.Handle("/healthz", handlers.HealthzHandler())
	mux.Handle("/validate/permissions", ah.Serve(permissions.NewValidatorHook()))

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
