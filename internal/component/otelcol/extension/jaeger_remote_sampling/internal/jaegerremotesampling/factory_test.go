// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package jaegerremotesampling

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/extension/extensiontest"
)

func TestCreateDefaultConfig(t *testing.T) {
	// prepare and test
	expected := &Config{
		HTTPServerConfig: &confighttp.ServerConfig{Endpoint: ":5778"},
		GRPCServerConfig: &configgrpc.ServerConfig{NetAddr: confignet.AddrConfig{
			Endpoint:  ":14250",
			Transport: confignet.TransportTypeTCP,
		}},
	}

	// test
	cfg := createDefaultConfig()

	// verify
	assert.Equal(t, expected, cfg)
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestCreateExtension(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	ext, err := createExtension(t.Context(), extensiontest.NewNopSettings(component.MustNewType("jaegerremotesampling")), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, ext)
}

func TestNewFactory(t *testing.T) {
	f := NewFactory()
	assert.NotNil(t, f)
}
