// Copyright 2023 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestServerConfigComplete(t *testing.T) {
	require := require.New(t)
	c := &ServerConfig{}
	err := c.Complete()
	require.NoError(err)

	require.EqualValues("token", c.Auth.Method)
	require.Equal(true, lo.FromPtr(c.Transport.TCPMux))
	require.Equal(true, lo.FromPtr(c.DetailedErrorsToClient))
	require.Equal(true, lo.FromPtr(c.OpenGFW.Enable))
	require.Equal(true, lo.FromPtr(c.OpenGFW.ProxyPolicy.Enable))
	require.Equal(85, c.OpenGFW.ProxyPolicy.HighThreshold)
	require.Equal(70, c.OpenGFW.ProxyPolicy.MediumThreshold)
	require.Equal(true, lo.FromPtr(c.OpenGFW.TrafficPolicy.Enable))
	require.Equal(true, lo.FromPtr(c.OpenGFW.TrafficPolicy.BlockWebTCP))
}

func TestAuthServerConfig_Complete(t *testing.T) {
	require := require.New(t)
	cfg := &AuthServerConfig{}
	err := cfg.Complete()
	require.NoError(err)
	require.EqualValues("token", cfg.Method)
}

func TestOpenGFWConfigComplete(t *testing.T) {
	require := require.New(t)

	cfg := OpenGFWServerConfig{
		Enable: lo.ToPtr(false),
		ProxyPolicy: OpenGFWProxyPolicyConfig{
			Enable:          lo.ToPtr(false),
			HighThreshold:   50,
			MediumThreshold: 70,
		},
		TrafficPolicy: OpenGFWTrafficPolicyConfig{
			Enable:      lo.ToPtr(false),
			BlockWebTCP: lo.ToPtr(false),
		},
	}
	cfg.Complete()

	require.Equal(false, lo.FromPtr(cfg.Enable))
	require.Equal(false, lo.FromPtr(cfg.ProxyPolicy.Enable))
	require.Equal(50, cfg.ProxyPolicy.HighThreshold)
	require.Equal(49, cfg.ProxyPolicy.MediumThreshold)
	require.Equal(false, lo.FromPtr(cfg.TrafficPolicy.Enable))
	require.Equal(false, lo.FromPtr(cfg.TrafficPolicy.BlockWebTCP))
}
