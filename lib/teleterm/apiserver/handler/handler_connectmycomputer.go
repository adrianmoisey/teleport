// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"strings"

	"github.com/gravitational/trace"

	api "github.com/gravitational/teleport/gen/proto/go/teleport/lib/teleterm/v1"
)

func (s *Handler) CreateConnectMyComputerRole(ctx context.Context, req *api.CreateConnectMyComputerRoleRequest) (*api.CreateConnectMyComputerRoleResponse, error) {
	res, err := s.DaemonService.CreateConnectMyComputerRole(ctx, req)
	return res, trace.Wrap(err)
}

func (s *Handler) CreateConnectMyComputerNodeToken(ctx context.Context, req *api.CreateConnectMyComputerNodeTokenRequest) (*api.CreateConnectMyComputerNodeTokenResponse, error) {
	token, err := s.DaemonService.CreateConnectMyComputerNodeToken(ctx, req.GetRootClusterUri())
	if err != nil {
		return nil, trace.Wrap(err)
	}

	apiLabels := APILabels{}
	for labelName, labelValues := range token.Labels {
		apiLabels = append(apiLabels, &api.Label{
			Name:  labelName,
			Value: strings.Join(labelValues, " "),
		})
	}

	response := &api.CreateConnectMyComputerNodeTokenResponse{
		Token:  token.Token,
		Labels: apiLabels,
	}

	return response, nil
}

func (s *Handler) DeleteConnectMyComputerToken(ctx context.Context, req *api.DeleteConnectMyComputerTokenRequest) (*api.DeleteConnectMyComputerTokenResponse, error) {
	res, err := s.DaemonService.DeleteConnectMyComputerToken(ctx, req)
	return res, trace.Wrap(err)
}
