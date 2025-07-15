package module

import (
	"context"
	modulepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/module"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ModuleService) GetModule(ctx context.Context, req *modulepb.GetModuleRequest) (*modulepb.GetModuleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &modulepb.GetModuleResponse{}
	return resp, nil
}
