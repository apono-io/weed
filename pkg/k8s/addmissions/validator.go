package addmissions

import (
	"context"
	admission "k8s.io/api/admission/v1"
)

type Validator func(ctx context.Context, request *admission.AdmissionRequest) (*ValidationResult, error)

type ValidationResult struct {
	Allowed bool
	Msg     string
}
