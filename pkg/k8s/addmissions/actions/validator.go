package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/apono-io/weed/pkg/k8s/addmissions"
	"github.com/apono-io/weed/pkg/k8s/api"
	"github.com/apono-io/weed/pkg/weed"
	admission "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

var (
	noIamRoleFoundErr = errors.New("unable to validate manifest, could not find iam role")
)

func NewValidatorHook(ctx context.Context, clientset *kubernetes.Clientset, weedClient weed.WeedClient) addmissions.Hook {
	v := validator{
		ctx:        ctx,
		clientset:  clientset,
		weedClient: weedClient,
	}

	return addmissions.Hook{
		Create: v.validate,
		Update: v.validate,
	}
}

type validator struct {
	ctx        context.Context
	clientset  *kubernetes.Clientset
	weedClient weed.WeedClient
}

func (v *validator) validate(_ context.Context, request *admission.AdmissionRequest) (*addmissions.ValidationResult, error) {
	var template corev1.PodTemplateSpec
	switch request.Kind {
	case metav1.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}:
		err := json.Unmarshal(request.Object.Raw, &template)
		if err != nil {
			return &addmissions.ValidationResult{Msg: fmt.Sprintf("failed to parse object, error: %v", err)}, nil
		}
	case metav1.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "Deployment"}:
		var dp appsv1.Deployment
		err := json.Unmarshal(request.Object.Raw, &dp)
		if err != nil {
			return &addmissions.ValidationResult{Msg: fmt.Sprintf("failed to parse object, error: %v", err)}, nil
		}
		template = dp.Spec.Template
	case metav1.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "StatefulSet"}:
		var ss appsv1.StatefulSet
		err := json.Unmarshal(request.Object.Raw, &ss)
		if err != nil {
			return &addmissions.ValidationResult{Msg: fmt.Sprintf("failed to parse object, error: %v", err)}, nil
		}

		template = ss.Spec.Template
	}

	return v.validateActions(template)
}

func (v *validator) validateActions(template corev1.PodTemplateSpec) (*addmissions.ValidationResult, error) {
	if permissionsCsv, exists := template.Annotations[api.RequiredActions]; exists {
		if strings.TrimSpace(permissionsCsv) == "" {
			return &addmissions.ValidationResult{Allowed: true}, nil
		}

		iamRoleArn, err := v.extractIamRoleArn(template)
		if err != nil {
			return nil, err
		}

		permissions := strings.Split(permissionsCsv, ",")
		missing, err := v.checkMissingActions(iamRoleArn, permissions)
		if err != nil {
			return nil, err
		}

		if len(missing) > 0 {
			return &addmissions.ValidationResult{Msg: fmt.Sprintf("Missing actions: %v", missing)}, nil
		}
	}

	return &addmissions.ValidationResult{Allowed: true}, nil
}

func (v *validator) checkMissingActions(iamRoleArn string, requiredPermissions []string) ([]string, error) {
	klog.Infof("Checking required actions for role: %s", iamRoleArn)

	find, err := v.weedClient.Find(requiredPermissions, iamRoleArn)
	if err != nil {
		return nil, err
	}

	return find.Missing, nil
}

func (v *validator) extractIamRoleArn(template corev1.PodTemplateSpec) (string, error) {
	if arn, exists := template.Annotations[api.RoleArn]; exists && strings.TrimSpace(arn) != "" {
		return strings.TrimSpace(arn), nil
	}

	if arn, exists := template.Annotations[api.Kube2IamRoleArn]; exists && strings.TrimSpace(arn) != "" {
		return strings.TrimSpace(arn), nil
	}

	return v.extractServiceAccountIamRoleArn(template)
}

func (v *validator) extractServiceAccountIamRoleArn(template corev1.PodTemplateSpec) (string, error) {
	serviceAccountName := template.Spec.ServiceAccountName
	if strings.TrimSpace(serviceAccountName) == "" {
		return "", noIamRoleFoundErr
	}

	serviceAccount, err := v.clientset.CoreV1().ServiceAccounts(template.Namespace).Get(v.ctx, serviceAccountName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if arn, exists := serviceAccount.Annotations[api.ServiceAccountIamRoleArn]; exists && strings.TrimSpace(arn) != "" {
		return strings.TrimSpace(arn), nil
	}

	return "", noIamRoleFoundErr
}
