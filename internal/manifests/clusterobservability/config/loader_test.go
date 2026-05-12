// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
	"github.com/open-telemetry/opentelemetry-operator/apis/v1beta1"
)

// TestBuildExportersConfig verifies the user-provided AnyConfig is embedded
// verbatim under the otlphttp key without any field-by-field translation.
func TestBuildExportersConfig(t *testing.T) {
	cases := []struct {
		name string
		in   v1alpha1.ClusterObservabilitySpec
		want map[string]any
	}{
		{
			name: "endpoint_only",
			in: v1alpha1.ClusterObservabilitySpec{
				Exporter: v1beta1.AnyConfig{Object: map[string]any{
					"endpoint": "https://otel.example.com:4318",
				}},
			},
			want: map[string]any{
				"otlphttp": map[string]any{
					"endpoint": "https://otel.example.com:4318",
				},
			},
		},
		{
			name: "passthrough_unknown_fields",
			in: v1alpha1.ClusterObservabilitySpec{
				Exporter: v1beta1.AnyConfig{Object: map[string]any{
					"endpoint":          "https://otel.example.com:4318",
					"compression":       "zstd", // a future compression algo unknown to operator
					"some_future_field": 42,
					"deeply": map[string]any{
						"nested": map[string]any{"value": true},
					},
				}},
			},
			want: map[string]any{
				"otlphttp": map[string]any{
					"endpoint":          "https://otel.example.com:4318",
					"compression":       "zstd",
					"some_future_field": 42,
					"deeply": map[string]any{
						"nested": map[string]any{"value": true},
					},
				},
			},
		},
		{
			name: "nil_object_yields_empty_otlphttp",
			in:   v1alpha1.ClusterObservabilitySpec{},
			want: map[string]any{
				"otlphttp": map[string]any{},
			},
		},
	}

	loader := NewConfigLoader()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := loader.buildExportersConfig(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}
