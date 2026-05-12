// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-opentelemetry-io-v1alpha1-clusterobservability,mutating=false,failurePolicy=fail,groups=opentelemetry.io,resources=clusterobservabilities,versions=v1alpha1,name=vclusterobservabilitycreateupdate.kb.io,sideEffects=none,admissionReviewVersions=v1
// +kubebuilder:webhook:verbs=delete,path=/validate-opentelemetry-io-v1alpha1-clusterobservability,mutating=false,failurePolicy=ignore,groups=opentelemetry.io,resources=clusterobservabilities,versions=v1alpha1,name=vclusterobservabilitydelete.kb.io,sideEffects=none,admissionReviewVersions=v1
// +kubebuilder:object:generate=false

// ClusterObservabilityWebhook validates ClusterObservability resources.
// The exporter spec is opaque (passed to the otlphttp exporter as-is) so
// validation only enforces that an endpoint of some kind is configured.
type ClusterObservabilityWebhook struct {
	logger logr.Logger
	scheme *runtime.Scheme
}

var _ admission.Validator[*v1alpha1.ClusterObservability] = &ClusterObservabilityWebhook{}

// NewClusterObservabilityWebhook constructs a new validating webhook.
func NewClusterObservabilityWebhook(logger logr.Logger, scheme *runtime.Scheme) *ClusterObservabilityWebhook {
	return &ClusterObservabilityWebhook{
		logger: logger,
		scheme: scheme,
	}
}

// ValidateCreate validates the ClusterObservability resource on create.
func (*ClusterObservabilityWebhook) ValidateCreate(_ context.Context, co *v1alpha1.ClusterObservability) (admission.Warnings, error) {
	return nil, validate(co)
}

// ValidateUpdate validates the ClusterObservability resource on update.
func (*ClusterObservabilityWebhook) ValidateUpdate(_ context.Context, _, co *v1alpha1.ClusterObservability) (admission.Warnings, error) {
	return nil, validate(co)
}

// ValidateDelete is a no-op (delete uses failurePolicy=ignore).
func (*ClusterObservabilityWebhook) ValidateDelete(_ context.Context, _ *v1alpha1.ClusterObservability) (admission.Warnings, error) {
	return nil, nil
}

// endpointKeys are the otlphttp exporter fields any one of which is sufficient
// to give the collector somewhere to send telemetry. The list mirrors the
// public otlphttpexporter README.
var endpointKeys = []string{
	"endpoint",
	"traces_endpoint",
	"metrics_endpoint",
	"logs_endpoint",
	"profiles_endpoint",
}

// validate checks the minimum requirements for a ClusterObservability spec.
// Schema-level validation (types, enums, etc.) is delegated to the collector
// at config-load time; the webhook only catches misconfigurations that would
// otherwise leave the operator silently producing a non-functional pipeline.
func validate(co *v1alpha1.ClusterObservability) error {
	exp := co.Spec.Exporter.Object
	if len(exp) == 0 {
		return errors.New("spec.exporter must be set and non-empty")
	}
	for _, k := range endpointKeys {
		if v, ok := exp[k]; ok {
			if s, isStr := v.(string); isStr && s != "" {
				return nil
			}
		}
	}
	return fmt.Errorf("spec.exporter must define one of %v", endpointKeys)
}

// SetupClusterObservabilityWebhook wires the webhook into the manager.
func SetupClusterObservabilityWebhook(mgr ctrl.Manager) error {
	w := NewClusterObservabilityWebhook(
		mgr.GetLogger().WithValues("handler", "ClusterObservabilityWebhook"),
		mgr.GetScheme(),
	)
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.ClusterObservability{}).
		WithValidator(w).
		Complete()
}
