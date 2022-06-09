package permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apono-io/weed/pkg/k8s/addmissions"
	"github.com/apono-io/weed/pkg/k8s/api"
	admission "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func NewValidatorHook() addmissions.Hook {
	validator := createValidator()
	return addmissions.Hook{
		Create: validator,
		Update: validator,
	}
}

func createValidator() addmissions.Validator {
	return func(_ context.Context, request *admission.AdmissionRequest) (*addmissions.ValidationResult, error) {
		var pom meta.PartialObjectMetadata
		err := json.Unmarshal(request.Object.Raw, &pom)
		if err != nil {
			return &addmissions.ValidationResult{Msg: fmt.Sprintf("failed to parse object, error: %v", err)}, nil
		}

		return validatePermissions(pom.ObjectMeta)
	}
}

func validatePermissions(object meta.ObjectMeta) (*addmissions.ValidationResult, error) {
	if permissionsCsv, exists := object.Annotations[api.RequiredPermissions]; exists {
		permissions := strings.Split(permissionsCsv, ",")
		missing := getMissingPermissions(permissions)
		if len(missing) > 0 {
			return &addmissions.ValidationResult{Msg: fmt.Sprintf("Missing permissions: %v", missing)}, nil
		}
	}

	return &addmissions.ValidationResult{Allowed: true}, nil
}

func getMissingPermissions(permissions []string) []string {
	existingPermissions := map[string]interface{}{
		"a": true,
	}

	var missingPermissions []string
	for _, permission := range permissions {
		if _, exists := existingPermissions[permission]; !exists {
			missingPermissions = append(missingPermissions, permission)
		}
	}

	return missingPermissions
}
