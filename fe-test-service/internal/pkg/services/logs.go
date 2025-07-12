// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package services

import (
	"context"
	_ "embed"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/unix4ever/fe-test-service/api/logs"
)

//go:embed data/talos.log
var data string

// LogsService implements LogsServiceServer.
type LogsService struct {
	logs.UnimplementedLogsServiceServer

	lines []string
}

// NewLogsService creates a new LogsService server.
func NewLogsService() *LogsService {
	return &LogsService{
		lines: strings.Split(strings.TrimSpace(data), "\n"),
	}
}

// List implements LogsServiceServer.
func (ls *LogsService) List(_ context.Context, request *logs.ListRequest) (*logs.ListResponse, error) {
	if request.Offset < 0 {
		return nil, status.Error(codes.InvalidArgument, "the offset can not be lower than 0")
	}

	if request.Limit < 0 {
		return nil, status.Error(codes.InvalidArgument, "the limit can not be lower than 0")
	}

	if request.Offset >= int32(len(ls.lines)) {
		return &logs.ListResponse{
			Lines: make([]string, 0),
		}, nil
	}

	res := ls.lines

	limit := len(res)
	offset := 0

	if request.Offset != 0 {
		offset = int(request.Offset)
	}

	if request.Limit != 0 {
		limit = offset + int(request.Limit)

		if limit > len(ls.lines) {
			limit = len(ls.lines)
		}
	}

	return &logs.ListResponse{
		Lines: res[offset:limit],
	}, nil
}
