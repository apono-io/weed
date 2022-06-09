package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apono-io/weed/pkg/k8s/addmissions"
	"io"
	admission "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	log "k8s.io/klog/v2"
	"net/http"
)

// admissionHandler represents the HTTP handler for an admission webhook
type admissionHandler struct {
	ctx     context.Context
	decoder runtime.Decoder
}

// AdmissionHandler returns an instance of AdmissionHandler
func AdmissionHandler(ctx context.Context) *admissionHandler {
	return &admissionHandler{
		ctx:     ctx,
		decoder: serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer(),
	}
}

// Serve returns a http.HandlerFunc for an admission webhook
func (h *admissionHandler) Serve(hook addmissions.Hook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprint("invalid method only POST requests are allowed"), http.StatusMethodNotAllowed)
			return
		}

		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, fmt.Sprint("only content type 'application/json' is supported"), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusBadRequest)
			return
		}

		var review admission.AdmissionReview
		if _, _, err := h.decoder.Decode(body, nil, &review); err != nil {
			http.Error(w, fmt.Sprintf("could not deserialize request: %v", err), http.StatusBadRequest)
			return
		}

		if review.Request == nil {
			http.Error(w, "malformed admission review: request is nil", http.StatusBadRequest)
			return
		}

		result, err := hook.Execute(h.ctx, review.Request)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		admissionResponse := review.DeepCopy()
		admissionResponse.Response = &admission.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: result.Allowed,
			Result:  &meta.Status{Message: result.Msg},
		}

		res, err := json.Marshal(admissionResponse)
		if err != nil {
			log.Error(err)
			http.Error(w, fmt.Sprintf("could not marshal response: %v", err), http.StatusInternalServerError)
			return
		}

		log.Infof("Webhook [%s - %s] - Allowed: %t", r.URL.Path, review.Request.Operation, result.Allowed)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(res)
	}
}
