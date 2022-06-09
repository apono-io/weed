package addmissions

import (
	"context"
	"fmt"
	admission "k8s.io/api/admission/v1"
)

type Hook struct {
	Create Validator
	Update Validator
}

func (h *Hook) Execute(ctx context.Context, request *admission.AdmissionRequest) (*ValidationResult, error) {
	operation := request.Operation

	var validator Validator
	switch operation {
	case admission.Create:
		validator = h.Create
	case admission.Update:
		validator = h.Update
	default:
		return &ValidationResult{Msg: fmt.Sprintf("Invalid operation: %s", operation)}, nil
	}

	if validator == nil {
		return nil, fmt.Errorf("operation %s is not registered", operation)
	}

	return validator(ctx, request)
}
