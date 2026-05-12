// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
	"github.com/open-telemetry/opentelemetry-operator/apis/v1beta1"
)

func TestClusterObservabilityWebhook_Validate(t *testing.T) {
	cases := []struct {
		name    string
		spec    v1alpha1.ClusterObservabilitySpec
		wantErr string
	}{
		{
			name: "valid_endpoint",
			spec: v1alpha1.ClusterObservabilitySpec{
				Exporter: v1beta1.AnyConfig{Object: map[string]any{
					"endpoint": "https://otel.example.com:4318",
				}},
			},
		},
		{
			name: "valid_traces_endpoint_only",
			spec: v1alpha1.ClusterObservabilitySpec{
				Exporter: v1beta1.AnyConfig{Object: map[string]any{
					"traces_endpoint": "https://otel.example.com:4318/v1/traces",
				}},
			},
		},
		{
			name: "valid_full_otlphttp_config",
			spec: v1alpha1.ClusterObservabilitySpec{
				Exporter: v1beta1.AnyConfig{Object: map[string]any{
					"endpoint":    "https://otel.example.com:4318",
					"compression": "gzip",
					"timeout":     "30s",
					"tls": map[string]any{
						"insecure": false,
					},
					"sending_queue": map[string]any{
						"enabled":    true,
						"queue_size": 1000,
					},
				}},
			},
		},
		{
			name:    "missing_exporter",
			spec:    v1alpha1.ClusterObservabilitySpec{},
			wantErr: "spec.exporter must be set and non-empty",
		},
		{
			name: "exporter_without_endpoint",
			spec: v1alpha1.ClusterObservabilitySpec{
				Exporter: v1beta1.AnyConfig{Object: map[string]any{
					"compression": "gzip",
					"timeout":     "30s",
				}},
			},
			wantErr: "spec.exporter must define one of",
		},
		{
			name: "exporter_with_empty_endpoint",
			spec: v1alpha1.ClusterObservabilitySpec{
				Exporter: v1beta1.AnyConfig{Object: map[string]any{
					"endpoint": "",
				}},
			},
			wantErr: "spec.exporter must define one of",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			co := &v1alpha1.ClusterObservability{Spec: tc.spec}
			w := &ClusterObservabilityWebhook{}

			_, err := w.ValidateCreate(context.Background(), co)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestClusterObservabilityWebhook_ValidateDelete_NoOp(t *testing.T) {
	w := &ClusterObservabilityWebhook{}
	warnings, err := w.ValidateDelete(context.Background(), &v1alpha1.ClusterObservability{})
	require.NoError(t, err)
	assert.Empty(t, warnings)
}
